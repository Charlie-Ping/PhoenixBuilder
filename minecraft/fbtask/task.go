package fbtask
import (
	"phoenixbuilder/minecraft/mctype"
	"phoenixbuilder/minecraft/parse"
	"phoenixbuilder/minecraft/builder"
	"phoenixbuilder/minecraft/command"
	"phoenixbuilder/minecraft/configuration"
	"phoenixbuilder/minecraft"
	"go.uber.org/atomic"
	"sync"
	"fmt"
	"time"
	"runtime"
	"strings"
	"github.com/google/uuid"
)

const (
	TaskStateUnknown     = 0
	TaskStateRunning     = 1
	TaskStatePaused      = 2
	TaskStateDied        = 3
	TaskStateCalculating = 4
	TaskStateSpecialBrk  = 5
)

type Task struct {
	TaskId int64
	CommandLine string
	OutputChannel chan *mctype.Module
	ContinueLock sync.Mutex
	State byte
	Type byte
	AsyncInfo
	Config *configuration.FullConfig
}

type AsyncInfo struct {
	Built int
	Total int
	BeginTime time.Time
}

var TaskIdCounter *atomic.Int64 = atomic.NewInt64(0)
var TaskMap sync.Map
var BrokSender chan string = make(chan string)

func GetStateDesc(st byte) string {
	if st == 0 {
		return "Unknown"
	}else if st==1 {
		return "Running"
	}else if st==2 {
		return "Paused"
	}else if st==3 {
		return "Died"
	}else if st==4 {
		return "Calculating"
	}else if st==5 {
		return "SpecialTask:Breaking"
	}
	return "???????"
}

func (task *Task) Finalize() {
	task.State = TaskStateDied
	TaskMap.Delete(task.TaskId)
}

func (task *Task) Pause() {
	if task.State == TaskStatePaused {
		return
	}
	task.ContinueLock.Lock()
	if task.State == TaskStateDied {
		task.ContinueLock.Unlock()
		return
	}
	task.State = TaskStatePaused
}

func (task *Task) Resume() {
	if task.State != TaskStatePaused {
		return
	}
	if task.Type==mctype.TaskTypeAsync {
		task.AsyncInfo.Total-=task.AsyncInfo.Built
		task.AsyncInfo.Built=0
	}
	task.State = TaskStateRunning
	task.ContinueLock.Unlock()
}

func (task *Task) Break() {
	if task.OutputChannel==nil {
		task.State=TaskStateSpecialBrk
		return
	}
	if task.State != TaskStatePaused {
		task.Pause()
	}
	if task.State == TaskStateDied {
		return
	}
	chann := task.OutputChannel
	for {
		blk, ok := <- chann
		if !ok {
			break
		}
		if false {
			fmt.Printf("%v\n",blk)
		}
	}
	if task.Type==mctype.TaskTypeAsync {
		// Avoid progress displaying
		if task.State != TaskStatePaused {
			return
		}
		task.State = TaskStateCalculating
		task.ContinueLock.Unlock()
		return
	}
	task.Resume()
}

func FindTask(taskId int64) *Task {
	t, _ := TaskMap.Load(taskId)
	ta, _ := t.(*Task)
	return ta
}

func CreateTask(commandLine string, conn *minecraft.Conn) *Task {
	cfg, err := parse.Parse(commandLine, configuration.GlobalFullConfig().Main())
	if err!=nil {
		command.Tellraw(conn, fmt.Sprintf("Failed to parse command: %v",err))
		return nil
	}
	fcfg := configuration.ConcatFullConfig(cfg, configuration.GlobalFullConfig().Delay())
	dcfg := fcfg.Delay()
	/*if cfg.Execute == "" {
		return nil
	}
	Needless since it will be checked in function module.
	*/
	und, _ := uuid.NewUUID()
	command.SendWSCommand("gamemode c", und, conn)
	blockschannel := make(chan *mctype.Module, 10240)
	task := &Task {
		TaskId: TaskIdCounter.Add(1),
		CommandLine: commandLine,
		OutputChannel: blockschannel,
		State: TaskStateCalculating,
		Type: configuration.GlobalFullConfig().Global().TaskCreationType,
		Config: fcfg,
	}
	taskid := task.TaskId
	TaskMap.Store(taskid, task)
	var asyncblockschannel chan *mctype.Module
	if task.Type==mctype.TaskTypeAsync {
		asyncblockschannel=blockschannel
		blockschannel=make(chan *mctype.Module)
		task.OutputChannel=blockschannel
		go func() {
			var blocks []*mctype.Module
			for {
				curblock, ok := <-asyncblockschannel
				if !ok {
					break
				}
				blocks=append(blocks,curblock)
			}
			task.State=TaskStateRunning
			t1 := time.Now()
			total := len(blocks)
			task.AsyncInfo=AsyncInfo {
				Built: 0,
				Total: total,
				BeginTime: t1,
			}
			for _, blk := range blocks {
				blockschannel <- blk
				task.AsyncInfo.Built++
			}
			close(blockschannel)
		} ()
	}else{
		task.State=TaskStateRunning
	}
	go func() {
		t1 := time.Now()
		blkscounter := 0
		tothresholdcounter := 0
		for {
			task.ContinueLock.Lock()
			task.ContinueLock.Unlock()
			curblock, ok := <-blockschannel
			if !ok {
				if blkscounter == 0 {
					command.Tellraw(conn, fmt.Sprintf("[Task %d] Nothing generated.",taskid))
					runtime.GC()
					task.Finalize()
					return
				}
				timeUsed := time.Now().Sub(t1)
				command.Tellraw(conn, fmt.Sprintf("[Task %d] %v block(s) have been changed.", taskid, blkscounter))
				command.Tellraw(conn, fmt.Sprintf("[Task %d] Time used: %v second(s)", taskid, timeUsed.Seconds()))
				command.Tellraw(conn, fmt.Sprintf("[Task %d] Average speed: %v blocks/second", taskid, float64(blkscounter)/timeUsed.Seconds()))
				runtime.GC()
				task.Finalize()
				return
			}
			if blkscounter%20 == 0 {
				u_d, _ := uuid.NewUUID()
				command.SendWSCommand(fmt.Sprintf("tp %d %d %d",curblock.Point.X,curblock.Point.Y,curblock.Point.Z),u_d, conn)
			}
			blkscounter++
			request := command.SetBlockRequest(curblock, cfg)
			err := command.SendSizukanaCommand(request, conn)
			if err != nil {
				panic(err)
			}
			if dcfg.DelayMode==mctype.DelayModeContinuous {
				time.Sleep(time.Duration(dcfg.Delay) * time.Microsecond)
			}else if dcfg.DelayMode==mctype.DelayModeDiscrete {
				tothresholdcounter++
				if tothresholdcounter>=dcfg.DelayThreshold {
					tothresholdcounter=0
					time.Sleep(time.Duration(dcfg.Delay) * time.Second)
				}
			}
		}
	} ()
	go func() {
		if task.Type==mctype.TaskTypeAsync {
			err := builder.Generate(cfg, asyncblockschannel)
			close(asyncblockschannel)
			if err != nil {
				command.Tellraw(conn, fmt.Sprintf("[Task %d] Error: %v", taskid, err))
			}
			return
		}
		err := builder.Generate(cfg, blockschannel)
		close(blockschannel)
		if err != nil {
			command.Tellraw(conn, fmt.Sprintf("[Task %d] Error: %v", taskid, err))
		}
	} ()
	return task
}


func InitTaskStatusDisplay(conn *minecraft.Conn) {
	go func() {
		for {
			str:=<-BrokSender
			command.Tellraw(conn,str)
		}
	} ()
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			<- ticker.C
			if configuration.GlobalFullConfig().Global().TaskDisplayMode == mctype.TaskDisplayNo {
				continue
			}
			var displayStrs []string
			TaskMap.Range(func (_tid interface{}, _v interface{}) bool {
				tid, _:=_tid.(int64)
				v, _:=_v.(*Task)
				addstr:=fmt.Sprintf("Task ID %d - %s - %s [%s]",tid,v.Config.Main().Execute,GetStateDesc(v.State),mctype.MakeTaskType(v.Type))
				if v.Type==mctype.TaskTypeAsync && v.State == TaskStateRunning {
					addstr=fmt.Sprintf("%s\nProgress: %s",addstr,ProgressThemes[0](&v.AsyncInfo))
				}
				displayStrs=append(displayStrs,addstr)
				return true
			})
			if len(displayStrs) == 0 {
				continue
			}
			command.Title(conn,strings.Join(displayStrs,"\n"))
		}
	} ()
}
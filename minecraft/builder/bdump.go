package builder

import (
	"github.com/andybalholm/brotli"
	"phoenixbuilder/minecraft/mctype"
	"fmt"
	"os"
	"encoding/binary"
)

func ReadBrString(br *brotli.Reader) (string, error) {
	str:=""
	c:=make([]byte,1)
	for {
		_, err:=br.Read(c)
		if(err!=nil) {
			return "",err
		}
		if(c[0]==0) {
			break
		}
		str+=string(c)
	}
	return str, nil
}

func BDump(config *mctype.MainConfig, blc chan *mctype.Module) error {
	file, err:=os.OpenFile(config.Path,os.O_RDONLY,0644)
	if err!=nil {
		return fmt.Errorf("Failed to open file: %v",err)
	}
	defer file.Close()
	{
		header3bytes:=make([]byte,3)
		_, err:=file.Read(header3bytes)
		if(err!=nil){
			return fmt.Errorf("Failed to read file, early EOF? File may be corrupted")
		}
		if(string(header3bytes)!="BD@"){
			return fmt.Errorf("Not a bdx file (Invalid file header)")
		}
	}
	br := brotli.NewReader(file)
	{
		tempbuf:=make([]byte,4)
		_, err:=br.Read(tempbuf)
		if(err!=nil) {
			return fmt.Errorf("Invalid file")
		}
		if(string(tempbuf)!="BDX\x00"){
			return fmt.Errorf("Not a bdx file (Invalid inner header)")
		}
	}
	{
		author, err:=ReadBrString(br)
		if(err!=nil){
			return fmt.Errorf("Failed to read author info, file may be corrupted")
		}
		mctype.ForwardedBrokSender<-fmt.Sprintf("Author: %s\n",author)
	}
	curcmdbuf:=make([]byte,1)
	brushPosition:=[]int{0,0,0}
	var blocksStrPool []string
	prevCmd:=0
	for {
		prevCmd=int(curcmdbuf[0])
		_, err:=br.Read(curcmdbuf)
		if err!=nil {
			return fmt.Errorf("Failed to get construction command, file may be corrupted")
		}
		cmd:=curcmdbuf[0]
		if(cmd==88){
			break
		}
		if(cmd==1) {
			bstr,err:=ReadBrString(br)
			if(err!=nil){
				return fmt.Errorf("Failed to get argument for cmd[pos:0], file may be corrupted!")
			}
			blocksStrPool=append(blocksStrPool,bstr)
			continue
		}else if(cmd==2){
			rdst:=make([]byte,2)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos1], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint16(rdst)
			brushPosition[0]+=int(jumpval)
			brushPosition[1]=0
			brushPosition[2]=0
		}else if(cmd==3){
			brushPosition[0]++
			brushPosition[1]=0
			brushPosition[2]=0
		}else if(cmd==4){
			rdst:=make([]byte,2)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos2], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint16(rdst)
			brushPosition[1]+=int(jumpval)
			brushPosition[2]=0
		}else if(cmd==5){
			brushPosition[1]++
			brushPosition[2]=0
		}else if(cmd==6){
			rdst:=make([]byte,2)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos3], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint16(rdst)
			brushPosition[2]+=int(jumpval)
		}else if(cmd==7){
			rdst:=make([]byte,2)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos4], file may be corrupted")
			}
			blockId:=binary.BigEndian.Uint16(rdst)
			_, err=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos5], file may be corrupted")
			}
			if(int(blockId)>=len(blocksStrPool)){
				fmt.Printf("WARNING: Invalid command")
				continue
			}
			blockData:=binary.BigEndian.Uint16(rdst)
			blockName:=&blocksStrPool[int(blockId)]
			blc<-&mctype.Module {
				Block: &mctype.Block {
					Name: blockName,
					Data: int16(blockData),
				},
				Point: mctype.Position {
					X: brushPosition[0]+config.Position.X,
					Y: brushPosition[1]+config.Position.Y,
					Z: brushPosition[2]+config.Position.Z,
				},
			}
		}else if(cmd==8){
			brushPosition[2]++
		}else if(cmd==9){
			// Command: NOP
		}else if(cmd==10){
			rdst:=make([]byte,4)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos6], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint32(rdst)
			brushPosition[0]+=int(jumpval)
			brushPosition[1]=0
			brushPosition[2]=0
		}else if(cmd==11){
			rdst:=make([]byte,4)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos7], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint32(rdst)
			brushPosition[1]+=int(jumpval)
			brushPosition[2]=0
		}else if(cmd==12){
			rdst:=make([]byte,4)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos8], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint32(rdst)
			brushPosition[2]+=int(jumpval)
		}else if(cmd==13){
			fmt.Printf("WARNING: BDump/Import: Use of reserved command")
		}else if(cmd==14){
			brushPosition[0]++
		}else if(cmd==15){
			brushPosition[0]--
		}else if(cmd==16){
			brushPosition[1]++
		}else if(cmd==17){
			brushPosition[1]--
		}else if(cmd==18){
			brushPosition[2]++
		}else if(cmd==19){
			brushPosition[2]--
		}else if(cmd==20){
			rdst:=make([]byte,2)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos9], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint16(rdst)
			brushPosition[0]+=int(int16(jumpval))
		}else if(cmd==21){
			rdst:=make([]byte,4)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos10], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint32(rdst)
			brushPosition[0]+=int(int32(jumpval))
		}else if(cmd==22){
			rdst:=make([]byte,2)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos11], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint16(rdst)
			brushPosition[1]+=int(int16(jumpval))
		}else if(cmd==23){
			rdst:=make([]byte,4)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos12], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint32(rdst)
			brushPosition[1]+=int(int32(jumpval))
		}else if(cmd==24){
			rdst:=make([]byte,2)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos13], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint16(rdst)
			brushPosition[2]+=int(int16(jumpval))
		}else if(cmd==25){
			rdst:=make([]byte,4)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos14], file may be corrupted")
			}
			jumpval:=binary.BigEndian.Uint32(rdst)
			brushPosition[2]+=int(int32(jumpval))
		}else if(cmd==26){
			fbuf:=make([]byte,4)
			_, err:=br.Read(fbuf)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos15], file may be corrupted")
			}
			cbmode:=binary.BigEndian.Uint32(fbuf)
			command, err:=ReadBrString(br)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos16], file may be corrupted")
			}
			cusname, err:=ReadBrString(br)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos17], file may be corrupted")
			}
			lasout, err:=ReadBrString(br)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos18], file may be corrupted")
			}
			_, err=br.Read(fbuf)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos19], file may be corrupted")
			}
			tickdelay:=int32(binary.BigEndian.Uint32(fbuf))
			_, err=br.Read(fbuf)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos20], file may be corrupted")
			}
			fbools:=[]bool{false,false,false,false}
			if fbuf[0]==1 {
				fbools[0]=true
			}
			if fbuf[1]==1 {
				fbools[1]=true
			}
			if fbuf[2]==1 {
				fbools[2]=true
			}
			if fbuf[3]==1 {
				fbools[3]=true
			}
			cbdata:=&mctype.CommandBlockData {
				Mode: cbmode,
				Command: command,
				CustomName: cusname,
				LastOutput: lasout,
				TickDelay: tickdelay,
				ExecuteOnFirstTick: fbools[0],
				TrackOutput: fbools[1],
				Conditional: fbools[2],
				NeedRedstone: fbools[3],
			}
			blc<-&mctype.Module {
				CommandBlockData: cbdata,
				Point: mctype.Position {
					X: brushPosition[0]+config.Position.X,
					Y: brushPosition[1]+config.Position.Y,
					Z: brushPosition[2]+config.Position.Z,
				},
			}
		}else if(cmd==27){
			rdst:=make([]byte,2)
			_, err:=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos4], file may be corrupted")
			}
			blockId:=binary.BigEndian.Uint16(rdst)
			_, err=br.Read(rdst)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos5], file may be corrupted")
			}
			if(int(blockId)>=len(blocksStrPool)){
				fmt.Printf("WARNING: Invalid command")
				continue
			}
			blockData:=binary.BigEndian.Uint16(rdst)
			blockName:=&blocksStrPool[int(blockId)]
			cmdl:=&mctype.Module {
				Block: &mctype.Block {
					Name: blockName,
					Data: int16(blockData),
				},
				Point: mctype.Position {
					X: brushPosition[0]+config.Position.X,
					Y: brushPosition[1]+config.Position.Y,
					Z: brushPosition[2]+config.Position.Z,
				},
			}
			fbuf:=make([]byte,4)
			_, err=br.Read(fbuf)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos15], file may be corrupted")
			}
			cbmode:=binary.BigEndian.Uint32(fbuf)
			command, err:=ReadBrString(br)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos16], file may be corrupted")
			}
			cusname, err:=ReadBrString(br)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos17], file may be corrupted")
			}
			lasout, err:=ReadBrString(br)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos18], file may be corrupted")
			}
			_, err=br.Read(fbuf)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos19], file may be corrupted")
			}
			tickdelay:=int32(binary.BigEndian.Uint32(fbuf))
			_, err=br.Read(fbuf)
			if(err!=nil) {
				return fmt.Errorf("Failed to get argument for cmd[pos20], file may be corrupted")
			}
			fbools:=[]bool{false,false,false,false}
			if fbuf[0]==1 {
				fbools[0]=true
			}
			if fbuf[1]==1 {
				fbools[1]=true
			}
			if fbuf[2]==1 {
				fbools[2]=true
			}
			if fbuf[3]==1 {
				fbools[3]=true
			}
			cbdata:=&mctype.CommandBlockData {
				Mode: cbmode,
				Command: command,
				CustomName: cusname,
				LastOutput: lasout,
				TickDelay: tickdelay,
				ExecuteOnFirstTick: fbools[0],
				TrackOutput: fbools[1],
				Conditional: fbools[2],
				NeedRedstone: fbools[3],
			}
			cmdl.CommandBlockData=cbdata
			blc<-cmdl
		}else{
			fmt.Printf("WARNING: BDump/Import: Unimplemented method found : %d\n",cmd)
			fmt.Printf("WARNING: BDump/Import: Previous command is: %d\n",prevCmd)
			fmt.Printf("WARNING: BDump/Import: Trying to ignore, it will probably cause an error!\n")
		}
	}
	return nil
}
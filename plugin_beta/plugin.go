// 这写插件框架有点折磨啊
// 最终版!一定是最终版!不能重写重新构思了!
// 感谢2401PT, awaqwqa, xX7912Xx, CAIMEOX, LNSSPsd... 所有帮助过咱的人!
package plugin_beta

import (
	"fmt"
	"io/ioutil"
	"log"

	"os"
	"path/filepath"
	expand "phoenixbuilder/fastbuilder/plugin_structs"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
	"plugin"
	"sort"
	"sync"
)

// PacketSender: plugin -> main
// PacketReceiver: main -> plugin

type Plugin struct {
	singleton bool
	block     bool // true if it blocks packets.
	priority  int
	handleNum int64
	name      string

	// handleWg sync.WaitGroup

	// locked when PluginManager Register and do nothing when PluginManager notify(packets)
	regMu sync.RWMutex
	// plugins get it and push Packet to main process.
	packetReceivers []chan packet.Packet

	// rule func(pk *packet.Packet) bool
}

type IPlugin interface {
	Init(*PluginManager)
	Handler(*PluginManager, packet.Packet)

	// assert pk and return ok. It should be simplified.
	// Should it be a single function instead of a method?
	Rule(packet.Packet) bool
}

type PluginManager struct {
	Method         expand.PluginBridge
	Expand         ExpandPluginBridgeImpl
	conn           *minecraft.Conn
	Logger         log.Logger
	regMu          sync.RWMutex
	pluginPriority []IPlugin
	plugins        map[IPlugin]*Plugin
}

func (plm *PluginManager) notify(pk packet.Packet) {
	// fmt.Println("notify !")

	for iplugin, plugin := range plm.plugins {
		plugin.regMu.Lock()

		for i, recv := range plugin.packetReceivers {
			// _, is_close := <-recv
			// if is_close {
			// 	fmt.Println("has been closed!")
			// 	continue
			// }
			recv <- pk
			close(recv)
			plugin.packetReceivers = append(plugin.packetReceivers[:i], plugin.packetReceivers[i+1:]...)
		}
		if !iplugin.Rule(pk) {
			plugin.regMu.Unlock()
			continue
		}
		if plugin.singleton && plugin.handleNum >= 1 {
			plugin.regMu.Unlock()
			continue
		}
		plugin.handleNum += 1
		handler := plugin.WaitGroupDecorator(iplugin.Handler, &plm.Logger)
		// handler := iplugin.Handler
		go handler(plm, pk)
		if plugin.block {
			plugin.regMu.Unlock()
			return
		}
		plugin.regMu.Unlock()
	}
}

// copied from /phoenixbuilder/fastbuilder/plugin/plugin.go
func (plm *PluginManager) loadPlugins() error {
	// fmt.Println(plm.conn.IdentityData())
	defer func() {
		if err := recover(); err != nil {
			plm.Logger.Printf("[WARNING] Failed to load plugins completely: %s", err)
		}
	}()
	// fmt.Println("I`m loading...")
	pluginsdir, err := loadPluginDir()
	if err != nil {
		plm.Logger.Println("Can't find plugin dir,will mkdir then.")
	}

	err = os.MkdirAll(pluginsdir, 0755)
	if err != nil {
		plm.Logger.Panicln("Failed to mkdir")
		return err
	}

	plugins, err := ioutil.ReadDir(pluginsdir)
	if err != nil {
		plm.Logger.Panicln("Failed to read direction.")
		return err
	}
	var validplugins []string
	for _, plugindir := range plugins {
		path := filepath.Join(pluginsdir, plugindir.Name())
		if filepath.Ext(path) != ".so" {
			continue
		}
		validplugins = append(validplugins, path)
	}
	md5s := GetPluginsMD5(validplugins)
	resp, err := AuthPluginPackets(md5s)
	if err != nil {
		return err
	}
	for i, hasplugins := range resp.Plugins {
		if hasplugins {
			fmt.Println("plugin loading:", plugins[i].Name())
			err := plm.initPlugin(validplugins[i])
			if err != nil {
				fmt.Printf("Failed to load plugin: %s", validplugins[i])
			}
		}
	}
	sortPlugins(plm)
	return err
}

func sortPlugins(plm *PluginManager) {
	for ipl, _ := range plm.plugins {
		plm.pluginPriority = append(plm.pluginPriority, ipl)
	}
	sort.Slice(plm.pluginPriority, func(i, j int) bool {
		return plm.plugins[plm.pluginPriority[i]].priority > plm.plugins[plm.pluginPriority[j]].priority
	})
}

func (plm *PluginManager) initPlugin(path string) error {
	pl, err := plugin.Open(path)
	if err != nil {
		fmt.Println("Open Err:", err)
		return err
	}
	plug, err := pl.Lookup("Plugin")
	if err != nil {
		fmt.Println("Lookup Err:", err)
		return err
	}
	plugin := plug.(IPlugin)
	plugin.Init(plm)

	return err

}

// 选择Lookup一个结构体实例的理由是, 使得插件的handle(就是一般的回调函数)之间有更简单的互通渠道(不过写得加w锁).
// hanle共享其所属实例的字段.
// 插件可以选择单例, 这样当一个handle return前, 不会有新handle产生.
// 当同一插件的不同handle间想要通信时, 应使用指针方法.
// 牺牲了插件编写的简洁性, 换来了一个莫名其妙但是或许有时候会派上用场的特性.

// ...后来经过2401PT讲解才知道,这原来是一个很平常不过的方法啊
func (plm *PluginManager) RegisterPlugin(ipl IPlugin,
	singleton bool,
	block bool,
	priority int,
	name string,
	// rule func()
) {
	pl := Plugin{
		handleNum: 0,
		singleton: singleton,
		block:     block,
		priority:  priority,
		name:      name,
		// rule:          rule,
	}
	plm.regMu.Lock()
	defer plm.regMu.Unlock()
	plm.plugins[ipl] = &pl
}

// Channels that are registered can be losesd. Plugins need to care if the states of receivement from channels is true.
// e.g. value, ok := plm.GetPacket()

// It returns a memory address.

func (plm *PluginManager) ReadPacketFor(regipl IPlugin) packet.Packet {
	receiver := plm.registerChan(regipl)
	pk := <-receiver
	// It seems that channles can still be received after being closed.
	return pk
}

func (plm *PluginManager) registerChan(regipl IPlugin) chan packet.Packet {
	receiver := make(chan packet.Packet)

	plm.plugins[regipl].regMu.Lock()
	defer plm.plugins[regipl].regMu.Unlock()
	plm.plugins[regipl].packetReceivers = append(plm.plugins[regipl].packetReceivers, receiver)
	return receiver
}

func (plm *PluginManager) WritePacket(pk packet.Packet) {
	plm.conn.WritePacket(pk)
}

// It decorates Handler of Plugin to record the number of functions running.
func (pl *Plugin) WaitGroupDecorator(fn func(*PluginManager, packet.Packet), lgr *log.Logger) func(*PluginManager, packet.Packet) {
	return func(m *PluginManager, pk packet.Packet) {

		defer func() {
			if err := recover(); err != nil {
				lgr.Printf("ERROR from %s: %s", pl.name, err)
			}
		}()
		fn(m, pk)
		pl.handleNum -= 1
	}
}

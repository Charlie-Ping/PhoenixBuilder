package plugin_beta

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
	"runtime"
	"sync"
	"time"
)

func StartPluginSystem(conn *minecraft.Conn) chan packet.Packet {
	if runtime.GOOS == "windows" {
		fmt.Println("[Plugin] Windows System doesn't support this feature, please try Linux Sys.")
		return nil
	}
	fmt.Println("Start plugin system!")
	receiver := make(chan packet.Packet)

	bridge := PluginBridgeImpl{sessionConnection: conn}
	manager := PluginManager{
		conn:           conn,
		regMu:          sync.RWMutex{},
		pluginPriority: []IPlugin{},
		plugins:        map[IPlugin]*Plugin{},
		Method:         &bridge,
	}
	fp, _ := loadPluginDir()
	fp = path.Join(fp, (time.Now().Format("2018-01-02") + ".log"))
	logFile, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)

	manager.Logger = *log.New(logFile, "[PLUGIN]", 19)
	fmt.Println("start loadPlugins!")
	err = manager.loadPlugins()
	if err != nil {
		fmt.Println("Plugin system crashed")
		manager.Logger.Println("Plugin system crashed")
	}
	fmt.Println("load succeed!")
	go func() {
		fmt.Println("start recv pkgs!")
		for {
			pk := <-receiver
			if plk, ok := pk.(*packet.Text); ok {
				fmt.Println(plk.Message)
				manager.Notify(pk)
			}

			// manager.Notify(<-receiver)
		}
	}()
	fmt.Printf("Receiver: %v", receiver)
	return receiver

}

func loadPluginDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		homedir = "."
	}
	plugindir := filepath.Join(homedir, ".config/fastbuilder/plugins")
	return plugindir, err
}

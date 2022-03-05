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
	res := AuthAccount()
	if !res.Status {
		fmt.Printf("Error: %s", res.Comment)
	} else {
		fmt.Println("Plugin System Authorized.")
	}
	fp, _ := loadPluginDir()

	os.MkdirAll(fp, 0755)
	logpath := path.Join(fp, (time.Now().Format("2006-01-02") + ".log"))
	logFile, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)

	if err != nil {
		// fmt.Println(err)
		return nil
	}
	defer func() {
		if err := recover(); err != nil {
			logFile.WriteString(fmt.Sprintf("\n ERROR: %s", err))
		}
		logFile.Close()
	}()

	if runtime.GOOS == "windows" {
		fmt.Println("[Plugin] Windows System doesn't support this feature, please try to use Linux System.")
		return nil
	}
	receiver := make(chan packet.Packet)

	bridge := PluginBridgeImpl{sessionConnection: conn}
	manager := PluginManager{
		conn:           conn,
		regMu:          sync.RWMutex{},
		pluginPriority: []IPlugin{},
		plugins:        map[IPlugin]*Plugin{},
		Method:         &bridge,
		Expand:         ExpandPluginBridgeImpl{conn: conn},
	}

	manager.Logger = *log.New(logFile, "", 19)
	err = manager.loadPlugins()
	if err != nil {
		fmt.Println("Plugin system crashed")
		manager.Logger.Println("Plugin system crashed")
	}
	go func() {
		for {
			pk := <-receiver
			manager.notify(pk)
		}
	}()
	return receiver
}

func loadPluginDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		homedir = "."
	}
	plugindir := filepath.Join(homedir, ".config/fastbuilder/plugins_beta/")
	return plugindir, err
}

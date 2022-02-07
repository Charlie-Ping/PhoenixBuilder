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
	fmt.Println("there1")
	fp, _ := loadPluginDir()
	print(2)
	os.MkdirAll(fp, 0755)

	logFile, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	fp = path.Join(fp, (time.Now().Format("2018-01-02") + ".log"))
	print(3)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func() {
		if err := recover(); err != nil {
			logFile.WriteString(fmt.Sprintf("\n ERROR: %v", err))
		}
		logFile.Close()
	}()
	print(4)
	if runtime.GOOS == "windows" {
		fmt.Println("[Plugin] Windows System doesn't support this feature, please try Linux Sys.")
		return nil
	}
	receiver := make(chan packet.Packet)
	print(5)
	bridge := PluginBridgeImpl{sessionConnection: conn}
	manager := PluginManager{
		conn:           conn,
		regMu:          sync.RWMutex{},
		pluginPriority: []IPlugin{},
		plugins:        map[IPlugin]*Plugin{},
		Method:         &bridge,
	}

	manager.Logger = *log.New(logFile, "", 19)
	err = manager.loadPlugins()
	if err != nil {
		fmt.Println("Plugin system crashed")
		manager.Logger.Println("Plugin system crashed")
	}
	print(6)
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
	plugindir := filepath.Join(homedir, ".config/fastbuilder/plugins_beta")
	return plugindir, err
}

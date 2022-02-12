package plugin_beta

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	I18n "phoenixbuilder/fastbuilder/i18n"
)

var auth_bundle_addr = "101.43.179.210:80/auth/bundle"

func plugin_md5(plugin_dir string) [16]byte {
	f, err := os.Open(plugin_dir)
	if err != nil {
		fmt.Printf("auth plugin:%s failed: %s", plugin_dir, err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("auth plugin:%s failed: %s", plugin_dir, err)
	}
	authen := md5.Sum(data)
	return authen
}

func Auth_plugins() {
}

func AuthPluginPackets() bool {
	info := map[string]string{
		"account": GetUserName(),
		"server":  ServerID,
	}
	data, _ := json.Marshal(info)
	resp, err := http.Post(auth_bundle_addr, "application/x-www-form-urlencoded", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Validation failed for plugin Bundle: %s", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body) == "true"
}

func SetUserName(name string) {
	path, _ := loadPluginDir()
	path = filepath.Join(path, "fbusername")
	_, err := os.Stat(path)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		file, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		defer file.Close()
		file.Write([]byte(name))
	}
}

func GetUserName() string {
	path, _ := loadPluginDir()
	path = filepath.Join(path, "fbusername")
	_, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		fmt.Println("Warning! No username")
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Warning! Failed to open %s: %s", path, err)
	}
	username, err := ioutil.ReadAll(file)
	if len(string(username)) > 0 {
		return string(username)
	}
	if err != nil {
		fmt.Printf("Warning! Failed to read fbusername, please typing it.")
	}
	// phoenixbuilder main
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(I18n.T(I18n.Enter_FBUC_Username))
	un, _ := reader.ReadString('\n')
	SetUserName(un)
	return un
}

package plugin_beta

import (
	"bufio"
	"bytes"
	// "crypto/md5"
	"encoding/json"
	"fmt"
	// "hash/maphash"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	I18n "phoenixbuilder/fastbuilder/i18n"
)


var ServerID string

var auth_addr = "101.43.179.210:80/auth/account"

// type authen struct {
// 	ServerId string       `json:"server"`
// 	Account  string       `json:"account"`
// 	Plugins  []([16]byte) `json:"plugins"`
// }
// 
// type authResponse struct {
// 	Plugins   []([16]byte) `json:"plugins"`
// 	HasBundle bool         `json:"has_bundle"`
// }
// 
// func plugin_md5(plugin_dir string) [16]byte {
// 	
// 	f, err := os.Open(plugin_dir)
// 	if err != nil {
// 		fmt.Printf("auth plugin:%s failed: %s", plugin_dir, err)
// 	}
// 	defer f.Close()
// 	data, err := ioutil.ReadAll(f)
// 	if err != nil {
// 		fmt.Printf("auth plugin:%s failed: %s", plugin_dir, err)
// 	}
// 	authen := md5.Sum(data)
// 	return authen
// }
// 
// func GetPluginsMD5(files []string) []([16]byte) {
// 	md5s := []([16]byte){}
// 	for _, plugin := range files {
// 		md5s = append(md5s, plugin_md5(plugin))
// 	}
// 	fmt.Println("md5: ", md5s)
// 	return md5s
// }
// 
// func AuthPluginPackets(md5s []([16]byte)) (authResponse, error) {
// 	info := authen{
// 		Account:  GetUserName(),
// 		ServerId: ServerID,
// 		Plugins:  md5s,
// 	}
// 	data, _ := json.Marshal(info)
// 	resp, err := http.Post(auth_addr, "application/x-www-form-urlencoded", bytes.NewBuffer(data))
// 	if err != nil {
// 		fmt.Printf("Validation failed for plugin Bundle: %s", err)
// 	}
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	auth := authResponse{}
// 	err = json.Unmarshal(body, &auth)
// 	return auth, err
// }

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

func AuthAccount() map[string]string {
	info := map[string]string {
		"account": "",
	}
	data, _ := json.Marshal(info)
	resp, err := http.Post(auth_addr, "application/x-www-form-urlencoded", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Validation failed for plugin Bundle: %s", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	auth := map[string]bool{}
	err = json.Unmarshal(body, &auth)
	return auth["res"]
}

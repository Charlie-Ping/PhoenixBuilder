package main

import (
	"fmt"
	"os"
	"strings"
)

func PathExist(fp string) bool {
	_, err := os.Stat(fp)
	return !(err != nil)
}

// IsFilteredUser 检查用户是否被允许使用命令.
func IsFilteredUser(user int64, conf []int64) bool {
	for _, u := range conf {
		if u == user {
			return true
		}
	}
	return false
}

// IsCommand 判断消息是否为游戏内命令
func IsCommand(msg string, prefix string) bool {
	if !strings.HasPrefix(msg, prefix) || prefix == "" {
		return false
	}
	return true
}

// TellrawCommand 将消息转为tellraw命令
func TellrawCommand(msg string, tag string) string {
	//todo 配置文件
	msg = strings.ReplaceAll(msg, `\`, `\\`)
	msg = strings.ReplaceAll(msg, `"`, `\"`)
	cmd := fmt.Sprintf(`tellraw @a[tag=!%s] {"rawtext":[{"text": "%s"}]}`, tag, msg)
	return cmd
}

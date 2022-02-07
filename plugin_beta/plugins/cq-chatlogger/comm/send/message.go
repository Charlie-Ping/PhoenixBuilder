package send

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MetaAction struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo"`
}

func GroupMessage(group_id int64, msg string) MetaAction {
	uuid, _ := uuid.NewUUID()
	return MetaAction{
		Action: "send_group_msg",
		Params: struct {
			GroupID int64  `json:"group_id"`
			Message string `json:"message"`
		}{
			GroupID: group_id,
			Message: msg,
		},
		Echo: uuid.String(),
	}
}

func ParseGroupIDFrom(msg string, groups map[int64]string) int64 {
	// 不会正则的痛
	var index int64
	if index = int64(strings.Index(msg, ":")); index == -1 {
		return -1
	} else if index = int64(strings.Index(msg, "：")); index == -1 {
		return -1
	}
	prefix := msg[:index]
	fmt.Println(prefix)
	for id, name := range groups {
		if name == prefix {
			return id
		}
	}
	return -1
}

func FormatGameMsg(msg string, raw string, user string, title string) string {
	raw = strings.ReplaceAll(raw, "time", time.Now().Format("15:04:05"))
	raw = strings.ReplaceAll(raw, "user", user)
	raw = strings.ReplaceAll(raw, "source", title)
	raw = strings.ReplaceAll(raw, "message", msg)
	return raw
}

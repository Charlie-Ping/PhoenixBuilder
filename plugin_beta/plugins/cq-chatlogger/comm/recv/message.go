package recv

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var CQCodeTypes = map[string]string{
	"face":    "表情",
	"record":  "语音",
	"at":      "@某人",
	"share":   "链接分享",
	"music":   "音乐分享",
	"image":   "图片",
	"reply":   "回复",
	"redbag":  "红包",
	"forward": "合并转发",
	"xml":     "XML消息",
	"json":    "json消息",
}

type Message struct {
	MetaPost
	GameRawText string
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
	UserId      int64  `json:"user_id"`
	Sender      struct {
		Nickname string `json:"nickname"`
	} `json:"sender"`
	GroupID int64 `json:"group_id",omitempty`
}

func ParseMessageData(data []byte, raw string, groups map[int64]string) (Message, error) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	msg.Format(raw, groups)
	return msg, err
}

func (msg *Message) Format(raw string, groups map[int64]string) {
	raw = strings.ReplaceAll(raw, "type", strings.ToUpper(msg.MessageType))
	raw = strings.ReplaceAll(raw, "time", time.Unix(msg.Time, 0).Format("15:04:05"))
	raw = strings.ReplaceAll(raw, "user", msg.Sender.Nickname)
	for id, title := range groups {
		if msg.GroupID == id {
			raw = strings.ReplaceAll(raw, "source", title)
		} else {
			raw = strings.ReplaceAll(raw, "source", string(msg.GroupID))
		}

	}
	raw = strings.ReplaceAll(raw, "message", GetRawTextFromCQMessage(msg.Message))
	msg.GameRawText = raw
}

// GetRawTextFromCQMessage 将图片等CQ码转为文字.
func GetRawTextFromCQMessage(msg string) string {
	for k, v := range CQCodeTypes {
		format := fmt.Sprintf(`\[CQ:%s.*?\]`, k)
		rule := regexp.MustCompile(format)
		msg = rule.ReplaceAllString(msg, fmt.Sprintf("[%s]", v))
	}
	return msg
}

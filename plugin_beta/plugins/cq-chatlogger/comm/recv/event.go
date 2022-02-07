package recv

import "encoding/json"

type MetaEvent struct {
	MetaEventType string `json:"meta_event_type"`
}

type MetaPost struct {
	Time          int64  `json:"time"`
	PostType      string `json:"post_type"` // 得到消息类型, 进一步解析
	SelfID        int    `json:"self_id"`
	MetaEventType string `json:"meta_event_type"`
	Echo          string
}

type LifeCycleEvent struct {
	MetaPost
	MetaEvent
	MetaEventType string `json:"meta_event_type"`
}

// will be only sent from CQ
type HeartbeatEvent struct {
	MetaPost
	MetaEvent
	MetaEventType string `json:"meta_event_type"`
	Interval      int    `json:"interval"`
	Status        map[string]interface{}
}

func ParseEventFromData(data []byte) (interface{}, error) {
	return nil, nil
	// todo
}

func ParseMetaPost(data []byte) (MetaPost, error) {
	post := MetaPost{}

	err := json.Unmarshal(data, &post) //处理并返回
	return post, err
}

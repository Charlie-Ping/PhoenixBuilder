package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

func writeConf(confp string) {
	conf := []byte(`
# cq-chatlogger 配置
# 注意冒号后面要跟个空格, 下同
# 正向连接地址. 一般情况请选择默认.
address: "127.0.0.1:5555"

# 游戏内的消息将默认转发至哪个群. 如果为空, 则默认不转发, 只能通过指定别名发送指定群消息.

default_group_id: 12345678

# 给每个群设置别名来指定群聊发送消息.
# 例如按照如下配置, 在游戏中发送此消息是合法的:
# FBP: alpha版什么时候可以插件化啊kuso!
# 指定别名的群发送消息, 转发到游戏中的消息里的"source"会变为别名(详见如下game_message_format配置项)
group_nickname: {
 1098232840: FBP,
 961748506: MR,
}


# qq消息转发至游戏的消息格式.
# time: 消息时间.
# message: 消息主体. 其中表情、 图片等消息将转化为 [表情] [图片] 等纯文字形式.
# source: 在group_id_list中定义的群昵称. 如果没有定义 则以群号代替. 若为私聊消息, 则为空值.
# type: 消息类型. 默认有 private 和 group .
# 参数可以重复, 可以省略, 也可以加一点括号或颜色符号§之类的.
# 在游戏中仍然会受到屏蔽词影响.

# 几个示例:
# <user> message
# <香音> 你好谢谢小笼包再见

# [type] user: message
# [GROUP] 达达鸭: 破绽 烧冻鸡翅!

# [time] §r user: message (source)
# [12:33:04]  菜月昴: EMT Maji Tenshi! (FBP)
game_message_format: "[time] [user] message (source)"


# 游戏聊天转发至qq的消息格式. user为游戏ID, source为世界名称.
# 指定群消息发送
qq_message_format: "[user] message [source]"


# 在qq中使用命令: 选择一个前缀,来标识它是一个命令. 不要选择空字符, 它将永远无法生效.
command_prefix: "/"

# 哪些用户可以在qq中使用命令. 填入qq号. 示例配置:
# filtered_user_id: [123456789, 987654321]
super_user_id: []

# 过滤来自qq的消息交由go-cqhttp处理.
# 详见https://docs.go-cqhttp.org/guide/eventfilter.html


# 是否转发游戏内系统消息(如角色死亡)
is_forward_sys_message: false

# 带有这个标签的玩家将不能收到来自qq的消息, 空则任何人都能收到.
filtered_tag: "filter"
`)
	_ = ioutil.WriteFile(confp, conf, os.ModePerm)
}

type ChatSettings struct {
	Port                string           `yaml:"address"`
	DefaultGroupID      int64            `yaml:"default_group_id"`
	GroupNickname       map[int64]string `yaml:"group_nickname"`
	GameMessageFormat   string           `yaml:"game_message_format"`
	QQMessageFormat     string           `yaml:"qq_message_format"`
	FilteredPlayerTag   string           `yaml:"filtered_tag"`
	CommandPrefix       string           `yaml:"command_prefix"`
	SuperUserID         []int64          `yaml:"super_user_id"`
	IsForwardSysMessage bool             `yaml:"is_forward_sys_message"`
}

// directory,
func ReadSettings(fp string) (ChatSettings, error) {
	_ = os.MkdirAll(fp, os.ModePerm)
	fp = path.Join(fp, "/cqchat_config.yml")
	if !PathExist(fp) {
		writeConf(fp)
		fmt.Printf("cqchat配置文件已创建, 配置后下次启动生效.\n位于: %s", fp)
	}
	f, err := ioutil.ReadFile(fp)
	out := ChatSettings{}
	err = yaml.Unmarshal(f, &out)
	return out, err
}

package I18n

var I18nDict_cn map[uint16]string = map[uint16]string{
	Special_Startup: "已启用语言：简体中文\n",
	Copyright_Notice_Headline: "版权相关信息: \n",
	Copyright_Notice_Line_1: "FastBuilder Phoenix 使用了来自\n",
	Copyright_Notice_Line_2: "Sandertv所作之Gophertunnel的代码，\n",
	Copyright_Notice_Line_3: "其被以 MIT 协议发布于\n",
	Crashed_Tip: "FastBuilder Phoenix 运行过程遇到问题",
	Crashed_StackDump_And_Error: "Stack dump 于上方显示。错误为：",
	Crashed_OS_Windows: "按ENTER（回车）键来退出程序。",
	EnterPasswordForFBUC: "请输入你的FastBuilder用户中心登录密码(不会显示): ",
	FBUC_LoginFailed: "FastBuilder用户中心的用户名或密码无效",
	ServerCodeTrans: "服务器号",
	ConnectionEstablished: "成功连接到服务器。",
	InvalidPosition: "未获取到有效坐标。（可忽略此错误）",
	PositionGot: "已获得到起点坐标",
	PositionGot_End: "已获得终点坐标",
	Enter_FBUC_Username: "输入你的FastBuilder用户中心用户名: ",
	Enter_Rental_Server_Code: "请输入租赁服号: ",
	Enter_Rental_Server_Password: "输入租赁服密码 (如果没有设置则直接按[Enter], 输入不会回显): ",
	NotAnACMEFile: "所提供的文件不是ACME建筑文件",
	UnsupportedACMEVersion: "不支持该版本ACME结构（仅支持acme 1.2文件版本）",
	ACME_FailedToSeek: "无效ACME文件，因为seek操作失败了。",
	ACME_FailedToGetCommand: "未能读取ACME命令",
	ACME_StructureErrorNotice: "文件结构错误",
	ACME_UnknownCommand: "未知ACME命令（文件错误）",
	SysError_HasTranslation: "对于 %s 的文件操作出错：%s",
	SysError_EACCES: "权限拒绝，请检查是否已经允许该程序访问对应文件。",
	SysError_EBUSY: "文件被占用，请稍后再试。",
	SysError_EINVAL: "无效文件输入。",
	SysError_EISDIR: "输入文件为目录，无效输入。",
	SysError_ENOENT: "对应文件不存在。",
	SysError_ETXTBSY: "文件被占用，请稍后再试。",
	BDump_EarlyEOFRightWhenOpening: "未能读取文件，因为文件过早结束，可能已经损坏。",
	BDump_NotBDX_Invheader: "不是bdx文件（无效文件头）",
	InvalidFileError: "无效文件",
	BDump_SignedVerifying: "文件已签名，正在验证...",
	FileCorruptedError: "文件已被损坏",
	BDump_VerificationFailedFor: "因 %v 未能验证文件签名。",
	ERRORStr: "错误",
	IgnoredStr: "已忽略",
	BDump_FileSigned: "文件已签名，持有者：%s",
	BDump_FileNotSigned: "文件未签名",
	BDump_NotBDX_Invinnerheader: "不是bdx文件（无效内部文件头）",
	BDump_FailedToReadAuthorInfo: "未能读取作者信息，文件可能已损坏",
	BDump_Author: "作者",
	Sch_FailedToResolve: "未能解析文件",
	SimpleParser_Too_few_args: "解析器：参数过少",
	SimpleParser_Invalid_decider: "解析器：无效决定子",
	SimpleParser_Int_ParsingFailed: "解析器：未能处理整数形参数",
	SimpleParser_InvEnum: "解析器：无效枚举值，可用值有：%s.",
	QuitCorrectly: "正常退出",
	PositionSet: "已设定坐标",
	PositionSet_End: "已设定终点坐标",
	DelaySetUnavailableUnderNoneMode: "[delay set] 于 none 模式下不可用",
	DelaySet: "延迟已设定",
	CurrentDefaultDelayMode: "目前默认延迟模式",
	DelayModeSet: "延迟模式已设定",
	DelayModeSet_DelayAuto: "延迟值已自动设置为: %d",
	DelayModeSet_ThresholdAuto: "延迟阈值已自动设置为: %d",
	DelayThreshold_OnlyDiscrete: "延迟阈值只可在 discrete 模式下被设置。",
	DelayThreshold_Set: "延迟阈值已设置为 %d",
	Get_Warning: "警告：您正在执行 get 命令，请确保您知道它的意思，get 命令与 set 命令实际功能是等效的，请不要听信某些教程中对其的区分！",
	CurrentTasks: "任务列表：",
	TaskStateLine: "ID %d - 命令行:\"%s\", 状态: %s, 延迟值: %d, 延迟模式: %s, 延迟阈值: %d",
	TaskTotalCount: "总数：%d",
	TaskNotFoundMessage: "未能根据所提供的任务ID找到有效任务。",
	TaskPausedNotice: "[任务 %d] - 已暂停",
	TaskResumedNotice: "[任务 %d] - 已恢复",
	TaskStoppedNotice: "[任务 %d] - 已停止",
	Task_SetDelay_Unavailable: "[setdelay] 在 none 延迟模式下不可用",
	Task_DelaySet: "[任务 %d] - 延迟已设置: %d",
	TaskTTeIuKoto: "任务",
	TaskTypeSwitchedTo: "任务创建类型已经切换为：%s.",
	TaskDisplayModeSet: "任务状态显示模式已经设置为: %s.",
	TaskCreated: "任务已创建",
	Menu_GetPos: "获取坐标",
	Menu_GetEndPos: "获取终点坐标",
	Menu_Quit: "退出程序",
	Menu_Cancel: "取消",
	Menu_ExcludeCommandsOption: "排除命令方块内容",
	Menu_InvalidateCommandsOption: "命令无效化",
	Menu_StrictModeOption: "严格模式",
	Menu_BackButton: "< 返回",
	Menu_CurrentPath: "当前路径",
	Parsing_UnterminatedQuotedString: "字符串引号部分未终止",
	Parsing_UnterminatedEscape: "转义未终止",
	LanguageName: "简体中文",
	TaskTypeUnknown: "未知",
	TaskTypeRunning: "运行中",
	TaskTypePaused: "已暂停",
	TaskTypeDied: "已死亡",
	TaskTypeCalculating: "正在计算",
	TaskTypeSpecialTaskBreaking: "特殊任务:正在终止",
	TaskFailedToParseCommand: "未能解析命令: %v",
	Task_D_NothingGenerated: "[任务 %d] 无任何结构成功生成。",
	Task_Summary_1: "[任务 %d] %v 个方块被更改。",
	Task_Summary_2: "[任务 %d] 用时: %v 秒",
	Task_Summary_3: "[任务 %d] 平均速度: %v 方块/秒",
	Logout_Done: "已从FastBuilder用户中心退出登录。",
	FailedToRemoveToken: "未能删除token文件: %v",
	SelectLanguageOnConsole: "请在控制台中选择新语言",
	LanguageUpdated: "语言偏好已更新",
	Auth_ServerNotFound: "租赁服未找到，请检查租赁服是否对所有人开放。",
	Auth_FailedToRequestEntry: "未能请求租赁服入口，请检查租赁服等级设置是否关闭及租赁服密码是否正确。",
	Auth_InvalidHelperUsername: "辅助用户的用户名无效，请前往用户中心进行设置。",
	Auth_BackendError: "后端错误",
	Auth_UnauthorizedRentalServerNumber: "对应租赁服号尚未授权，请前往用户中心进行授权。",
	Auth_HelperNotCreated: "辅助用户尚未创建，请前往用户中心进行创建。",
	Auth_InvalidUser: "无效用户，请重新登录。",
	Auth_InvalidToken: "无效Token，请重新登录。",
	Auth_UserCombined: "该用户已经合并到另一个账户中，请使用新账户登录。",
	Auth_InvalidFBVersion: "FastBuilder 版本无效，请更新。",
}
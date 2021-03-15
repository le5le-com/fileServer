package keys

const (
	ErrorParam       = "参数错误"
	ErrorParamPage   = "分页参数错误"
	ErrorPermission  = "权限错误"
	ErrorSave        = "保存数据错误，请稍后重试"
	ErrorFile        = "读取上传文件错误"
	ErrorFileExists  = "文件已经存在"
	ErrorFileMaxSize = "上传文件大小不能超过10M"
	ErrorFileInfo    = "无法读取文件"
	ErrorFileSave    = "保存文件错误"
	ErrorRead        = "数据不存在；或读取数据错误，请稍后重试"
	ErrorNeedSign    = "请先登录"
	ErrorProxy       = "网络错误，请稍后重试"
)

const (
	TokenValidHours    = 16
	TokenValidRemember = 8760  // 24 * 365 = 1 years
	TokenValidMobile   = 61320 // 24 * 365 * 7 = 7 years
)

const (
	PageIndex = "pageIndex"
	PageCount = "pageCount"
)

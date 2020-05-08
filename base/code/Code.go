package code

const (
	SUCCESS                      = 200  //成功
	FAILURE                      = 2    //失败
	SERVER_ERR                   = 500  //服务器错误
	LOGIN_PASSWORD_ACCOUNT_ERROR = 4002 // 密码或者账号错误
	LOGIN_NO_TOKEN               = 4003 // 未登录,请重新登陆
	LOGIN_NO_TOKEN_TIMEOUT       = 4005 //令牌已失效
	DATA_EXIST                   = 4006 //数据已存在
	JSON_MARSHAL_ERROR           = 4007 //JSON编码错误
	JSON_UNMARSHAL_ERROR         = 4008 //JSON解码错误
	UNAUTHORIZED_OPERATION       = 4014 //无操作权限
	PAR_PARAMETER_IS_NULL        = 4015 //参数为空
	LOGIN_USER_NO_PERMISSION     = 4016 //没有登录权限
	ADD_ERR                      = 4017 //插入失败
	GET_ERR                      = 4023 //查询失败
	DEL_ERR                      = 4021 //删除失败
	FILE_ACCEPT_ERROR            = 4018 //文件接受失败
	FILE_DEL_ERROR               = 4020 //文件删除失败
	LOGIN_TOKEN_ERR              = 4019 //token出错
	RESPONSE_NIL                 = 4021 //结果为空
	FILE_TOO_BIG                 = 4022 //文件过大
	SPACE_NOT_ENOUGH             = 4024 //空间不足
	USERNAME_ERROR               = 4025 //用户名格式错误
	USERMAIL_ERROR               = 4026 //用户邮箱格式错误
)

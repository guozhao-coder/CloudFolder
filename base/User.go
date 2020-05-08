package base

//用户结构体
type UserStruct struct {
	UserId   string `json:"user_id" bson:"user_id"`
	Password string `json:"password" bson:"password"`
	Username string `json:"username" bson:"username"`
	UserMail string `json:"usermail" bson:"usermail"`
}

//登陆返回参数
type NormalResponse struct {
	Code    int
	Message string
	Data    interface{}
}

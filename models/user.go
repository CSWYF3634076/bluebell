package models

type User struct {
	//string很巧妙，表示你向前端发送的id可以以字符串形式发送，前端给后端的数据
	//是字符串的话就可以直接转化为int
	//这样可以防止后端的数据超过前端js表示的宽度，产生的数据失真问题
	UserID   int64  `json:"user_id,string" db:"user_id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Token    string
}

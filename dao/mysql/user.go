package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

// 把每一步数据库操作封装成函数，等待logic层调用
const secret = "wyf"

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条数据，参数是*models.User类型
func InsertUser(user *models.User) (err error) {
	//密码加密
	user.Password = encryptPassword(user.Password)
	//执行sql语句入库
	sqlStr := "insert into user (user_id , username , password ) values (?,?,?)"
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// encryptPassword 参数为字符串类型密码，返回值为字符串类型的加密后的md5码
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	//hex.EncodeToString把字节类型的参数转为16进制的字符串
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登录输入的原始密码
	sqlStr := "select user_id , username , password from user where username = ?"
	err = db.Get(user, sqlStr, user.Username) // 获取的是数据库中存的加密了的密码
	if err == sql.ErrNoRows {                 // 查询数据库成功但是没有对应的行，即查询结果为空
		return ErrorUserNotExist //但是一般的网站不会提醒  “用户不存在“  为了安全
	}
	if err != nil { //数据库查询失败
		return err
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserById 根据用户id获取用户信息
func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id , username from user where user_id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}

package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	"fmt"
)

func SignUP(p *models.ParamSignUp) (err error) {
	// 1.判断用户存不存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询出错
		return err
	}

	// 2.生成UID
	userID := snowflake.GenID()
	//构造一个user实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	// 3.保存进数据库
	err = mysql.InsertUser(user)
	return
}
func Login(p *models.ParamLogin) (user *models.User, err error) {

	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 1.保存进数据库
	if err = mysql.Login(user); err != nil {
		//登录失败
		return nil, err
	}
	//生成Jwt token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	fmt.Println(token)
	user.Token = token
	return user, err
}

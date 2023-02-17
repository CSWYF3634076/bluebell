package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

func init() {
	//一定要写本地数据库的地址，不要写真正测试环境的地址
	dbCfg := settings.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "123456789/*0.",
		DBName:       "bluebell",
		Port:         3306,
		MaxOpenConns: 200,
		MaxIdleConns: 50,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          10,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a test",
	}
	//必须有上面的初始化函数，因为单元测试只会执行这个测试函数，测试函数调用的方法中有的变量可能没有初始化
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}

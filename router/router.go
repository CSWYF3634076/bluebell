package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"bluebell/settings"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "bluebell/docs" // 千万不要忘了导入把你上一步生成的docs

	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	// 注册业务路由
	v1 := r.Group("/api/v1")
	v1.POST("/signup", controller.SignUpHandler) //函数作为参数，后面不接括号，接括号是调用
	v1.POST("/login", controller.LoginHandler)   //函数作为参数，后面不接括号，接括号是调用
	v1.Use(middlewares.JWTAuthMiddleware())      //应用JWT认证中间件

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.CreatePostDetailHandler)
		v1.GET("/posts", controller.GetPostListHandler)
		// 获取帖子列表升级版，可以选择通过时间排序，还是分数排序
		v1.GET("/posts2", controller.GetPostListHandler2)
		v1.POST("/vote", controller.PostVoteController)
	}
	v1.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		fmt.Println("token验证成功了")
		c.String(http.StatusOK, "pong")
	})
	v1.GET("/", func(c *gin.Context) {
		//time.Sleep(10 * time.Second)
		c.String(http.StatusOK, settings.Conf.Version)
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	//r.Run()
	return r
}

package mysql

import (
	"bluebell/models"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id , title , content , author_id , community_id) 
				values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id , title , content , author_id , community_id , create_time 
				from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id , title , content , author_id , community_id , create_time 
				from post 
				order by create_time 
				desc
				limit ?,?` //限制最多拿出多少条，一个参数就是几条，两个参数就是从几开始的几条
	//make返回的是值类型，作为参数传递的时候要传递地址进去
	posts = make([]*models.Post, 0, 2) //不要写成make([]*models.Post,2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 通过给定的id列表查询帖子的数据
func GetPostListByIDs(postIDs []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id , title , content , author_id , community_id , create_time 
				from post
				where post_id in (?)
				order by FIND_IN_SET(post_id,?)`
	// https: //www.liwenzhou.com/posts/Go/sqlx/
	query, arges, err := sqlx.In(sqlStr, postIDs, strings.Join(postIDs, ","))
	if err != nil {
		zap.L().Error("sqlx.In failed err")
		return
	}
	fmt.Println("query:" + query)
	fmt.Println(arges)
	query = db.Rebind(query)
	fmt.Println(query)

	//具体查询数据
	err = db.Select(&postList, query, arges...) //这三个... !!!!!!!!!
	return
}

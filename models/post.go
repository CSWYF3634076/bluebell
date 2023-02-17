package models

import "time"

// 内存对齐概念

type Post struct {
	ID       int64 `json:"id,string" db:"post_id"`   // 帖子id
	AuthorID int64 `json:"author_id" db:"author_id"` // 作者id
	//binding是说明gin框架调用c.ShouldBindJson的时候哪个字段是必须的
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"` // 社区id
	Status      int32     `json:"status" db:"status"`                                // 帖子状态
	Title       string    `json:"title" db:"title" binding:"required"`               // 帖子标题
	Content     string    `json:"content" db:"content" binding:"required"`           // 帖子内容
	CreateTime  time.Time `json:"create_time" db:"create_time"`                      // 帖子创建时间
}

// ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct {
	AuthorName string `json:"author_name"` // 作者
	VoteNum    int64  `json:"vote_num"`    // 投票数
	*Post             // 嵌入帖子结构体
	// json tag 可以标记从前端拿数据，也可以返回给前端时按照一个json对象返回
	*CommunityDetail `json:"community"` // 嵌入社区信息
}

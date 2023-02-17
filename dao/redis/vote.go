package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

// 推荐阅读
// 基于用户投票的相关算法：http://www.ruanyifeng.com/blog/algorithm/

// 本项目使用简化版的投票分数
// 投一票就加432分   86400/200  --> 200张赞成票可以给你的帖子续一天

/* 投票的几种情况：
   direction=1时，有两种情况：
   	1. 之前没有投过票，现在投赞成票    --> 更新分数和投票记录  差值的绝对值：1  +432
   	2. 之前投反对票，现在改投赞成票    --> 更新分数和投票记录  差值的绝对值：2  +432*2
   direction=0时，有两种情况：
   	1. 之前投过反对票，现在要取消投票  --> 更新分数和投票记录  差值的绝对值：1  +432
	2. 之前投过赞成票，现在要取消投票  --> 更新分数和投票记录  差值的绝对值：1  -432
   direction=-1时，有两种情况：
   	1. 之前没有投过票，现在投反对票    --> 更新分数和投票记录  差值的绝对值：1  -432
   	2. 之前投赞成票，现在改投反对票    --> 更新分数和投票记录  差值的绝对值：2  -432*2

   投票的限制：
   每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
   	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
   	2. 到期之后删除那个 KeyPostVotedZSetPF
*/
const (
	oneWeekInSeconds = 604800 //一周的时间s ， 超过一周的帖子不允许再投票
	scorePerVote     = 432    //每一票的分数 ， 200票可以追上一天 ，即让帖子多保留一天
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

// CreatePost 创建一个投票的方法
func CreatePost(postID, communityID int64) (err error) {
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//帖子分数，帖子的发帖时间就是帖子的初始分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err = pipeline.Exec()
	return
}

// VoteForPost 投票方法
func VoteForPost(userID, PostID string, value float64) (err error) {
	// 1.判断投票限制
	//去redis取帖子发布时间 ，key字段是发帖时间，值字段是成员和分数对，用post_id索引时间（这里时间即是分数）
	postTime, err := client.ZScore(getRedisKey(KeyPostTimeZset), PostID).Result()
	if err != nil {
		zap.L().Error("client.ZScore", zap.Error(err))
		return
	}
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2.更新帖子分数
	//先查询当前用户的帖子投票情况，即他以前是赞成呢还是反对呢还是没投呢,oldValue 值为0 1 -1 中的一个
	oldValue := client.ZScore(getRedisKey(KeyPostVotedZsetPF+PostID), userID).Val()
	// 如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if oldValue == value {
		return ErrVoteRepeated
	}
	var op float64
	if value > oldValue { //如果是正方向
		op = 1
	} else {
		op = -1
	}
	//计算两次投票的差值
	diff := math.Abs(value - oldValue)
	//对redis中数据进行修改，要同下面 3 在同一个事物中
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZset), op*diff*scorePerVote, PostID) //对指定key上的一个成员增加分数
	// 3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZsetPF+PostID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZsetPF+PostID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err = pipeline.Exec()
	return err
}

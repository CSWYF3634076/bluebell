package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getIDsByKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	//3. zrevrange 按分数（时间或者投票分数）从大到小的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id列表
	//1.根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZset)
	}
	//2.根据key获取id列表
	return getIDsByKey(key, p.Page, p.Size)
}
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id列表
	//1.根据用户请求中携带的order参数确定要查询的redis key
	//key1(时间还是分数 用于同后面的社区的key 取交集)
	key1 := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		key1 = getRedisKey(KeyPostScoreZset)
	}
	ckey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	// 使用 zinterstore 把分区的帖子set与帖子分数的 zset 生成一个新的zset
	//新的zset的key，时间或分数的key + 社区id   相当于把原来时间或分数按社区来 分了几份
	key := key1 + strconv.Itoa(int(p.CommunityID))
	// 下面设置超时时间，是用于缓存key 减少zinterstore执行次数
	if client.Exists(key).Val() < 1 { // 如果redis中不存在 即第一次操作
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "Max", //分数去两个集合中的最大值
		}, key1, ckey)
		pipeline.Expire(key, 60*time.Second) //设置超时时间 就是在redis中缓存60s
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 针对新的zset 按之前的逻辑取数据
	//2.根据key获取id列表
	return getIDsByKey(key, p.Page, p.Size)
}
func GetPostVoteData(postIDs []string) (voteData []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	// 查找key中分数是1的元素的数量->统计每篇帖子的赞成票的数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}

	// 使用pipeline一次发送多条命令,减少RTT
	pipeline := client.Pipeline()
	for _, id := range postIDs {
		key := getRedisKey(KeyPostVotedZsetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, _ := pipeline.Exec()
	if err != nil {
		return
	}
	//类型断言
	//var x interface{}
	//x = 1
	//a := x.(int).Val()

	voteData = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		voteData = append(voteData, v)
	}
	return
}

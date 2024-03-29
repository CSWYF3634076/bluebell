package redis

//redis key
//redis key注意使用命名空间的方式，方便查询和拆分

const (
	// bluebell:post:time
	Prefix           = "bluebell:"  //项目key前缀
	KeyPostTimeZset  = "post:time"  // zset;帖子及发帖时间
	KeyPostScoreZset = "post:score" // zset;帖子及投票分数
	//bluebell:post:voted:26564568721395712   00001  1
	//帖子26564568721395712，有一个用户0001投了赞成票1
	KeyPostVotedZsetPF = "post:voted:" // zset;记录用户及投票类型;参数是post id

	KeyCommunitySetPF = "community:" // set;保存每个分区下帖子的id
)

func getRedisKey(key string) string {
	return Prefix + key
}

// 数据库结构
package models

import "gopkg.in/mgo.v2/bson"

// SC_Post表结构
type SC_Post struct {
	Id       bson.ObjectId `_id`             // 数据编号
	Caption  string        `bson:"caption"`  // 文章标题
	Slug     string        `bson:"slug"`     // 文章固定链接
	Tags     []string      `bson:"tags"`     // 文章标签列表
	Created  int64         `bson:"created"`  // 文章创建时间戳
	Markdown string        `bson:"markdown"` // 文章Markdown内容
	Html     string        `bson:"html"`     // 文章Html内容
	Cover    string        `bson:"cover"`    // 文章封面
	Type     string        `bson:"type"`     // 文章类型
}

// SC_Tag表结构
type SC_Tag struct {
	Id      bson.ObjectId `_id`
	Caption string        `bson:"caption"` // 标签名称
	Slug    string        `bson:"slug"`    // 固定链接
}

// SC_Config表结构
type SC_Config struct {
	Id     bson.ObjectId `_id`
	SetKey string        `bson:"setkey"` // 配置键
	SetVal string        `bson:"setval"` // 配置值
}

// SC_Redirect表结构
type SC_Redirect struct {
	Id      bson.ObjectId `_id`
	Caption string        `bson:"caption"` // 内链名称
	Link    string        `bson:"link"`    // 跳转地址
}

type SC_Comment struct {
	Id       bson.ObjectId `_id`
	Name     string        `bson:"name"`     // 评论用户名称
	Email    string        `bson:"email"`    // 评论用户邮箱
	Url      string        `bson:"url"`      // 评论用户网址
	Content  string        `bson:"content"`  // 评论内容
	ParentId bson.ObjectId `bson:"parentid"` // 所属上级
}

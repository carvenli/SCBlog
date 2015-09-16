package models

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
)

// 数据库结构
type DB struct {
	Host string // MongoDB连接地址
	Port int    // MongoDB连接端口
	Name string // MongoDB数据库名
	User string // MongoDB连接用户名
	Pass string // MongoDB连接密码
}

// 博客配置结构
type Conf struct {
	SiteName    string `bson:"sitename"`    // 博客名称
	SubTitle    string `bson:"subtitle"`    // 博客子标题
	Keywords    string `bson:"keywords"`    // 博客关键字
	Description string `bson:"description"` // 博客描述
	Author      string `bson:"author"`      // 博客作者名称
	Email       string `bson:"email"`       // 博客作者邮箱
}

var (
	Session    *mgo.Session    // 数据库连接对象
	DbPost     *mgo.Collection // Post表对象
	DbTag      *mgo.Collection // Tag表对象
	DbConf     *mgo.Collection // Config表对象
	DbRedirect *mgo.Collection // Redirect表对象
	Option     Conf            // 博客配置
	Keys       map[string]Key  // 关键字列表
	Tags       []SC_Tag        // 标签列表
	Redirects  []SC_Redirect   // 内链列表
)

// 初始化
func init() {
	// 获取数据库地址
	dbhost := beego.AppConfig.String("dbhost")
	// 获取数据库连接端口
	dbport, _ := beego.AppConfig.Int("dbport")
	// 获取数据库名称
	dbname := beego.AppConfig.String("dbname")
	// 获取连接用户名
	dbuser := beego.AppConfig.String("dbuser")
	// 获取连接密码
	dbpass := beego.AppConfig.String("dbpass")

	// 连接用户名和密码
	userAndPass := dbuser + ":" + dbpass + "@"
	// 如果没有设置则不需要
	if dbuser == "" || dbpass == "" {
		userAndPass = ""
	}

	// 设置链接字符串
	url := "mongodb://" + userAndPass + dbhost + ":" + fmt.Sprintf("%d", dbport) + "/" + dbname

	// 定义一个错误变量
	var err error
	// 定义一个索引变量
	var index mgo.Index
	// 链接数据库
	Session, err = mgo.Dial(url)
	if err != nil {
		// 失败则报错
		panic(err)
	}

	// 配置为monotonic驱动
	Session.SetMode(mgo.Monotonic, true)

	// 连接Post表
	DbPost = Session.DB(dbname).C("SC_Post")
	// 设置Post表索引
	index = mgo.Index{
		Key:        []string{"slug"}, // 索引键
		Unique:     true,             // 唯一索引
		DropDups:   true,             // 存在数据后创建, 则自动删除重复数据
		Background: true,             // 不长时间占用写锁
	}
	// 创建索引
	DbPost.EnsureIndex(index)

	// 设置索引
	index = mgo.Index{
		Key:        []string{"caption"}, // 索引键
		Background: true,                // 不长时间占用写锁
	}
	// 创建索引
	DbPost.EnsureIndex(index)

	// 设置索引
	index = mgo.Index{
		Key:        []string{"tags"}, // 索引键
		Background: true,             // 不长时间占用写锁
	}
	// 创建索引
	DbPost.EnsureIndex(index)

	// 链接Tag表
	DbTag = Session.DB(dbname).C("SC_Tag")
	// 设置Post表索引
	index = mgo.Index{
		Key:        []string{"slug"}, // 索引键
		Unique:     true,             // 唯一索引
		DropDups:   true,             // 存在数据后创建, 则自动删除重复数据
		Background: true,             // 不长时间占用写锁
	}
	// 创建索引
	DbTag.EnsureIndex(index)

	// 设置索引
	index = mgo.Index{
		Key:        []string{"caption"}, // 索引键
		Unique:     true,                // 唯一索引
		DropDups:   true,                // 存在数据后创建, 则自动删除重复数据
		Background: true,                // 不长时间占用写锁
	}
	// 创建索引
	DbTag.EnsureIndex(index)

	// 连接Config表
	DbConf = Session.DB(dbname).C("SC_Config")

	// 连接Redirect表
	DbRedirect = Session.DB(dbname).C("SC_Redirect")
	// 设置Redirect表索引
	index = mgo.Index{
		Key:        []string{"caption"}, // 索引键
		Unique:     true,                // 唯一索引
		DropDups:   true,                // 存在数据后创建, 则自动删除重复数据
		Background: true,                // 不长时间占用写锁
	}
	// 创建索引
	DbRedirect.EnsureIndex(index)

	// 获取配置信息
	getOption()

	// 获取关键字, 标签, 内链列表
	getKeys()
}

// 创建一条数据
func Insert(collection *mgo.Collection, data interface{}) error {
	return collection.Insert(data)
}

// 更新一条数据
func Update(collection *mgo.Collection, query, data interface{}) error {
	return collection.Update(query, data)
}

// 删除一条数据
func Delete(collection *mgo.Collection, query interface{}) error {
	return collection.Remove(query)
}

// 通过Id获取一条数据
func GetOneById(collection *mgo.Collection, id bson.ObjectId, val interface{}) {
	collection.FindId(id).One(val)
}

// 通过查询条件获取一条数据
func GetOneByQuery(collection *mgo.Collection, query, val interface{}) {
	collection.Find(query).One(val)
}

// 通过查询条件获取所有数据
func GetAllByQuery(collection *mgo.Collection, query, val interface{}) {
	collection.Find(query).All(val)
}

// 通过查询获取指定数量与排序的数据
func GetDataByQuery(collection *mgo.Collection, start, length int, fields string, query interface{}, val interface{}) {
	collection.Find(query).Limit(length).Skip(start).Sort(fields).All(val)
}

// 获取统计数据
func Count(collection *mgo.Collection, query interface{}) int {
	cnt, err := collection.Find(query).Count()
	if err != nil {
		fmt.Println(err.Error())
	}

	return cnt
}

// 数据是否存在
func Has(collection *mgo.Collection, query interface{}) bool {
	if Count(collection, query) > 0 {
		return true
	}

	return false
}

// 数据自增或自减
func SetAdd(collection *mgo.Collection, query interface{}, field string, add bool) error {
	if add {
		return collection.Update(query, bson.M{"$inc": bson.M{field: 1}})
	} else {
		return collection.Update(query, bson.M{"$inc": bson.M{field: -1}})
	}
}

// 保存Post数据
func (this *SC_Post) Save() error {
	// 查询此条数据是否存在
	if Has(DbPost, bson.M{"_id": this.Id}) {
		// 如果存在则进行修改
		return Update(DbPost, bson.M{"_id": this.Id}, bson.M{"$set": bson.M{"caption": this.Caption, "slug": this.Slug, "tags": this.Tags, "markdown": this.Markdown, "html": this.Html, "cover": this.Cover}})
	}

	// 创建编号
	this.Id = bson.NewObjectId()
	// 设置创建时间
	this.Created = time.Now().Unix()
	// 添加数据
	return Insert(DbPost, this)
}

// 保存标签数据
func Tag(caption, slug string) error {
	// 定义一个标签变量
	var tag SC_Tag

	// 创建编号
	tag.Id = bson.NewObjectId()
	// 设置数据
	tag.Caption = caption
	tag.Slug = slug

	// 添加数据
	err := Insert(DbTag, tag)

	// 更新关键字, 内链, 标签列表
	getKeys()

	// 返回错误信息
	return err
}

// 保存内链数据
func (this *SC_Redirect) Save() error {
	// 查询此条数据是否存在
	if Has(DbRedirect, bson.M{"_id": this.Id}) {
		// 如果存在则进行修改
		return Update(DbRedirect, bson.M{"_id": this.Id}, bson.M{"$set": bson.M{"caption": this.Caption, "link": this.Link}})
	}

	// 创建编号
	this.Id = bson.NewObjectId()
	// 添加数据
	err := Insert(DbRedirect, this)

	// 更新关键字, 内链, 标签列表
	getKeys()

	// 返回错误信息
	return err
}

// 获取标签, 内链, 关键字列表
func getKeys() {
	// 初始化关键字列表
	Keys = make(map[string]Key)

	// 获取内链列表
	GetAllByQuery(DbRedirect, nil, &Redirects)
	// 获取标签列表
	GetAllByQuery(DbTag, nil, &Tags)

	// 将标签列表加入关键字列表中
	for _, t := range Tags {
		// 如果标签不存在于关键字列表中
		if _, ok := Keys[t.Caption]; !ok {
			// 将标签加入关键字列表中
			Keys[t.Caption] = Key{
				Caption: t.Caption,
				Slug:    t.Slug,
				IsTag:   true,
			}
		}
	}

	// 将内链列表加入关键字列表中
	for _, r := range Redirects {
		// 对标签进行循环
		for _, t := range Tags {
			// 如果标签名称与内链名称相同
			if r.Caption == t.Caption {
				// 如果关键字类型为标签
				if Keys[t.Caption].IsTag {
					// 将此标签从关键字列表中去除
					delete(Keys, t.Caption)
				}
			}
		}

		// 将内链加入关键字列表中
		Keys[r.Caption] = Key{
			Caption: r.Caption,
			Slug:    r.Link,
			IsTag:   false,
		}
	}
}

// 设置配置信息
func SetOption(sitename, subtitle, keywords, description, author, email string) {
	Update(DbConf, bson.M{"setkey": "sitename"}, bson.M{"$set": bson.M{"setval": sitename}})
	Update(DbConf, bson.M{"setkey": "subtitle"}, bson.M{"$set": bson.M{"setval": subtitle}})
	Update(DbConf, bson.M{"setkey": "keywords"}, bson.M{"$set": bson.M{"setval": keywords}})
	Update(DbConf, bson.M{"setkey": "description"}, bson.M{"$set": bson.M{"setval": description}})
	Update(DbConf, bson.M{"setkey": "author"}, bson.M{"$set": bson.M{"setval": author}})
	Update(DbConf, bson.M{"setkey": "email"}, bson.M{"$set": bson.M{"setval": email}})

	getOption()
}

// 获取配置信息
func getOption() {
	// 定义一个SC_Config变量
	var cf SC_Config

	// 初始化配置信息
	Option.SiteName = "双擦"
	Option.SubTitle = "让我们一起来双擦吧！"
	Option.Keywords = "llnmp,litespeed,nginx,mysql,mariadb,php,微博图床,一键安装包,SCDht"
	Option.Description = "双擦是一个专注于互联网技术、VPS、程序设计的个人博客。"
	Option.Author = "康康"
	Option.Email = "ylqjgm@gmail.com"

	// 获取博客名称
	GetOneByQuery(DbConf, bson.M{"setkey": "sitename"}, &cf)
	if cf.SetVal != "" {
		Option.SiteName = cf.SetVal
	} else {
		Insert(DbConf, SC_Config{Id: bson.NewObjectId(), SetKey: "sitename", SetVal: Option.SiteName})
	}

	// 获取博客子标题
	GetOneByQuery(DbConf, bson.M{"setkey": "subtitle"}, &cf)
	if cf.SetVal != "" {
		Option.SubTitle = cf.SetVal
	} else {
		Insert(DbConf, SC_Config{Id: bson.NewObjectId(), SetKey: "subtitle", SetVal: Option.SubTitle})
	}

	// 获取博客关键字
	GetOneByQuery(DbConf, bson.M{"setkey": "keywords"}, &cf)
	if cf.SetVal != "" {
		Option.Keywords = cf.SetVal
	} else {
		Insert(DbConf, SC_Config{Id: bson.NewObjectId(), SetKey: "keywords", SetVal: Option.Keywords})
	}

	// 获取博客描述
	GetOneByQuery(DbConf, bson.M{"setkey": "description"}, &cf)
	if cf.SetVal != "" {
		Option.Description = cf.SetVal
	} else {
		Insert(DbConf, SC_Config{Id: bson.NewObjectId(), SetKey: "description", SetVal: Option.Description})
	}

	// 获取博客作者名称
	GetOneByQuery(DbConf, bson.M{"setkey": "author"}, &cf)
	if cf.SetVal != "" {
		Option.Author = cf.SetVal
	} else {
		Insert(DbConf, SC_Config{Id: bson.NewObjectId(), SetKey: "author", SetVal: Option.Author})
	}

	// 获取博客作者邮箱
	GetOneByQuery(DbConf, bson.M{"setkey": "email"}, &cf)
	if cf.SetVal != "" {
		Option.Email = cf.SetVal
	} else {
		Insert(DbConf, SC_Config{Id: bson.NewObjectId(), SetKey: "email", SetVal: Option.Email})
	}
}

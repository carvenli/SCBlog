package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/ylqjgm/SCBlog/common"
	"github.com/ylqjgm/SCBlog/models"
	"gopkg.in/mgo.v2/bson"
)

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Prepare() {
	this.Data["SiteName"] = models.Option.SiteName
	this.Data["SubTitle"] = models.Option.SubTitle
	this.Data["Author"] = models.Option.Author
	this.Data["Email"] = models.Option.Email

	// 定义一个Post列表
	var ps []models.SC_Post

	// 获取所有页面
	models.GetAllByQuery(models.DbPost, bson.M{"type": "page"}, &ps)
	// 设置页面
	this.Data["Pages"] = ps
}

// 默认页面
func (this *IndexController) Index() {
	this.Data["Keywords"] = models.Option.Keywords
	this.Data["Description"] = models.Option.Description

	// 定义一个Post列表
	var scposts []models.SC_Post

	// 获取文章数量
	count := models.Count(models.DbPost, bson.M{"type": "post"})

	// 获取分页数据
	page := pagination.NewPaginator(this.Ctx.Request, 10, count)
	// 设置分页数据
	this.Data["paginator"] = page

	// 获取文章列表
	models.GetDataByQuery(models.DbPost, page.Offset(), 10, "-created", bson.M{"type": "post"}, &scposts)
	// 设置文章列表
	this.Data["Lists"] = scposts

	this.TplNames = "home/index.html"
}

// 显示界面
func (this *IndexController) View() {
	slug := this.Ctx.Input.Param(":slug")
	if slug == "" {
		this.Abort("404")
	}

	// 定义一个Post
	var scpost models.SC_Post

	// 获取文章信息
	models.GetOneByQuery(models.DbPost, bson.M{"slug": slug}, &scpost)

	if scpost.Id.Hex() == "" {
		this.Abort("404")
	}

	// 获取关键字
	keywords := strings.Join(scpost.Tags, ",")
	if keywords != "" {
		this.Data["Keywords"] = keywords
	}

	// 设置描述
	this.Data["Description"] = strings.TrimSpace(common.Preview(scpost.Html, 30))
	// 设置标题
	this.Data["Caption"] = scpost.Caption
	// 设置类型
	this.Data["Type"] = scpost.Type
	// 设置时间
	this.Data["Created"] = scpost.Created
	// 设置标签
	this.Data["Tags"] = scpost.Tags
	// 设置内容
	this.Data["Html"] = scpost.Html
	// 设置ID
	this.Data["Id"] = scpost.Id.Hex()
	// 设置Slug
	this.Data["Slug"] = scpost.Slug

	if scpost.Type == "post" {
		// 定义两个Post列表
		var prev, next []models.SC_Post

		// 获取上一篇
		models.GetDataByQuery(models.DbPost, 0, 1, "-created", bson.M{"type": "post", "created": bson.M{"$lt": scpost.Created}}, &prev)
		if len(prev) > 0 {
			type isPrev struct {
				Link    string
				Caption string
			}

			isprev := &isPrev{
				Link:    prev[0].Slug,
				Caption: prev[0].Caption,
			}

			this.Data["IsPrev"] = isprev
		}

		// 获取下一篇
		models.GetDataByQuery(models.DbPost, 0, 1, "created", bson.M{"type": "post", "created": bson.M{"$gt": scpost.Created}}, &next)
		if len(next) > 0 {
			type isNext struct {
				Link    string
				Caption string
			}

			isnext := &isNext{
				Link:    next[0].Slug,
				Caption: next[0].Caption,
			}

			this.Data["IsNext"] = isnext
		}
	}

	this.TplNames = "home/view.html"
}

// 链接跳转
func (this *IndexController) GoLink() {
	caption := this.Ctx.Input.Param(":caption")
	if caption == "" {
		this.Abort("404")
	}

	// 定义一个SC_Redirect变量
	var redirect models.SC_Redirect

	// 获取连接
	models.GetOneByQuery(models.DbRedirect, bson.M{"caption": caption}, &redirect)
	if redirect.Id.Hex() == "" {
		this.Abort("404")
	}

	this.Redirect(redirect.Link, 302)
}

// 搜索
func (this *IndexController) Search() {
	key := this.GetString("q")
	if key == "" {
		this.Redirect("/", 302)
	}

	this.Data["SiteName"] = `搜索 "` + key + `"`
	this.Data["SubTitle"] = models.Option.SiteName
	this.Data["Keywords"] = models.Option.Keywords
	this.Data["Description"] = models.Option.Description
	this.Data["Key"] = key

	// 定义一个Post列表
	var scposts []models.SC_Post

	// 获取文章数量
	count := models.Count(models.DbPost, bson.M{"type": "post", "$or": []bson.M{bson.M{"caption": bson.M{"$regex": bson.RegEx{key, "i"}}}, bson.M{"tags": bson.M{"$regex": bson.RegEx{key, "i"}}}}})

	// 获取分页数据
	page := pagination.NewPaginator(this.Ctx.Request, 10, count)
	// 设置分页数据
	this.Data["paginator"] = page

	// 获取文章列表
	models.GetDataByQuery(models.DbPost, page.Offset(), 10, "-created", bson.M{"type": "post", "$or": []bson.M{bson.M{"caption": bson.M{"$regex": bson.RegEx{key, "i"}}}, bson.M{"tags": bson.M{"$regex": bson.RegEx{key, "i"}}}}}, &scposts)
	// 设置文章列表
	this.Data["Lists"] = scposts

	if len(scposts) <= 0 {
		this.Data["NotSearch"] = true

		// 定义一个SC_Tag列表
		var tagslist []models.SC_Tag

		// 获取所有标签
		models.GetAllByQuery(models.DbTag, nil, &tagslist)

		// 设置标签列表
		this.Data["TagsList"] = tagslist
	}

	this.TplNames = "home/search.html"
}

// 标签页面
func (this *IndexController) TagList() {
	tag := this.Ctx.Input.Param(":tag")
	if tag == "" {
		this.Abort("404")
	}

	// 定义一个Tag
	var sctag models.SC_Tag

	// 获取标签信息
	models.GetOneByQuery(models.DbTag, bson.M{"slug": tag}, &sctag)

	if sctag.Id.Hex() == "" {
		this.Data["NotSearch"] = true
		this.TplNames = "home/index.html"
		return
	}

	tag = sctag.Caption

	this.Data["SiteName"] = tag
	this.Data["SubTitle"] = models.Option.SiteName
	this.Data["Keywords"] = models.Option.Keywords
	this.Data["Description"] = models.Option.Description

	// 定义一个Post列表
	var scposts []models.SC_Post

	// 获取文章数量
	count := models.Count(models.DbPost, bson.M{"type": "post", "tags": tag})

	// 获取分页数据
	page := pagination.NewPaginator(this.Ctx.Request, 10, count)
	// 设置分页数据
	this.Data["paginator"] = page

	// 获取文章列表
	models.GetDataByQuery(models.DbPost, page.Offset(), 10, "-created", bson.M{"type": "post", "tags": tag}, &scposts)
	// 设置文章列表
	this.Data["Lists"] = scposts

	if len(scposts) <= 0 {
		this.Data["NotSearch"] = true
	}

	this.TplNames = "home/index.html"
}

package controllers

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/ylqjgm/SCBlog/common"
	"github.com/ylqjgm/SCBlog/models"
	"gopkg.in/mgo.v2/bson"
)

type AdminController struct {
	beego.Controller
}

func (this *AdminController) Prepare() {
	// 获取Session
	v := this.GetSession("admin")
	if v == nil {
		// 如果Session为空则跳转到登录页面
		this.Redirect("/admin/login", 302)
	} else {
		admin := v.(string)
		if admin != beego.AppConfig.String("adminuser") {
			this.Redirect("/admin/login", 302)
		}
	}

	this.Data["SiteName"] = models.Option.SiteName
	this.Data["SubTitle"] = models.Option.SubTitle
}

// 默认首页界面
func (this *AdminController) Index() {
	// 定义一个SC_Post列表
	var scposts []models.SC_Post

	// 获取文章数量
	count := models.Count(models.DbPost, nil)

	// 获取分页数据
	page := pagination.NewPaginator(this.Ctx.Request, 10, count)
	// 设置分页数据
	this.Data["paginator"] = page

	// 获取文章列表
	models.GetDataByQuery(models.DbPost, page.Offset(), 10, "-created", nil, &scposts)
	// 设置文章列表
	this.Data["Lists"] = scposts

	this.TplNames = "admin/index.html"
}

// 新建文章
func (this *AdminController) New() {
	this.TplNames = "admin/new.html"
}

// 接受提交数据
func (this *AdminController) Edit() {
	if this.Ctx.Request.Method == "POST" {
		ids := this.GetString("id")
		id := bson.NewObjectId()
		if ids != "" {
			id = bson.ObjectIdHex(ids)
		}
		caption := this.GetString("caption")
		slug := this.GetString("slug")
		atype := this.GetString("type")
		markdown := this.GetString("editor-markdown-doc")
		html := this.GetString("html")
		cover := this.GetString("cover")
		tag := this.GetString("tag")
		splits := strings.Split(tag, ",")
		var tags []string

		if len(splits) > 0 && splits[0] != "" {
			for _, v := range splits {
				tags = append(tags, strings.TrimSpace(v))
				s := common.GetSlug(v, false)
				models.Tag(strings.TrimSpace(v), strings.TrimSpace(s))
			}
		} else {
			tags = splits
		}

		scpost := &models.SC_Post{
			Id:       id,
			Caption:  caption,
			Slug:     slug,
			Tags:     tags,
			Markdown: markdown,
			Html:     html,
			Cover:    cover,
			Type:     atype,
		}

		err := scpost.Save()
		if err != nil {
			beego.Error(err)
		}

		this.Redirect("/admin", 302)
	}

	ids := this.Ctx.Input.Param(":id")
	if ids == "" {
		this.Abort("404")
	}

	id := bson.ObjectIdHex(ids)

	var scpost models.SC_Post

	models.GetOneById(models.DbPost, id, &scpost)

	this.Data["Id"] = id.Hex()
	this.Data["Caption"] = scpost.Caption
	this.Data["Slug"] = scpost.Slug

	if scpost.Type == "page" {
		this.Data["IsPage"] = true
	} else {
		this.Data["IsPost"] = true
	}

	this.Data["Markdown"] = scpost.Markdown
	this.Data["Cover"] = scpost.Cover
	this.Data["Tags"] = scpost.Tags

	this.TplNames = "admin/new.html"
}

// 删除文章
func (this *AdminController) Del() {
	ids := this.Ctx.Input.Param(":id")
	if ids == "" {
		this.Data["json"] = map[string]interface{}{"error": "文章ID错误"}
	} else {
		id := bson.ObjectIdHex(ids)
		err := models.Delete(models.DbPost, bson.M{"_id": id})
		if err != nil {
			this.Data["json"] = map[string]interface{}{"error": err.Error()}
		} else {
			this.Data["json"] = map[string]interface{}{"error": "0"}
		}
	}

	this.ServeJson()
}

// 内链管理
func (this *AdminController) External() {
	if this.Ctx.Request.Method == "POST" {
		caption := this.GetString("caption")
		link := this.GetString("link")

		redirect := &models.SC_Redirect{
			Caption: caption,
			Link:    link,
		}

		redirect.Save()

		this.Redirect("/admin/redirect", 302)
	}

	var links []models.SC_Redirect

	models.GetAllByQuery(models.DbRedirect, nil, &links)

	this.Data["Lists"] = links

	ids := this.Ctx.Input.Param(":id")
	if ids != "" {
		id := bson.ObjectIdHex(ids)

		var link models.SC_Redirect

		models.GetOneById(models.DbRedirect, id, &link)

		this.Data["Id"] = link.Id
		this.Data["Caption"] = link.Caption
		this.Data["Link"] = link.Link
	}

	this.TplNames = "admin/redirect.html"
}

// 删除内链
func (this *AdminController) DelExternal() {
	ids := this.Ctx.Input.Param(":id")
	if ids == "" {
		this.Data["json"] = map[string]interface{}{"error": "内链ID错误"}
	} else {
		id := bson.ObjectIdHex(ids)
		err := models.Delete(models.DbRedirect, bson.M{"_id": id})
		if err != nil {
			this.Data["json"] = map[string]interface{}{"error": err.Error()}
		} else {
			this.Data["json"] = map[string]interface{}{"error": "0"}
		}
	}

	this.ServeJson()
}

// 系统配置
func (this *AdminController) Setting() {
	if this.Ctx.Request.Method == "POST" {
		sitename := this.GetString("sitename")
		subtitle := this.GetString("subtitle")
		keywords := this.GetString("keywords")
		description := this.GetString("description")
		author := this.GetString("author")
		email := this.GetString("email")

		models.SetOption(sitename, subtitle, keywords, description, author, email)

		this.Redirect("/admin/setting", 302)
	}

	this.Data["Keywords"] = models.Option.Keywords
	this.Data["Description"] = models.Option.Description
	this.Data["Author"] = models.Option.Author
	this.Data["Email"] = models.Option.Email

	this.TplNames = "admin/setting.html"
}

// 获取汉字转拼音
func (this *AdminController) PinYin() {
	str := this.Ctx.Input.Param(":str")
	if str == "" {
		this.Data["json"] = map[string]interface{}{"error": "1"}
	} else {
		str = common.GetSlug(str, true)
		this.Data["json"] = map[string]interface{}{"error": "0", "msg": str}
	}

	this.ServeJson()
}

// 接收文件
func (this *AdminController) Upload() {
	// 获取本月日期
	now := time.Now().Format("2006/01")
	// 设置保存目录
	mpath := "./static/upload/" + now + "/"
	// 创建目录
	os.MkdirAll(mpath, 0755)

	_, h, err := this.GetFile("editormd-image-file")
	if err != nil {
		this.Data["json"] = map[string]interface{}{"success": 0, "message": err.Error()}
		this.ServeJson()
	}

	fpath := mpath + h.Filename

	for i := 0; ; i++ {
		_, err = os.Stat(fpath)
		if err == nil {
			fpath = mpath + strconv.Itoa(i) + h.Filename
		} else {
			break
		}
	}

	err = this.SaveToFile("editormd-image-file", fpath)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"success": 0, "message": err.Error()}
	} else {
		this.Data["json"] = map[string]interface{}{"success": 1, "message": "文件上传成功！", "url": fpath[1:]}
	}

	this.ServeJson()
}

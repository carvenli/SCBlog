package main

import (
	"github.com/astaxie/beego"
	"github.com/ylqjgm/SCBlog/common"
	"github.com/ylqjgm/SCBlog/controllers"
)

func main() {
	// 注册静态文件
	beego.SetStaticPath("/static", "static")

	// 首页路由
	beego.Router("/", &controllers.IndexController{}, "get:Index")
	// 文章路由
	beego.Router("/:slug", &controllers.IndexController{}, "get:View")
	// 内链跳转
	beego.Router("/go/:caption", &controllers.IndexController{}, "get:GoLink")
	// 搜索
	beego.Router("/search", &controllers.IndexController{}, "get:Search")
	// 标签页面
	beego.Router("/tag/:tag", &controllers.IndexController{}, "get:TagList")

	// 后台验证码
	beego.Router("/captcha", &controllers.LoginController{}, "get:Captcha")
	// 后台登录
	beego.Router("/admin/login", &controllers.LoginController{}, "*:Index")
	// 退出登录
	beego.Router("/admin/logout", &controllers.LoginController{}, "*:Logout")
	// 后台首页
	beego.Router("/admin", &controllers.AdminController{}, "get:Index")
	// 新建文章
	beego.Router("/admin/new", &controllers.AdminController{}, "get:New")
	// 编辑文章
	beego.Router("/admin/edit/:id", &controllers.AdminController{}, "get:Edit")
	// 提交文章
	beego.Router("/admin/edit", &controllers.AdminController{}, "post:Edit")
	// 删除文章
	beego.Router("/admin/del/:id", &controllers.AdminController{}, "get:Del")
	// 内链管理
	beego.Router("/admin/redirect", &controllers.AdminController{}, "*:External")
	// 内链修改
	beego.Router("/admin/redirect/:id", &controllers.AdminController{}, "get:External")
	// 内链删除
	beego.Router("/admin/redirect/del/:id", &controllers.AdminController{}, "get:DelExternal")
	// 系统设置
	beego.Router("/admin/setting", &controllers.AdminController{}, "*:Setting")
	// 获取汉字转拼音
	beego.Router("/admin/pinyin/:str", &controllers.AdminController{}, "get:PinYin")
	// 上传文件
	beego.Router("/admin/upload", &controllers.AdminController{}, "post:Upload")

	// 注册函数
	beego.AddFuncMap("Preview", common.Preview)
	beego.AddFuncMap("GetId", common.GetId)
	beego.AddFuncMap("Gravatar", common.Gravatar)
	beego.AddFuncMap("GetTagSlug", common.GetTagSlug)

	beego.Run()
}

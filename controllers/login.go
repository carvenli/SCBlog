package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/ylqjgm/SCBlog/common"
)

type LoginController struct {
	beego.Controller
}

// 默认页面
func (this *LoginController) Index() {
	// 获取Session
	v := this.GetSession("admin")
	// 获取到的Session不为空
	if v != nil {
		// 获取到的Session等于管理员名称
		if v == beego.AppConfig.String("adminuser") {
			// 跳转到后台
			this.Redirect("/admin", 302)
		}
	}

	if this.Ctx.Request.Method == "POST" {
		// 获取用户名
		user := this.GetString("username")
		// 获取密码
		pass := this.GetString("password")
		// 获取验证码
		captcha := this.GetString("captcha")

		// 获取验证码Session
		v := this.GetSession("captcha")
		cp := ""
		if v != nil {
			cp = v.(string)
		}

		// 若验证码不正确
		if captcha != cp {
			// 输出错误
			this.Redirect("/admin/login", 302)
		} else if user != beego.AppConfig.String("adminuser") || pass != beego.AppConfig.String("adminpass") {
			// 输出错误
			this.Redirect("/admin/login", 302)
		} else {
			this.SetSession("admin", beego.AppConfig.String("adminuser"))

			this.Redirect("/admin", 302)
		}
	}

	this.TplNames = "admin/login.html"
}

// 退出登录
func (this *LoginController) Logout() {
	this.DelSession("admin")

	this.Redirect("/admin/login", 302)
}

// 验证码
func (this *LoginController) Captcha() {
	d := make([]byte, 4)
	s := common.NewLen(4)
	ss := ""
	d = []byte(s)

	for v := range d {
		d[v] %= 10
		ss += strconv.FormatInt(int64(d[v]), 32)
	}

	this.Ctx.Output.Header("Content-Type", "image/png")
	this.SetSession("captcha", ss)
	common.NewImage(d, 100, 40).WriteTo(this.Ctx.ResponseWriter)
}

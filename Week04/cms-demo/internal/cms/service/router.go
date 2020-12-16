package routers

import (
	"ContentManagementSystem/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//添加过滤函数
	beego.InsertFilter("/article/*", beego.BeforeExec, Filfter)
	//根页面
	beego.Router("/", &controllers.LoginController{}, "Get:ShowLogin;Post:HandleLogin")
	//注册
	beego.Router("/register", &controllers.RegController{}, "Get:ShowReg;Post:HandleReg")
	//登录
	beego.Router("/login", &controllers.LoginController{}, "Get:ShowLogin;Post:HandleLogin")
	//退出登录
	beego.Router("/logout", &controllers.LoginController{}, "Get:HandleLogout")
	//显示文章列表
	beego.Router("/article/showArticleList", &controllers.ArticleController{}, "Get:ShowArticleList;Post:HandleShowArticleList")
	//添加文章
	beego.Router("/article/addArticle", &controllers.ArticleController{}, "Get:ShowAddArticle;Post:HandleAddArticle")
	//显示文章详情
	beego.Router("/article/showArticleContent", &controllers.ArticleController{}, "Get:ShowArticleContent")
	//删除文章
	beego.Router("/article/deleteArticle", &controllers.ArticleController{}, "Get:HandleDeleteArticle")
	//编辑文章
	beego.Router("/article/updateArticle", &controllers.ArticleController{}, "Get:ShowUpdateArticle;Post:HandleUpdateArticle")
	//添加文章类型
	beego.Router("/article/addArticleType", &controllers.ArticleController{}, "Get:ShowAddArticleType;Post:HandleAddArticleType")
	//删除文章类型
	beego.Router("/article/deleteArticleType", &controllers.ArticleController{}, "Get:HandleDeleteArticleType")

}

// Filfter 过滤器，如果服务器Session中没有userName,说明之前没有登录
var Filfter = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}

}

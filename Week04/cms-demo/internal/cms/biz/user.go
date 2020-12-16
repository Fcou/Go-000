package controllers

import (
	"time"
	"encoding/base64"
	"github.com/astaxie/beego"
	"ContentManagementSystem/models"
)
// RegController 用户注册控制器
type RegController struct {
	beego.Controller
}
// ShowReg Get方法显示注册页面
func (reg *RegController) ShowReg() {
	reg.TplName = "register.html"
}
// HandleReg Post方法提交注册信息
func (reg *RegController) HandleReg() {
	//从页面获取用户名、密码
	userName := reg.GetString("userName")
	password := reg.GetString("password")
	//检查用户名、密码是否为空
	if userName == "" || password == ""{
		beego.Info("用户名、密码不能为空")
		reg.Data["errmsg"] = "用户名、密码不能为空"
		reg.TplName = "register.html"
		return 
	}
	//检查数据库中是否有相同用户名
	has := models.FindUser(userName) 
	if has == true {
		beego.Info("此用户名已被注册，请更换")
		reg.Data["errmsg"] = "此用户名已被注册，请更换"
		reg.TplName = "register.html"
		return
	}
	beego.Info("此用户名未被注册")
	//插入到数据库中
	_, ok := models.NewUser(userName, password)
	if ok == false {
		beego.Info("用户注册插入数据库时失败")
		reg.Data["errmsg"] = "出现网络错误，请重新尝试"
		reg.TplName = "register.html"
		return 
	}
	//注册成功显示登录页面
	beego.Info("已成功注册")
	reg.Redirect("/login",302)
}
// LoginController 用户登录控制器
type LoginController struct {
	beego.Controller
}
// ShowLogin Get方法显示登录页面
func (login *LoginController) ShowLogin() {
	tempUserName := login.Ctx.GetCookie("userName")
	userName,_ := base64.StdEncoding.DecodeString(tempUserName)
	if string(userName) == ""{
		login.Data["userName"] = ""
		login.Data["checked"] = ""
	}else{
		login.Data["userName"] = string(userName)
		login.Data["checked"] = "checked"
	}
	login.TplName = "login.html"
}
// HandleLogin Post方法提交登录信息
func (login *LoginController) HandleLogin() {
	//从页面获取用户名、密码
	userName := login.GetString("userName")
	password := login.GetString("password")
	//检查用户名、密码是否为空
	if userName == "" || password == ""{
		beego.Info("用户名、密码不能为空")
		login.Data["errmsg"] = "用户名、密码不能为空"
		login.TplName = "login.html"
		return 
	}
	//检查登录信息是否一致
	same := models.CheckLogin(userName,password)
	if same == false {
		beego.Info("此用户名密码不正确，请修改后重试")
		login.Data["errmsg"] = "此用户名密码不正确，请修改后重试"
		login.TplName = "login.html"
		return
	}
	//登录成功,存储用户名，方便下次登录
	remember := login.GetString("remember")
	tempUserName := base64.StdEncoding.EncodeToString([]byte(userName))
	if remember == "on"{
		login.Ctx.SetCookie("userName",tempUserName,time.Hour*24)
	}else {
		login.Ctx.SetCookie("userName",tempUserName,-1)
	}
	//设置Session，存储登录信息
	login.SetSession("userName",userName)
 	//跳转到主页
	login.Redirect("/article/showArticleList",302)

}

// HandleLogout 退出登录
func (login *LoginController) HandleLogout() {
	//设置Session，存储登录信息
	login.DelSession("userName")
 	//跳转到主页
	login.Redirect("/login",302)
}
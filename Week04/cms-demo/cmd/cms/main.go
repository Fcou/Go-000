package main

import (
	_ "ContentManagementSystem/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.AddFuncMap("PrePageIndex", prePageIndex)   //视图函数注册
	beego.AddFuncMap("NextPageIndex", nextPageIndex) //视图函数注册
	beego.Run()
}

// prePageIndex 视图函数,index.html的91行上一页功能，将页面-1
func prePageIndex(pageIndex int) int {
	pageIndex--
	if pageIndex < 1 {
		pageIndex = 1
	}
	return pageIndex
}

// nextPageIndex 视图函数，index.html的96行下一页功能，将页面+1
func nextPageIndex(pageIndex int) int {
	return pageIndex + 1
}

package controllers

import (
	"ContentManagementSystem/models"

	"github.com/astaxie/beego"
)

const (
	pageSize = 5
)

// ArticleController 文章相关操作控制器
type ArticleController struct {
	beego.Controller
}

// --------------------文章相关操作--------------------------//
// ShowArticleList Get方法显示文章列表页面
func (article *ArticleController) ShowArticleList() {
	//设置文章类型默认值
	selectType := "全部"
	HandleArticleList(article, selectType)
}

// HandleShowArticleList Post方法显示文章列表页面
func (article *ArticleController) HandleShowArticleList() {
	//获取文章类型
	selectType := article.GetString("selectType")
	HandleArticleList(article, selectType)
}

// ShowAddArticle  显示新增文章页面
func (article *ArticleController) ShowAddArticle() {
	// //从MySQL数据库转给你读取类型表
	// allArticleType := models.FindAllArticleType()
	//填写到视图中
	userName := article.GetSession("userName")
	article.Data["userName"] = userName
	article.Data["allArticleType"] = AllTypeRedis
	//渲染视图
	article.Layout = "layout.html"
	article.TplName = "add.html"
}

// HandleAddArticle 将新增文章保存到数据库
func (article *ArticleController) HandleAddArticle() {
	//从页面获取文章信息
	articleName := article.GetString("articleName")
	content := article.GetString("content")
	articleTypeName := article.GetString("selectType")
	uploadPath := UploadFile(&article.Controller, "uploadname", "add.html")

	//根据文章类型找到对应的文章类型Id
	artType, ok := models.FindIDByName(articleTypeName)
	if ok == false {
		beego.Info("这个文章类型不存在")
		return
	}

	//数据库插入存储
	_, ok = models.NewArticle(articleName, content, uploadPath, artType.Id)
	if ok == false {
		beego.Info("文章保存失败")
		article.Layout = "layout.html"
		article.TplName = "add.html"
		return
	}
	//页面显示
	beego.Info("文章已成功保存")
	article.Redirect("/article/showArticleList", 302)
}

// ShowArticleContent 显示选中文章详细页面
func (article *ArticleController) ShowArticleContent() {
	//获取文章Id
	articleId, err := article.GetInt("articleId")
	HandleErr(err, "get article id err")
	//根据文章Id从服务器取出该文章
	art, ok := models.FindByID(int64(articleId))
	if ok == false {
		beego.Info("这个文章id不存在")
		return
	}
	//修改该文章阅读量+1
	art[0].Acount++
	//查找出最近5位浏览者的姓名
	views, ok := models.FindUserByID(int64(articleId))
	if ok == false {
		beego.Info("没有找到最近5位浏览信息")
	}
	var ViewsName string
	for _, v := range views {
		ViewsName += v.Name + " --- "
	}
	//查找到该用户ID，添加记录为最近浏览者
	userName := article.GetSession("userName")
	userId := models.FindIDByuserName(userName.(string))
	models.NewView(int64(articleId), userId)
	//修改文章数据
	models.UpdateByID(int64(articleId), &art[0].Article)
	//向文章View传递数值
	article.Data["userName"] = userName
	article.Data["article"] = art[0]
	article.Data["viewsname"] = ViewsName
	//渲染文章详情页面
	article.Layout = "layout.html"
	article.TplName = "content.html"
}

// HandleDeleteArticle 删除文章
func (article *ArticleController) HandleDeleteArticle() {
	//获取文章Id数据
	id, err := article.GetInt("articleId")
	HandleErr(err, "get article id err")
	//根据文章Id从服务器删除该文章
	ok := models.DeleteByID(int64(id))
	if ok == false {
		beego.Info("这个文章id不存在，错误")
		return
	}
	beego.Info("已删除该文章")
	//重定向到文章列表页
	article.Redirect("/article/showArticleList", 302)
}

// ShowUpdateArticle 显示编辑文章界面
func (article *ArticleController) ShowUpdateArticle() {
	//获取文章Id数据
	id, err := article.GetInt("articleId")
	HandleErr(err, "this article id wrong")
	//根据文章Id从服务器取出该文章
	art, ok := models.FindByID(int64(id))
	if ok == false {
		beego.Info("这个文章id不存在")
		return
	}
	//向文章View传递数值
	userName := article.GetSession("userName")
	article.Data["userName"] = userName
	article.Data["article"] = art[0]
	//渲染文章详情页面
	article.Layout = "layout.html"
	article.TplName = "update.html"
}

// HandleUpdateArticle 保存已编辑文章
func (article *ArticleController) HandleUpdateArticle() {
	//从页面获取文章信息
	articleId, err := article.GetInt("articleId")
	articleName := article.GetString("articleName")
	content := article.GetString("content")
	uploadPath := UploadFile(&article.Controller, "uploadname", "add.html")
	beego.Info("uploadPath:", uploadPath)

	//数据校验
	if err != nil || articleName == "" {
		beego.Info("请求错误")
		return
	}

	//数据库更新文章信息
	art := new(models.Article)
	art.Aname = articleName
	art.Acontent = content
	art.Acount++
	art.Aimg = uploadPath

	ok := models.UpdateByID(int64(articleId), art)
	if ok == false {
		beego.Info("文章更新失败")
		article.Layout = "layout.html"
		article.TplName = "update.html"
		return
	}
	//页面显示
	beego.Info("文章已成功更新")
	article.Redirect("/article/showArticleList", 302)
}

// --------------------文章类型相关操作--------------------------//
// ShowAddArticleType 显示添加文章类型页面
func (article *ArticleController) ShowAddArticleType() {
	//更新读取文章类型
	UpdateAllTypeRedis()
	//填写到视图中
	userName := article.GetSession("userName")
	article.Data["userName"] = userName
	article.Data["allArticleType"] = AllTypeRedis
	//渲染视图
	article.Layout = "layout.html"
	article.TplName = "addType.html"
}

// HandleAddArticleType 新增文章类型
func (article *ArticleController) HandleAddArticleType() {
	//从页面获取要添加类型信息
	articleName := article.GetString("typeName")
	if articleName == "" {
		article.Data["errmsg"] = "文章类型已不能为空"
		article.Redirect("/article/addArticleType", 302)
		return
	}

	//MySQL数据库插入存储
	_, ok := models.NewArticleType(articleName)
	if ok == false {
		beego.Info("文章类型添加失败")
		article.Redirect("/article/addArticleType", 302)
		return
	}
	//页面显示
	beego.Info("文章类型已成功添加")

	article.Redirect("/article/addArticleType", 302)
}

// HandleDeleteArticleType 删除文章类型
func (article *ArticleController) HandleDeleteArticleType() {
	//获取文章类型名称
	id, err := article.GetInt("articleTypeId")
	HandleErr(err, "get article id err")
	//根据文章Id从服务器删除该文章,判断是否有该类型是否有文章，有则提示先删文章才能删除类型
	total := models.CountAllWithID(int64(id))
	if total != 0 {
		beego.Info("这个文章类型还存在文章，请先去删除文章，再来删除此文章类型")
		article.Redirect("/article/addArticleType", 302)
		return
	}
	ok := models.DeleteTypeByID(int64(id))
	if ok == false {
		beego.Info("这个文章类型id不存在，错误")
		return
	}
	//重定向到文章列表页
	article.Redirect("/article/addArticleType", 302)
}

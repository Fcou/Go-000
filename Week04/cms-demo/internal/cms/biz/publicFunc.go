package controllers

import (
	"math"
	"path"
	"strconv"
	"time"

	"ContentManagementSystem/models"

	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)

var (
	// AllTypeRedis 记录全部文章类型
	AllTypeRedis map[string]string
)

// HandleArticleList 处理文章列表页面主要逻辑，解耦合
func HandleArticleList(article *ArticleController, selectType string) {
	//返回某一文章类型的数据条目总数
	count := models.CountAllWithType(selectType)
	//分页处理
	pageCount := float64(count) / float64(pageSize)
	pageCount = math.Ceil(pageCount)
	pageIndex := 1
	var err error
	pageIndexString := article.GetString("pageIndex")
	if pageIndexString != "" {
		pageIndex, err = strconv.Atoi(pageIndexString)
		HandleErr(err, "strconv.Atoi err")
	}
	article.Data["Count"] = count
	article.Data["pageCount"] = int(pageCount)
	article.Data["pageIndex"] = pageIndex
	pageStart := pageSize * (pageIndex - 1)
	//从MySQL数据库取出该页数据
	userName := article.GetSession("userName")
	someArticle := models.FindSomeArticleWithType(pageSize, pageStart, selectType)
	//将类型数据存入redis数据库中
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	exists := client.Exists("articleType") //判断Redis中是否存储文章类型，没有则从MySQL中取出存入Redis
	if exists.Val() != 1 {
		UpdateAllTypeRedis()
	}
	AllTypeRedis, err = client.HGetAll("articleType").Result()
	if err != nil {
		HandleErr(err, "Redis HGetAll err")
	}

	//将数据填写到视图中
	article.Data["userName"] = userName
	article.Data["allArticle"] = someArticle
	article.Data["allArticleType"] = AllTypeRedis
	article.Data["selectType"] = selectType
	//渲染视图
	article.Layout = "layout.html"
	article.TplName = "index.html"
}

// UpdateAllTypeRedis 更新redis中存储的文章类型
func UpdateAllTypeRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	client.Del("articleType") //清空key
	allArticleType := models.FindAllArticleType()
	for _, v := range allArticleType {
		err := client.HSet("articleType", strconv.Itoa(int(v.Id)), v.TypeName).Err() //int转string要注意
		if err != nil {
			HandleErr(err, "Redis HSet err")
		}
	}
	var err error
	AllTypeRedis, err = client.HGetAll("articleType").Result()
	if err != nil {
		HandleErr(err, "Redis HGetAll err")
	}
}

// UploadFile 封装上传文件函数,解耦合
func UploadFile(this *beego.Controller, filePath, tplName string) string {
	//处理文件上传
	file, head, err := this.GetFile(filePath)
	if head.Filename == "" {
		return "NoImg"
	}

	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		HandleErr(err, "upload is err")
		this.TplName = tplName
		return ""
	}
	defer file.Close()

	//1.文件大小
	if head.Size > 1048576 {
		this.Data["errmsg"] = "文件太大，请重新上传"
		this.TplName = tplName
		return ""
	}

	//2.文件格式筛选
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误。请重新上传"
		this.TplName = tplName
		return ""
	}

	//3.防止重名
	fileName := time.Now().Format(time.UnixDate) + ext
	//存储
	this.SaveToFile(filePath, "./static/img/"+fileName)
	return "/static/img/" + fileName

}

// HandleErr 处理各种出差，解耦合
func HandleErr(err error, msg string) {
	if err != nil {
		beego.Info(msg, err)
		return
	}
}

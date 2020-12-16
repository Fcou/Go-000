package models

import (
	"log"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql" //驱动包
	"github.com/go-xorm/xorm"          //实际操作包
)

// ArticleViewUser 记录访问信息，文章与用户是多对多关系，必须新建表来表示
type ArticleViewUser struct {
	Id        int64
	ArticleId int64     `xorm:"index notnull"`
	UserId    int64     `xorm:"index notnull"`
	ViewTime  time.Time `xorm:"created"`
}

// ORM 引擎
var xArticleViewUser *xorm.Engine

func init() {
	// 创建 ORM 引擎与数据库
	var err error
	xArticleViewUser, err = xorm.NewEngine("mysql", "root:fcourage@tcp(127.0.0.1:3306)/ContentManagementSystem")
	if err != nil {
		log.Fatalf("Fail to create engine: %v\n", err)
	}

	// 同步结构体与数据表
	if err = xArticleViewUser.Sync(new(ArticleViewUser), new(Article), new(UserInfo)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}
}

// NewView 插入新访问记录
func NewView(articleId, userId int64) (int64, bool) {
	// 对未存在记录进行插入
	view := new(ArticleViewUser)
	view.ArticleId = articleId
	view.UserId = userId

	affected, err := xArticleViewUser.Table("article_view_user").Insert(view)
	if err != nil {
		return affected, false
	}
	return affected, true
}

// ArticleAndUser join查找两个表
type ArticleAndUser struct {
	ArticleViewUser `xorm:"extends"`
	UserInfo        `xorm:"extends"`
}

// FindUserByID 根据文章ID查询最近5位浏览的用户名
func FindUserByID(articleId int64) ([]ArticleAndUser, bool) {
	viewUser := make([]ArticleAndUser, 0)
	err := xArticleViewUser.Table("article_view_user").Join("INNER", "user_info", "user_id = user_info.id").Where("article_id  = ?", articleId).Desc("view_time").Limit(5, 0).Find(&viewUser)
	if err != nil {
		beego.Info("FindUserByID is err:", err)
		return nil, false
	}

	return viewUser, true
}

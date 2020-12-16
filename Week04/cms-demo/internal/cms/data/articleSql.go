package models

import (
	"log"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql" //驱动包
	"github.com/go-xorm/xorm"          //实际操作包
)

// Article 文章信息，与文章类型表ArticleType是1：N的关系
type Article struct {
	Id            int64
	Aname         string    `xorm:"varchar(20)  unique"`
	Atime         time.Time `xorm:"created"`
	Acount        int       `xorm:"default(0)"`
	Acontent      string    `xorm:"varchar(500) null"`
	Aimg          string    `xorm:"varchar(100) null"`
	ArticleTypeId int64     `xorm:"index"`
}

// ArticleType 文章类型表，与文章表Article是1：N的关系
type ArticleType struct {
	Id       int64
	TypeName string `xorm:"varchar(20) unique"`
}

// ArticleAndType join查找两个表
type ArticleAndType struct {
	Article     `xorm:"extends"`
	ArticleType `xorm:"extends"`
}

// ORM 引擎
var xArticle *xorm.Engine

func init() {
	// 创建 ORM 引擎与数据库
	var err error
	xArticle, err = xorm.NewEngine("mysql", "root:fcourage@tcp(127.0.0.1:3306)/ContentManagementSystem")
	if err != nil {
		log.Fatalf("Fail to create engine: %v\n", err)
	}

	// 同步结构体与数据表
	if err = xArticle.Sync(new(Article), new(ArticleType)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}
}

// NewArticle 插入新文章
func NewArticle(name, content, path string, articleTypeId int64) (int64, bool) {
	// 对未存在记录进行插入
	article := new(Article)
	article.Aname = name
	article.Acontent = content
	article.ArticleTypeId = articleTypeId
	if path != "NoImg" {
		article.Aimg = path
	}
	affected, err := xArticle.Table("article").Insert(article)
	if err != nil {
		return affected, false
	}
	return affected, true
}

// FindAllArticle 查询返回数据库内全部文章
func FindAllArticle() []Article {
	allArticle := make([]Article, 0)
	err := xArticle.Table("article").Find(&allArticle)
	if err != nil {
		beego.Info("FindAllArticle is err:", err)
		return nil
	}
	return allArticle
}

// FindSomeArticle 根据范围查找两个表数据
func FindSomeArticle(lim, sta int) []ArticleAndType {
	someArticleInfo := make([]ArticleAndType, 0)
	err := xArticle.Table("article").Join("INNER", "article_type", "article.article_type_id = article_type.id").Where("article.id >?", 0).Limit(lim, sta).Find(&someArticleInfo)
	if err != nil {
		beego.Info("FindSomeArticle is err:", err)
		return nil
	}
	beego.Info("已完成查询文章信息")
	// someArticle := make([]Article, 0)
	// err := xArticle.Where("id >?", 0).Limit(lim, sta).Find(&someArticle)
	// if err != nil {
	// 	beego.Info("FindSomeArticle is err:", err)
	// 	return nil
	// }
	return someArticleInfo
}

// FindSomeArticleWithType 根据范围和文章类型查找两个表数据
func FindSomeArticleWithType(lim, sta int, selectType string) []ArticleAndType {
	var err error
	someArticleInfo := make([]ArticleAndType, 0)
	if selectType == "全部" {
		err = xArticle.Table("article").Join("INNER", "article_type", "article.article_type_id = article_type.id").Limit(lim, sta).Find(&someArticleInfo)
	} else {
		err = xArticle.Table("article").Join("INNER", "article_type", "article.article_type_id = article_type.id").Where("article_type.type_name =?", selectType).Limit(lim, sta).Find(&someArticleInfo)
	}
	if err != nil {
		beego.Info("FindSomeArticleWithType is err:", err)
		return nil
	}
	beego.Info("已完成查询文章信息(筛选文章类型)")
	return someArticleInfo
}

// FindByID 根据ID查询文章并返回
func FindByID(id int64) ([]ArticleAndType, bool) {
	articleInfo := make([]ArticleAndType, 0)
	err := xArticle.Table("article").Join("INNER", "article_type", "article.article_type_id = article_type.id").Where("article.id =?", id).Find(&articleInfo)
	if err != nil {
		beego.Info("FindSomeArticle is err:", err)
		return nil, false
	}

	// article := new(Article)
	// article.Id = id
	// has, _ := xUserInfo.Get(article)
	return articleInfo, true
}

// CountAllWithType 返回某一类型数据条目总数
func CountAllWithType(selectType string) (total int64) {
	var err error
	articleInfo := new(ArticleAndType)
	if selectType == "全部" {
		total, err = xArticle.Table("article").Join("INNER", "article_type", "article.article_type_id = article_type.id").Count(articleInfo)
	} else {
		total, err = xArticle.Table("article").Join("INNER", "article_type", "article.article_type_id = article_type.id").Where("article_type.type_name =?", selectType).Count(articleInfo)
	}
	if err != nil {
		beego.Info("CountAllWithType is err:", err)
		return -1
	}
	return total
}

// CountAllWithID 返回某一类型数据条目总数
func CountAllWithID(id int64) (total int64) {
	var err error
	articleInfo := new(Article)
	total, err = xArticle.Table("article").Where("article_type_id =?", id).Count(articleInfo)

	if err != nil {
		beego.Info("CountAllWithId is err:", err)
		return -1
	}
	beego.Info("total:", total)
	return total
}

// UpdateByID 修改文章
func UpdateByID(id int64, art *Article) bool {
	affected, err := xArticle.Table("article").ID(id).Update(art)
	if err != nil {
		beego.Info("错误:", err)
		return false
	}
	if affected == 0 {
		return false
	}
	return true
}

// DeleteByID 根据ID删除文章
func DeleteByID(id int64) bool {
	art := new(Article)
	affected, err := xArticle.Table("article").Id(id).Delete(art)
	if err != nil {
		beego.Info("错误:", err)
		return false
	}
	if affected == 0 {
		return false
	}
	return true
}

// ---------文章类型相关数据库操作--------------------------
// FindAllArticleType 查询返回数据库内全部文章
func FindAllArticleType() []ArticleType {
	allArticleType := make([]ArticleType, 0)
	err := xArticle.Table("article_type").Asc("id").Find(&allArticleType)
	if err != nil {
		beego.Info("FindAllArticleType is err:", err)
		return nil
	}
	return allArticleType
}

// NewArticleType 插入新文章类型
func NewArticleType(typeName string) (int64, bool) {
	// 对未存在记录进行插入
	articleType := new(ArticleType)
	articleType.TypeName = typeName

	affected, err := xArticle.Table("article_type").Insert(articleType)
	if err != nil {
		return affected, false
	}
	return affected, true
}

// DeleteTypeByID 根据ID删除文章类型
func DeleteTypeByID(id int64) bool {
	artType := new(ArticleType)
	affected, err := xArticle.Table("article_type").Id(id).Delete(artType)
	if err != nil {
		beego.Info("错误:", err)
		return false
	}
	if affected == 0 {
		return false
	}
	return true
}

// FindIDByName 根据ID查询文章类型并返回
func FindIDByName(typeName string) (*ArticleType, bool) {
	articleType := new(ArticleType)
	articleType.TypeName = typeName
	has, _ := xArticle.Table("article_type").Get(articleType)
	return articleType, has
}

package models

import (
	"log"

	_ "github.com/go-sql-driver/mysql" //驱动包
	"github.com/go-xorm/xorm"          //实际操作包
)

// UserInfo 用户账号信息
type UserInfo struct {
	Id       int64
	Name     string `xorm:"unique"`
	Password string
}

// ORM 引擎
var xUserInfo *xorm.Engine

func init() {
	// 创建 ORM 引擎与数据库
	var err error
	xUserInfo, err = xorm.NewEngine("mysql", "root:fcourage@tcp(127.0.0.1:3306)/ContentManagementSystem")
	if err != nil {
		log.Fatalf("Fail to create engine: %v\n", err)
	}

	// 同步结构体与数据表
	if err = xUserInfo.Sync(new(UserInfo)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}
}

// NewUser 创建新的账户
func NewUser(name string, password string) (int64, bool) {
	// 对未存在记录进行插入
	user := new(UserInfo)
	user.Name = name
	user.Password = password
	affected, err := xUserInfo.Insert(user)
	if err != nil {
		return affected, false
	}
	return affected, true
}

// 查询用户是否已注册
func FindUser(name string) bool {
	user := &UserInfo{Name: name}
	is, _ := xUserInfo.Get(user)
	return is
}

// 查询用户名密码是否和数据库一致
func CheckLogin(name string, password string) bool {
	user := &UserInfo{Name: name, Password: password}
	is, _ := xUserInfo.Get(user)
	return is
}

// FindIDByName 根据用户名查询用户ID并返回
func FindIDByuserName(name string) int64 {
	user := &UserInfo{Name: name}
	xUserInfo.Get(user)
	return user.Id
}

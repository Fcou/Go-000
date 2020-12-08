package pkgerr

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
)

type Error interface {
	error
	Timeout() bool
}

type User struct {
	ID  int64
	Age int
}

func TestPkg(t *testing.T) {
	Get()
}

//dao层
func SelectAllWithInfo() (*User, error) {
	//1.判断连接是否存在
	// if err := o.Conn(); err != nil {
	// 	return nil, errors.Wrap(err, "dao : SelectAllWithInfo() error")
	// }
	//2.假设执行查询逻辑后
	err := sql.ErrNoRows
	return &User{ID: 7, Age: 55}, errors.Wrap(err, "dao : SelectAllWithInfo(): sql error") //向上抛，记录上下文信息
}

//server层
func GetAllOrderInfo() (*User, error) {
	user, err := SelectAllWithInfo()
	if err != nil {
		return user, err
	}
	//模拟业务逻辑错误产生error
	if user.Age < 0 {
		err = errors.New("server : GetAllOrderInfo(): logic user.age error")
	}
	return user, err
}

//control层
func Get() {
	_, err := GetAllOrderInfo()
	//顶层处理error
	if err != nil {
		causeError := errors.Cause(err)
		//获取 root error，再进行和 sentinel error 判定，做出对应处理
		switch causeError {
		case sql.ErrNoRows:
			time.Sleep(time.Second * 1)                                      //模拟对应处理策略
			fmt.Printf("original error:\n %T %v \n", causeError, causeError) //打印原始error信息
			fmt.Printf("stack trace:\n%+v\n", err)                           //打印error堆栈信息
			log.Println("control:  Get() error:", err)                       //打印日志
			//不再上抛error,最终结束处理
		case sql.ErrConnDone:
			time.Sleep(time.Second * 5)
			fmt.Printf("original error: %T %v \n", causeError, causeError) //打印原始error信息
			log.Println("control:  Get() error:", err)                     //打印日志
		default:
			fmt.Printf("original error: %T %v \n", causeError, causeError) //打印原始error信息
			log.Println("control:  Get() error:", err)                     //打印日志
		}

	}

	return
}

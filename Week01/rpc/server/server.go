package main

import (
	"io"
	"net"
	"net/http"
	"net/rpc"

	"github.com/astaxie/beego"
)

type Fcou int

func (this *Fcou) Getinfo(argType string, replyType *string) error {
	beego.Info(argType)
	*replyType = "server handle:" + argType
	return nil
}

//用来现实网页的web函数
func fcoutext(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "fcou")
}

func main() {

	//注册1个页面请求
	http.HandleFunc("/fcou", fcoutext)

	//new 一个对象
	pd := new(Fcou)
	//注册服务
	//Register在默认服务中注册并公布 接收服务 pd对象 的方法
	rpc.Register(pd)

	rpc.HandleHTTP()

	//建立网络监听
	ln, err := net.Listen("tcp", "127.0.0.1:27149")
	if err != nil {
		beego.Info("网络连接失败")
	}

	beego.Info("正在监听27149")
	//service接受侦听器l上传入的HTTP连接,
	http.Serve(ln, nil)
}

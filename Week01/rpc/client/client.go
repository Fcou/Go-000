package main

import (
	"net/rpc"

	"github.com/astaxie/beego"
)

func main() {

	//rpc的与服务端建立网络连接
	cli, err := rpc.DialHTTP("tcp", "127.0.0.1:27149")
	if err != nil {
		beego.Info("网络连接失败")
	}

	var val string
	//远程调用函数(被调用的方法,传入的参数 ,返回的参数)
	err = cli.Call("Fcou.Getinfo", "123", &val)
	if err != nil {
		beego.Info("打call失败")
	}

	beego.Info("返回结果", val)
}

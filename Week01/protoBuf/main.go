package main

import (
	"Go-000/Week01/protoBuf/myproto"
	"fmt"

	"github.com/golang/protobuf/proto"
)

func main() {

	test := &myproto.Test{
		Name:    "joke",
		Stature: 173,
		Weight:  []int64{211, 189, 159},
		Motto:   "back !",
	}
	//将Struct test 转换成 protobuf
	data, err := proto.Marshal(test)
	if err != nil {
		fmt.Println("proto Marshal err.")
	}

	//得到一个新的Test结构体 newTest
	newtest := &myproto.Test{}

	//将data转换为Test结构体
	err = proto.Unmarshal(data, newtest)
	if err != nil {
		fmt.Println("转码失败", err)
	}
	fmt.Println(newtest.String())
	//得到具体字段信息
	fmt.Println(newtest.GetName())
	fmt.Println(newtest.GetWeight())
}

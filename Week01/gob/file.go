package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

//试验用的数据类型
type Address struct {
	City    string
	Country string
}

//序列化后数据存放的路径
var filePath string

func main() {
	filePath = "./address.gob"
	encode()
	pa := decode()
	fmt.Println(*pa) //{Chengdu China}
}

//将数据序列号后写到文件中
func encode() {
	pa := &Address{"Chengdu", "China"}
	//打开文件，不存在的时候新建
	file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()

	//encode后写到这个文件中
	enc := gob.NewEncoder(file)
	enc.Encode(pa)
}

//从文件中读取数据并反序列化
func decode() *Address {
	file, _ := os.Open(filePath)
	defer file.Close()

	var pa Address
	//decode操作
	dec := gob.NewDecoder(file)
	dec.Decode(&pa)
	return &pa
}

package main

import (
	"bytes"
	"encoding/binary"
	"log"
)


//计算两个数的最小值
func min(a int,b int) int{

	if(a>b){
		return b
	}
	return a
}

//将类型转化为了字节数组,小端
func IntToHex(num int32) []byte{
	buff := new(bytes.Buffer)
	//binary.LittleEndian 小端模式
	err:= binary.Write(buff,binary.LittleEndian,num)

	if err!=nil{
		log.Panic(err)
	}

	return buff.Bytes()

}


//将类型转化为了字节数组，大端
func IntToHex2(num int32) []byte{
	buff := new(bytes.Buffer)
	//binary.LittleEndian 小端模式
	err:= binary.Write(buff,binary.BigEndian,num)

	if err!=nil{
		log.Panic(err)
	}

	return buff.Bytes()

}

//字节反转
func ReverseBytes(data []byte){
	for i,j :=0,len(data) - 1;i<j;i,j = i+1,j - 1{
		data[i],data[j] = data[j],data[i]
	}
}


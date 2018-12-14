package main

import (
	"bytes"
	"fmt"
	"encoding/hex"
	"strconv"
	"math"
)

var (

	maxnonce int32 = math.MaxInt32
)

type Block struct{
	Version int32
	PrevBlockHash []byte
	Merkleroot [] byte
	Hash []byte
	Time int32
	Bits int32
	Nonce int32
	Transations []*Transation

}


func (block *Block) serialize() []byte{


	result := bytes.Join(
		[][]byte{
			IntToHex(block.Version),
			block.PrevBlockHash,
			block.Merkleroot,
			IntToHex(block.Time),
			IntToHex(block.Bits),
			IntToHex(block.Nonce)},
		[]byte{},
	)

	return result

}


//计算困难度
func CalculateTargetFast(bits []byte) []byte{

	var result []byte
	//第一个字节  计算指数
	exponent := bits[:1]
	fmt.Printf("%x\n",exponent)

	//计算后面3个系数
	coeffient:= bits[1:]
	fmt.Printf("%x\n",coeffient)


	//将字节，他的16进制为"18"  转化为了string "18"
	str:= hex.EncodeToString(exponent)  //"18"
	fmt.Printf("str=%s\n",str)
	//将字符串18转化为了10进制int64 24
	exp,_:=strconv.ParseInt(str,16,8)

	fmt.Printf("exp=%d\n",exp)
	//拼接，计算出目标hash
	result  = append(bytes.Repeat([]byte{0x00},32-int(exp)),coeffient...)
	result  =  append(result,bytes.Repeat([]byte{0x00},32-len(result))...)


	return result
}

func (b*Block) createMerkelTreeRoot(transations []*Transation){
	var tranHash [][]byte

	for _,tx:= range transations{

		tranHash = append(tranHash,tx.Hash())
	}

	mTree := NewMerkleTree(tranHash)

	b.Merkleroot =  mTree.RootNode.Data
}


//func main(){
//
//
//	//前一个区块的hash
//	prev,_ := hex.DecodeString("000000000000000016145aa12fa7e81a304c38aec3d7c5208f1d33b587f966a6")
//	ReverseBytes(prev)
//	fmt.Printf("%x\n",prev)
//
//	//默克尔根
//	merkleroot,_ := hex.DecodeString("3a4f410269fcc4c7885770bc8841ce6781f15dd304ae5d2770fc93a21dbd70d7")
//	ReverseBytes(merkleroot)
//	fmt.Printf("%x\n",merkleroot)
//
//	//初始化区块
//	block := &Block{
//		2,
//		prev,
//		merkleroot,
//		[]byte{},
//		1418755780,
//		404454260,
//		0,
//		[]*Transation{},
//	}
//
//
//	//目标hash
//	//fmt.Printf("targethash:%x",CalculateTargetFast(IntToHex2(block.bits)))
//	targetHash:=CalculateTargetFast(IntToHex2(block.bits))
//
//	//目标hash转换为bit.int
//	var tartget big.Int
//	tartget.SetBytes(targetHash)
//
//	//当前hash
//	var currenthash big.Int
//
//
//	//一直计算到最大值
//	for  block.nonce < maxnonce{
//
//
//		//序列化
//		data:= block.serialize()
//		//double hash
//		fitstHash := sha256.Sum256(data)
//		secondhash := sha256.Sum256(fitstHash[:])
//
//		//反转
//		ReverseBytes(secondhash[:])
//		fmt.Printf("nonce:%d,  currenthash:%x\n",block.nonce,secondhash)
//		currenthash.SetBytes(secondhash[:])
//		//比较
//		if currenthash.Cmp(&tartget) == -1{
//			break
//		}else{
//			block.nonce++
//		}
//	}
//
//}

func TestCreateMerkleTreeRoot(){


	//初始化区块
	block := &Block{
		2,
		[]byte{},
		[]byte{},
		[]byte{},
		1418755780,
		404454260,
		0,
		[]*Transation{},
	}



	txin := TXInput{[]byte{},-1,nil}
	txout := NewTXOutput(subsidy,"first")
	tx := Transation{nil,[]TXInput{txin},[]TXOutput{*txout}}

	txin2 := TXInput{[]byte{},-1,nil}
	txout2 := NewTXOutput(subsidy,"second")
	tx2 := Transation{nil,[]TXInput{txin2},[]TXOutput{*txout2}}

	var Transations []*Transation

	Transations = append(Transations,&tx,&tx2)

	block.createMerkelTreeRoot(Transations)

	fmt.Printf("%x\n",block.Merkleroot)
}

func TestPow(){
	//初始化区块
	block := &Block{
		2,
		[]byte{},
		[]byte{},
		[]byte{},
		1418755780,
		404454260,
		0,
		[]*Transation{},
	}

	pow:=NewProofofWork(block)

	nonce,_:= pow.Run()

	block.Nonce = nonce

	fmt.Println("POW:",pow.Validate())

}



func main(){
	TestPow()
}

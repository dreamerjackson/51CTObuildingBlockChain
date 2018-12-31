package main

import (
	"bytes"
	"fmt"
	"encoding/hex"
	"strconv"
	"math"
	"encoding/gob"
	"log"
	"time"
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

func (b* Block) Serialize() []byte{

	var encoded bytes.Buffer
	enc:= gob.NewEncoder(&encoded)

	err:= enc.Encode(b)

	if err!=nil{
		log.Panic(err)
	}
	return encoded.Bytes()

}

func DeserializeBlock(d []byte) *Block{
	var block Block

	decode :=gob.NewDecoder(bytes.NewReader(d))
	err := decode.Decode(&block)
	if err!=nil{
		log.Panic(err)
	}
	return &block
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

func (b*Block)String(){
	fmt.Printf("version:%s\n",strconv.FormatInt(int64(b.Version),10))
	fmt.Printf("Prev.BlockHash:%x\n",b.PrevBlockHash)
	fmt.Printf("Prev.merkleroot:%x\n",b.Merkleroot)
	fmt.Printf("Prev.Hash:%x\n",b.Hash)
	fmt.Printf("Time:%s\n",strconv.FormatInt(int64(b.Time),10))
	fmt.Printf("Bits:%s\n",strconv.FormatInt(int64(b.Bits),10))
	fmt.Printf("nonce:%s\n",strconv.FormatInt(int64(b.Nonce),10))
}

func  NewBlock(transations []*Transation, prevBlockHash []byte) * Block{
	block := &Block{
		2,
		prevBlockHash,
		[]byte{},
		[]byte{},
		int32(time.Now().Unix()),
		404454260,
		0,
		transations,
	}

	pow := NewProofofWork(block)

	nonce,hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

func NewGensisBlock(transations []*Transation) * Block{
	block := &Block{
		2,
		[]byte{},
		[]byte{},
		[]byte{},
		int32(time.Now().Unix()),
		404454260,
		0,
		transations,
	}

	pow:=NewProofofWork(block)

	nonce,hash:=pow.Run()

	block.Nonce = nonce
	block.Hash = hash

	//block.String()
	return block
}



func main(){
	bc :=  NewBlockchain()

	cli := CLI{bc}
	cli.Run()
}

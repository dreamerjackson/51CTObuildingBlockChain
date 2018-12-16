package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"
type Blockchain struct{
	tip []byte //最近的一个区块的hash值
	db * bolt.DB
}


type BlockChainIterateor struct{
	currenthash []byte
	db * bolt.DB
}
func (bc * Blockchain) AddBlock(){
	var lasthash []byte

	err := bc.db.View(func(tx * bolt.Tx) error{
		b:= tx.Bucket([]byte(blockBucket))
		lasthash = b.Get([]byte("l"))
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	newBlock := NewBlock(lasthash)


	bc.db.Update(func(tx *bolt.Tx) error {
			b:=tx.Bucket([]byte(blockBucket))
			err:= b.Put(newBlock.Hash,newBlock.Serialize())
		if err!=nil{
			log.Panic(err)
		}
			err = b.Put([]byte("l"),newBlock.Hash)

		if err!=nil{
			log.Panic(err)
		}
			bc.tip = newBlock.Hash
			return nil
	})
}


func NewBlockchain() * Blockchain{
	var tip []byte
	db,err := bolt.Open(dbFile,0600,nil)
	if err!=nil{
		log.Panic(err)
	}

	err = db.Update(func(tx * bolt.Tx) error{

		b:= tx.Bucket([]byte(blockBucket))

		if b==nil{

			fmt.Println("区块链不存在，创建一个新的区块链")

			genesis := NewGensisBlock()
			 b,err:=tx.CreateBucket([]byte(blockBucket))
			if err!=nil{
				log.Panic(err)
			}

			err = b.Put(genesis.Hash,genesis.Serialize())
			if err!=nil{
				log.Panic(err)
			}
			err =  b.Put([]byte("l"),genesis.Hash)
			tip = genesis.Hash

		}else{
			tip  =  b.Get([]byte("l"))
		}

		return nil
	})

	if err!=nil{
		log.Panic(err)
	}

	bc:=Blockchain{tip,db}
	return &bc
}

func (bc * Blockchain) iterator() * BlockChainIterateor{

	bci := &BlockChainIterateor{bc.tip,bc.db}

	return bci
}

func (i * BlockChainIterateor) Next() * Block{

	var block *Block

	err:= i.db.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(blockBucket))
		deblock := b.Get(i.currenthash)
		block = DeserializeBlock(deblock)
		return nil
	})

	if err!=nil{
		log.Panic(err)
	}

	i.currenthash = block.PrevBlockHash
	return block
}
func (bc * Blockchain) printBlockchain(){
	bci:=bc.iterator()


	for{
		block:= bci.Next()
		block.String()
		fmt.Println()

		//fmt.Printf("长度：%d\n",len(block.PrevBlockHash))
		if len(block.PrevBlockHash)==0{
			break
		}

	}

}
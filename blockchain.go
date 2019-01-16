package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"encoding/hex"
	"crypto/ecdsa"
	"bytes"
	"errors"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"


const genesisData =  "jonson blockchain"
type Blockchain struct{
	tip []byte //最近的一个区块的hash值
	db * bolt.DB
}


type BlockChainIterateor struct{
	currenthash []byte
	db * bolt.DB
}
func (bc * Blockchain) MineBlock(transations []*Transation ) *Block{

	for _,tx := range transations{
		if bc.VerifyTransation(tx)!=true{
			log.Panic("ERROR:INVALID transation")
		}else{
			fmt.Println("Verify Success")
		}
	}

	var lasthash []byte
	var lastheight int32
	err := bc.db.View(func(tx * bolt.Tx) error{
		b:= tx.Bucket([]byte(blockBucket))
		lasthash = b.Get([]byte("l"))
		blockdata:= b.Get(lasthash)
		block:=DeserializeBlock(blockdata)

		lastheight = block.Height
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	newBlock := NewBlock(transations,lasthash,lastheight+1)


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
		return newBlock
}


func NewBlockchain(address string) * Blockchain{
	var tip []byte
	db,err := bolt.Open(dbFile,0600,nil)
	if err!=nil{
		log.Panic(err)
	}

	err = db.Update(func(tx * bolt.Tx) error{
		b:= tx.Bucket([]byte(blockBucket))
		if b==nil{

			fmt.Println("区块链不存在，创建一个新的区块链")
			transation := NewCoinbaseTX(address,genesisData)

			genesis := NewGensisBlock([]*Transation{transation})
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

	 set:=UTXOSet{&bc}
	 set.Reindex()

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

func (bc * Blockchain) FindUnspentTransations(pubkeyhash []byte) []Transation{
	var unspentTXs []Transation //所有未花费的交易

	spendTXOs := make(map[string][]int) //  string 交易的哈希值 --->   []int 已经被花费的输出的序号   存储已经花费的交易

	bci := bc.iterator()

	for{
		block:=bci.Next()

		for _,tx := range block.Transations{
			txID := hex.EncodeToString(tx.ID)

		output:
			for outIdx,out := range tx.Vout{
				if spendTXOs[txID] !=nil{
					for _,spentOut:=range spendTXOs[txID]{
						if spentOut == outIdx{
							continue output
						}
					}
				}

				if out.CanBeUnlockedWith(pubkeyhash){
					unspentTXs = append(unspentTXs,*tx)
				}
			}

			if tx.IsCoinBase()==false{
				for _,in :=range tx.Vin{
					if in.canUnlockOutputWith(pubkeyhash){
						inTxId := hex.EncodeToString(in.TXid)
						spendTXOs[inTxId] = append(spendTXOs[inTxId],in.Voutindex)
					}
				}
			}

		}

		if len(block.PrevBlockHash) == 0{
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindUTXO(pubkeyhash []byte) []TXOutput{
	var UTXOs []TXOutput

	unspendTransations := bc.FindUnspentTransations(pubkeyhash)

	for _,tx :=range unspendTransations{
		for _,out := range tx.Vout{
			if out.CanBeUnlockedWith(pubkeyhash){
				UTXOs = append(UTXOs,out)
			}
		}
	}

	return UTXOs


}

func (bc *Blockchain) FindSpendableOutputs(pubkeyhash []byte,amount int) (int,map[string][]int) {

	unspentOutputs := make(map[string][]int)


	unspentTXs := bc.FindUnspentTransations(pubkeyhash)

	accumulated :=0



	Work:
	for _,tx := range unspentTXs{
		txID := hex.EncodeToString(tx.ID)

		for outIdx,out := range tx.Vout{
			if out.CanBeUnlockedWith(pubkeyhash) && accumulated < amount{

				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID],outIdx)

				if accumulated >=amount{
					break Work
				}
			}
		}
	}

	return accumulated,unspentOutputs
}


func (bc*Blockchain) SignTransation(tx *Transation,prikey ecdsa.PrivateKey){

	prevTXs :=  make(map[string]Transation)

	for _,vin := range tx.Vin{
		prevTX ,err :=  bc.FindTransationById(vin.TXid)

		if err !=nil{
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(prikey,prevTXs)
}

func (bc * Blockchain) FindTransationById(ID []byte)(Transation,error){

	bci := bc.iterator()


	for{
		block := bci.Next()
			for _,tx := range block.Transations{

				if bytes.Compare(tx.ID,ID)==0{
					return *tx,nil
				}
			}




		if len(block.PrevBlockHash) ==0{
			break
		}
	}

	return Transation{},errors.New("transation is not find ")
}

func (bc *Blockchain) VerifyTransation(tx * Transation) bool{
	prevTXs := make(map[string]Transation)

	for _,vin := range tx.Vin{
		prevTx,err:= bc.FindTransationById(vin.TXid)

		if err !=nil{
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	return tx.Verify(prevTXs)
}

func (bc * Blockchain) FindALLUTXO() map[string]TXOutputs{
	UTXO := make(map[string]TXOutputs)

	spentTXs := make(map[string][]int)

	bci := bc.iterator()

	for{

		block:= bci.Next()

		for _,tx :=range block.Transations{
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx,out := range tx.Vout{

				if spentTXs[txID] !=nil{
					for _,spendOutIds := range spentTXs[txID]{
						if spendOutIds == outIdx{
							continue Outputs
						}
					}
				}
				outs:=UTXO[txID]
				outs.Outputs = append(outs.Outputs,out)
				UTXO[txID] = outs
			}


			if tx.IsCoinBase()==false{
				for _,in := range tx.Vin{
					inTXID := hex.EncodeToString(in.TXid)

					spentTXs[inTXID] = append(spentTXs[inTXID],in.Voutindex)
				}
			}

		}


		if len(block.PrevBlockHash)==0{
			break
		}
	}

	return UTXO

}



func (bc *Blockchain) GetBestHeight() int32{
	var lastBlock Block

	err := bc.db.View(func(tx*bolt.Tx) error{
		b:=tx.Bucket([]byte(blockBucket))
		lastHash := b.Get([]byte("l"))
		blockdata := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockdata)

		return nil
	})

	if err !=nil{
		log.Panic(err)
	}
	return lastBlock.Height
}
func (bc *Blockchain) Getblockhash() [][]byte {
	var blocks [][]byte

	bci:=bc.iterator()

	for{
		block:=bci.Next()

		blocks= append(blocks,block.Hash)

		if len(block.PrevBlockHash) ==0{
			break
		}
	}





	return blocks
}

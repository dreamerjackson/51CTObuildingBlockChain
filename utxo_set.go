package main

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
)
type UTXOSet struct{
	bchain * Blockchain
}



const utxoBucket = "chainset"


func (u UTXOSet) Reindex(){

	db:=u.bchain.db

	bucketName  :=[]byte(utxoBucket)


	err := db.Update(func(tx * bolt.Tx) error{
		err2 := tx.DeleteBucket(bucketName)

		if err2 !=nil && err2 != bolt.ErrBucketNotFound{
			log.Panic(err2)
		}


		_,err3 := tx.CreateBucket(bucketName)
		if err3 !=nil{
			log.Panic(err3)
		}

		return nil

	})



	if err !=nil {
		log.Panic(err)
	}

	UTXO := u.bchain.FindALLUTXO()

	err4 := db.Update(func(tx *bolt.Tx) error{
		b:= tx.Bucket(bucketName)

		for txID,outs := range UTXO{
			key,err5:= hex.DecodeString(txID)
			if err5!=nil{
				log.Panic(err5)
			}

			err6:=b.Put(key,outs.Serialize())
			if err6 !=nil{
				log.Panic(err6)
			}
		}

		return nil

	})

	if err4!=nil{

		log.Panic(err4)
	}

}

func (u UTXOSet) FindUTXObyPubkeyHash(pubkeyhash []byte) []TXOutput{
	var UTXOs []TXOutput

	db := u.bchain.db

	err := db.View(func(tx * bolt.Tx) error {

		b:= tx.Bucket([]byte(utxoBucket))

		c := b.Cursor()

		for k,v:=c.First();k!=nil;k,v=c.Next(){
			outs := DeserializeOutputs(v)

			for _,out := range outs.Outputs{
				if out.CanBeUnlockedWith(pubkeyhash){
					UTXOs = append(UTXOs,out)
				}
			}

		}
		return nil
	})
	if err!=nil{

		log.Panic(err)
	}
		return UTXOs
}
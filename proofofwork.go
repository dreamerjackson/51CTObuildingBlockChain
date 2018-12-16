package main

import (
	"math/big"
	"bytes"
	"crypto/sha256"
)

type  ProofOfWork struct{
	block * Block
	tartget * big.Int
}
const targetBits = 16


func NewProofofWork(b*Block) * ProofOfWork{

	target := big.NewInt(1)
	target.Lsh(target,uint(256-targetBits))
	pow := &ProofOfWork{b,target}
	return pow
}

func (pow * ProofOfWork) prepareData(nonce int32) []byte{

	data := bytes.Join(
		[][]byte{
			IntToHex(pow.block.Version),
			pow.block.PrevBlockHash,
			pow.block.Merkleroot,
			IntToHex(pow.block.Time),
			IntToHex(pow.block.Bits),
			IntToHex(nonce)},
		[]byte{},
	)
	return data
}

func (pow * ProofOfWork) Run() (int32,[]byte){

		var nonce int32
		var secondhash [32]byte
		nonce = 0
		var currenthash big.Int

		for  nonce < maxnonce{


			//序列化
			data:= pow.prepareData(nonce)
			//double hash
			fitstHash := sha256.Sum256(data)
			secondhash = sha256.Sum256(fitstHash[:])
		//	fmt.Printf("%x\n",secondhash)


			currenthash.SetBytes(secondhash[:])
			//比较
			if currenthash.Cmp(pow.tartget) == -1{
				break
			}else{
				nonce++
			}
		}


		return nonce,secondhash[:]
}

func (pow * ProofOfWork) Validate() bool{
	var hashInt big.Int

	data:=pow.prepareData(pow.block.Nonce)

	fitstHash := sha256.Sum256(data)
	secondhash := sha256.Sum256(fitstHash[:])
	hashInt.SetBytes(secondhash[:])
	isValid:= hashInt.Cmp(pow.tartget) == -1

	return isValid
}
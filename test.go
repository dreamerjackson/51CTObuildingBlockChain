package main

import "fmt"

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
		0,
	}



	txin := TXInput{[]byte{},-1,nil,nil}
	txout := NewTXOutput(subsidy,"first")
	tx := Transation{nil,[]TXInput{txin},[]TXOutput{*txout}}

	txin2 := TXInput{[]byte{},-1,nil,nil}
	txout2 := NewTXOutput(subsidy,"second")
	tx2 := Transation{nil,[]TXInput{txin2},[]TXOutput{*txout2}}

	var Transations []*Transation

	Transations = append(Transations,&tx,&tx2)

	block.createMerkelTreeRoot(Transations)

	fmt.Printf("%x\n",block.Merkleroot)
}

func TestNewSerialize(){

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
		0,
	}

	deBlock:=DeserializeBlock(block.Serialize())

	deBlock.String()


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
		0,
	}

	pow:=NewProofofWork(block)

	nonce,_:= pow.Run()

	block.Nonce = nonce

	fmt.Println("POW:",pow.Validate())

}

func TestBoltDB(){
	blockchain := NewBlockchain("1NeBzmfLDxinqHwNdzoA5y8c5fYgZgiUds")
	blockchain.MineBlock([]*Transation{})
	blockchain.MineBlock([]*Transation{})
	blockchain.printBlockchain()
}
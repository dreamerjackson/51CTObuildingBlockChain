package main

import (
	"fmt"
	"net"
	"log"
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
)
const nodeversion = 0x00

var nodeAddress string
var blockInTransit  [][]byte
const commonLength = 12
type Version struct {
	Version int
	BestHeight int32
	AddrFrom string
}


func (ver *Version) String(){
	fmt.Printf("Version:%d\n",ver.Version)
	fmt.Printf("BestHeight:%d\n",ver.BestHeight)
	fmt.Printf("AddrFrom:%s\n",ver.AddrFrom)
}


var knownNodes = []string{"localhost:3000"}

func StartServer(nodeID,minerAddress string,bc*Blockchain){

	nodeAddress = fmt.Sprintf("localhost:%s",nodeID)
	ln,err:= net.Listen("tcp",nodeAddress)
	defer ln.Close()

	//bc := NewBlockchain("1NeBzmfLDxinqHwNdzoA5y8c5fYgZgiUds")

	if nodeAddress !=knownNodes[0]{
		sendVersion(knownNodes[0],bc)
	}

	for{

		conn,err2:=ln.Accept()

		if err2 != nil{
			log.Panic(err)
		}
		go handleConnction(conn,bc)
	}
}


func handleConnction(conn net.Conn, bc *Blockchain) {

	request,err := ioutil.ReadAll(conn)

	if err !=nil{
		log.Panic(err)
	}

		//获取命令
	command:= bytesToCommand(request[:commonLength])
	fmt.Println(command)
	switch command {
	case "version":
	fmt.Printf("\nstr:获取version\n")
		handleVersion(request,bc)
	case "getblocks":
		handleGetBlock(request,bc)
	case "inv":
		handleInv(request,bc)
	case "getdata":
		handleGetData(request,bc)
	case "block":
		handleBlock(request,bc)
	}

}
func handleBlock(request []byte, bc *Blockchain) {

	var buff bytes.Buffer

	var payload blocksend

	buff.Write(request[commonLength:])
	dec:= gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err !=nil{
		log.Panic(err)
	}

	blockdata:= payload.Block

	block:= DeserializeBlock(blockdata)
	bc.AddBlock(block)
	fmt.Printf("Recieve a new Block")

	if len(blockInTransit) >0{
		blockHash:= blockInTransit[0]
		sendGetData(payload.AddrFrom,"block",blockHash)
		blockInTransit = blockInTransit[1:]
	}else{
		set:= UTXOSet{bc}
		set.Reindex()
	}
}
func handleGetData(request []byte, bc *Blockchain) {
	fmt.Println("jkjkjhgffghjjkjkkkkkkkkkkkkkkk")
	var buff bytes.Buffer
	var payload getdata
	buff.Write(request[commonLength:])
	dec:=gob.NewDecoder(&buff)
	err:= dec.Decode(&payload)
	if err !=nil{
		log.Panic(err)
	}

	if payload.Type=="block"{
		fmt.Printf("payload.ID:%x\n",payload.ID)
		block ,err:= bc.GetBlock([]byte(payload.ID))
		//fmt.Println("g5")
		if err!=nil{
			log.Panic(err)
		}
		fmt.Println("g6: ",payload.AddrFrom)
		sendBlock(payload.AddrFrom,&block)
	}

}

type blocksend struct {
	AddrFrom string
	Block []byte
}

func sendBlock(addr string, block *Block) {
	fmt.Println("发送block: ",addr)
	data:= blocksend{nodeAddress,block.Serialize()}
	payload := gobEncode(data)

	request:= append(commandToBytes("block"),payload...)

	sendData(addr,request)
}
func handleInv(request []byte, bc *Blockchain) {
	var buff bytes.Buffer

	var payload inv
	buff.Write(request[commonLength:])

	dec:= gob.NewDecoder(&buff)
	err:= dec.Decode(&payload)

	if err !=nil{
		log.Panic(err)
	}


	fmt.Printf("Recieve inventory %d,%s",len(payload.Items),payload.Type)


	if payload.Type =="block"{
		blockInTransit = payload.Items

		blockHash:= payload.Items[0]
		sendGetData(payload.AddrFrom,"block",blockHash)

		newInTransit := [][]byte{}

		for _,b:= range blockInTransit{
			if bytes.Compare(b,blockHash)!=0{
				newInTransit = append(newInTransit,b)
			}
		}
		blockInTransit =  newInTransit
	}
}

type getdata struct {
	AddrFrom string
	Type string
	ID []byte
}

func sendGetData(addr string, kind string, id []byte) {
	payload:= gobEncode(getdata{nodeAddress,kind,id})

	request:= append(commandToBytes("getdata"),payload...)

	sendData(addr,request)
}
func handleGetBlock(request []byte, bc *Blockchain) {


	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commonLength:])

	dec:= gob.NewDecoder(&buff)
	err:= dec.Decode(&payload)

	if err !=nil{
		log.Panic(err)
	}

	//block:= bc.Getblockhash()
	block:= bc.GetBlockHashScope(payload.LowHeight,payload.HighHeight)
	fmt.Println("sendenv: ",payload.Addrfrom)
	sendInv(payload.Addrfrom,"block",block)
}

type inv struct {
	AddrFrom string
	Type string
	Items [][]byte
}

func sendInv(addr string, kind string, items [][]byte) {
	inventory:= inv{nodeAddress,kind,items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"),payload...)

	sendData(addr,request)
}

func handleVersion(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload Version
	buff.Write(request[commonLength:])

	dec:=gob.NewDecoder(&buff)

	err:= dec.Decode(&payload)

	if err!=nil{
		log.Panic(err)
	}
	payload.String()
	myBestHeight := bc.GetBestHeight()

	foreignerBestHeight :=  payload.BestHeight


	if myBestHeight < foreignerBestHeight{
		sendGetBlock(payload.AddrFrom,myBestHeight+1,foreignerBestHeight)
	}else{

		sendVersion(payload.AddrFrom,bc)
	}


	if !nodeIsKnow(payload.AddrFrom){
		knownNodes = append(knownNodes,payload.AddrFrom)
	}

}

type getblocks struct {
	Addrfrom    string    //命令发送方地址，用于对方应答回来
	LowHeight   int32     //区块高度--低
	HighHeight  int32     //区块高度--高
}

func sendGetBlock(address string, low int32, high int32) {
	payload:=  gobEncode(getblocks{nodeAddress,low, high})

	request:= append(commandToBytes("getblocks"),payload...)

	sendData(address,request)


}
func nodeIsKnow(addr string) bool {


	for _,node :=range knownNodes{
		if node ==addr{
			return true
		}
	}

	return false
}
func sendVersion(addr string, bc *Blockchain) {
	bestHeight :=bc.GetBestHeight()

	payload := gobEncode(Version{nodeversion,bestHeight,nodeAddress})
	request:=append(commandToBytes("version"),payload...)
	sendData(addr,request)
}
func sendData(addr string, data []byte) {
	con,err := net.Dial("tcp",addr)
	if err !=nil{

		fmt.Printf("%s is no available",addr)


		var updateNodes []string

		for _,node:=range knownNodes{
			if node !=addr{
				updateNodes = append(updateNodes,node)
			}
		}

		knownNodes = updateNodes
	}
	defer con.Close()

	_,err = io.Copy(con,bytes.NewReader(data))

	if err !=nil{
		log.Panic(err)
	}

}
func commandToBytes(command string) []byte {
	var bytes [commonLength]byte

	for i,c:= range command{
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string{
	var command []byte

	for _,b:=range bytes{
		if b!=0x00{
			command = append(command,b)
		}
	}

	return fmt.Sprintf("%s",command)
}




func gobEncode(data interface{}) []byte{
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)

	if err!=nil{
		log.Panic(err)

	}

	return buff.Bytes()

}
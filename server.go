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

	}

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
		//sendGetBlock(payload.AddrFrom)
	}else{

		sendVersion(payload.AddrFrom,bc)
	}


	if !nodeIsKnow(payload.AddrFrom){
		knownNodes = append(knownNodes,payload.AddrFrom)
	}

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
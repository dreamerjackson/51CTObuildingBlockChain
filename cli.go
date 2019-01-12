package main

import (
	"os"
	"fmt"
	"flag"
	"log"
)

type CLI struct{

	bc * Blockchain
}


func (cli * CLI) addBlock(){
	cli.bc.MineBlock([]*Transation{})
}
func (cli * CLI) validateArgs(){
	if len(os.Args) < 1{
		fmt.Println("参数小于1")
		os.Exit(1)
	}
	fmt.Println(os.Args)
}

func (cli * CLI ) printChain(){

	cli.bc.printBlockchain()
}

func (cli*CLI) send(from,to string,amount int){

	tx:= NewUTXOTransation(from,to,amount,cli.bc)

	cli.bc.MineBlock([]*Transation{tx})

	fmt.Printf("Success!")
}
func (cli * CLI) getBalance(address string){


	balance := 0

	decodeAddress := Base58Decode([]byte(address))
	pubkeyHash:= decodeAddress[1:len(decodeAddress)-4]

	set := UTXOSet{cli.bc}

	//UTXOs := cli.bc.FindUTXO(pubkeyHash)
	UTXOs := set.FindUTXObyPubkeyHash(pubkeyHash)
	for _,out := range UTXOs{
		balance += out.Value
	}

	fmt.Printf("\nbalance of '%s':%d\n",address,balance)
}



func (cli *CLI) createWallet(){
		wallets,_:=NewWallets()
		address :=wallets.CreateWallet()
		wallets.SaveToFile()
		fmt.Printf("your address:%s\n",address)

}


func (cli * CLI) listAddress(){


	wallets,err:=NewWallets()
	if err!=nil{
		log.Panic(err)
	}

	addresses :=  wallets.getAddress()

	for _,address := range addresses{

		fmt.Println(address)
	}

}

func (cli * CLI) printUsage(){

	fmt.Println("USages:")
	fmt.Println("addblock-增加区块:")
	fmt.Println("printChain:打印区块链")
}
func (cli * CLI) Run(){
	cli.validateArgs()

	addBlockCmd  := flag.NewFlagSet("addblock",flag.ExitOnError)
	printChianCmd  := flag.NewFlagSet("printChian",flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address","","the address to get balance of ")

	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)

	sendFrom := sendCmd.String("from","","source wallet address")
	sendTo := sendCmd.String("to","","Destination wallet address")
	sendAmount := sendCmd.Int("amount",0,"Amount to send")


	createWalletCMD := flag.NewFlagSet("createwallet",flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listaddress",flag.ExitOnError)

	switch os.Args[1]{
	case "createwallet":
		err:=createWalletCMD.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "listaddress":
		err:=listAddressCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}

	case "send":
		err:=sendCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "getbalance":
		err:=getBalanceCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "addblock":
		err:=addBlockCmd.Parse(os.Args[2:])

		if err!=nil{
			log.Panic(err)
		}
	case "printChian":
		err:=printChianCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}

	default:
		cli.printUsage()
		os.Exit(1)
	}


	if addBlockCmd.Parsed(){
		cli.addBlock()
	}


	if printChianCmd.Parsed(){
		cli.printChain()
	}

	if getBalanceCmd.Parsed(){
			if *getBalanceAddress == ""{
				os.Exit(1)

			}
			cli.getBalance(*getBalanceAddress)
	}

	if sendCmd.Parsed(){
		if *sendFrom == "" || *sendTo=="" || *sendAmount <=0{
			os.Exit(1)

		}
		fmt.Println(*sendFrom,*sendTo,*sendAmount)
		cli.send(*sendFrom,*sendTo,*sendAmount)
	}

	if createWalletCMD.Parsed(){
		cli.createWallet()
	}
	if listAddressCmd.Parsed(){
		cli.listAddress()
	}

}
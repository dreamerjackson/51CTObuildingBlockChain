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


func (cli * CLI) getBalance(address string){
	balance := 0
	UTXOs := cli.bc.FindUTXO(address)

	for _,out := range UTXOs{
		balance += out.Value
	}

	fmt.Printf("\nbalance of '%s':%d\n",address,balance)
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


	switch os.Args[1]{
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

}
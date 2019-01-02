package main

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"io/ioutil"
	"crypto/elliptic"
	"os"
)

const walletFile = "wallet.dat"

type Wallets struct{

	Walletsstore map[string]*Wallet
}

func NewWallets() (*Wallets,error){
	wallets := Wallets{}

	wallets.Walletsstore =  make(map[string]*Wallet)

	err:= wallets.LoadFromFile()

	return &wallets,err


}

func (ws *Wallets) CreateWallet() string{

	wallet := Newwallet()

	address := fmt.Sprintf("%s",wallet.GetAddress())
	ws.Walletsstore[address] = wallet

	return address
}

func (ws * Wallets) GetWallet(address string) Wallet{
	return *ws.Walletsstore[address]
}

func (ws * Wallets) getAddress() []string{

	var addresses []string

	for address,_ := range ws.Walletsstore{

		addresses =  append(addresses,address)
	}


	return addresses
}


func (ws * Wallets) LoadFromFile() error{
	if _ ,err :=  os.Stat(walletFile);os.IsNotExist(err){
		return err
	}

		fileContent,err:=ioutil.ReadFile(walletFile)

		if err !=nil{
			log.Panic(err)
		}

		var wallets Wallets
		gob.Register(elliptic.P256())
		decoder := gob.NewDecoder(bytes.NewReader(fileContent))
		err = decoder.Decode(&wallets)

		if err !=nil{
			log.Panic(err)
		}

		ws.Walletsstore =  wallets.Walletsstore

		return nil
}


func (ws *Wallets) SaveToFile(){

	var content bytes.Buffer

	gob.Register(elliptic.P256())
	encoder:= gob.NewEncoder(&content)

	err := encoder.Encode(ws)
	if err !=nil{

		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile,content.Bytes(),0777)
	if err !=nil{

		log.Panic(err)
	}





}
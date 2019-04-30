package main

func main(){
	bc :=  NewBlockchain("1NeBzmfLDxinqHwNdzoA5y8c5fYgZgiUds")

	cli := CLI{bc}
	cli.Run()

	//wallet := Newwallet()
	//
	//
	//fmt.Printf("私钥：%x\n",wallet.PrivateKey.D.Bytes())
	//fmt.Printf("公钥：%x\n", wallet.Publickey)
	//fmt.Printf("地址：%x\n", wallet.GetAddress())
	//address,_:=hex.DecodeString("3146536551465558616f7631635a434d4e424a343834707663616754676765473936")
	//fmt.Printf("%d\n",ValidateAddress(address))
}

package main

import (
  "fmt"
  "log"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/ethclient"
  instapay "github.com/sslab-instapay/instapay-go-server/contracts"
)

func main() {
  client, err := ethclient.Dial("ws://141.223.121.139:8881")
  if err != nil {
    log.Fatal(err)
  }

  address := common.HexToAddress("0xD92151dd38931d83afB31A8206ca410cef2A9862")
  instance, err := instapay.NewInstapay(address, client)
  if err != nil {
    fmt.Printf("ERROR 1: ")
    log.Fatal(err)
  }

  readme, err := instance.Readme(nil)
  if err != nil {
    fmt.Printf("ERROR 2: ")
    log.Fatal(err)
  }

  fmt.Println(readme)
}

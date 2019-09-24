package main

import (
  "fmt"
  "log"
  "sync"
  "github.com/ethereum/go-ethereum/ethclient"
  "github.com/ethereum/go-ethereum/common"
  instapay "github.com/sslab-instapay/instapay-go-server/contracts"
)

func rt(wg *sync.WaitGroup, client *ethclient.Client) {
  defer wg.Done()

  fmt.Printf("rt routine start !\n")

  address := common.HexToAddress("0x3016947BE73dcb877401Ee33802aC8fA6feE631E")
  instance, err := instapay.NewInstapay(address, client)
  if err != nil {
    log.Fatal(err)
  }

  readme, err := instance.Readme(nil)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("%d\n", readme)
}

func dummy() {
  fmt.Printf("I am dummy\n")
}

func Start() {
  var wg sync.WaitGroup
  wg.Add(1)

  fmt.Printf("Start()\n")

  client, err := ethclient.Dial("ws://141.223.121.139:8881")
  if err != nil {
    log.Fatal(err)
  }

  go dummy()
  go rt(&wg, client)
  wg.Wait()

  fmt.Printf("End of Start()\n")
}

func main() {
  Start()
}

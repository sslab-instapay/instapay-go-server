package main

import (
  "fmt"
  "log"
  "github.com/sslab-instapay/instapay-go-server/repository"
)

func main() {
  info, err := repository.GetClientInfo("0xD03A2CC08755eC7D75887f0997195654b928893e")
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println((*info).IP)
}

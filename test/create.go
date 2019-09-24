package main

import (
  "context"
  "crypto/ecdsa"
  "fmt"
  "log"
  "math/big"
  "github.com/ethereum/go-ethereum/accounts/abi/bind"
  "github.com/ethereum/go-ethereum/crypto"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/ethclient"
  instapay "github.com/sslab-instapay/instapay-go-server/contracts"
)

func main() {
  client, err := ethclient.Dial("ws://141.223.121.139:8881")
  if err != nil {
    log.Fatal(err)
  }

  // loading instapay contract on the blockchain
  address := common.HexToAddress("0x0FD3915b4b5a474896c6b5DC194489Bde6820815")  // change to correct address
  instance, err := instapay.NewInstapay(address, client)
  if err != nil {
    log.Fatal(err)
  }

  // loading my public key, nonce and gas price e113ff405699b7779fbe278ee237f2988b1e6769d586d8803860d49f28359fbd
  privateKey, err := crypto.HexToECDSA("240af81838ad22e8baa5c6223c7c7e112b091ba50e6fb396c0dc2b84cf034169")
  if err != nil {
    log.Fatal(err)
  }

  publicKey := privateKey.Public()
  publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
  if !ok {
    log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
  }

  fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

  nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
  if err != nil {
    log.Fatal(err)
  }

  gasPrice, err := client.SuggestGasPrice(context.Background())
  if err != nil {
    log.Fatal(err)
  }

  // composing a transaction
  auth := bind.NewKeyedTransactor(privateKey)
  auth.Nonce = big.NewInt(int64(nonce))
  auth.Value = big.NewInt(8000000000000000000) // in wei
  auth.GasLimit = uint64(2000000) // in units
  auth.GasPrice = gasPrice

  receiver := common.HexToAddress("0x78902c58006916201F65f52f7834e467877f0500")

  tx, err := instance.CreateChannel(auth, receiver)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("tx sent: %s\n", tx.Hash().Hex())
}

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

  // loading my public key, nonce and gas price
  privateKey, err := crypto.HexToECDSA("e113ff405699b7779fbe278ee237f2988b1e6769d586d8803860d49f28359fbd")
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
  auth.Value = big.NewInt(0)
  auth.GasLimit = uint64(2000000) // in units
  auth.GasPrice = gasPrice

  id := big.NewInt(15)
  owner_bal := big.NewInt(7)
  receiver_bal := big.NewInt(1)

  tx, err := instance.CloseChannel(auth, id, owner_bal, receiver_bal)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("tx sent: %s\n", tx.Hash().Hex())
}

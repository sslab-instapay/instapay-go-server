package main

import (
	"log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"github.com/ethereum/go-ethereum"
	"math/big"
	"fmt"
	"reflect"
	instapay "github.com/sslab-instapay/instapay-go-server/contracts"
)

var EthereumConfig = map[string]string{
	/* web3 and ethereum */
	"rpcHost":          "141.223.121.139",
	"rpcPort":          "8555",
	"wsHost":           "141.223.121.139", //141.223.121.139
	"wsPort":           "8881",
	"contractAddr":     "0x0FD3915b4b5a474896c6b5DC194489Bde6820815",   // change to correct address

	/* grpc configuration */
	"serverGrpcHost": "141.223.121.139",
	"serverGrpcPort": "50004",
	"serverProto":    "",
	"server":         "",
	"myGrpcPort":     "", //process.argv[3]
	"clientProto":    "",
	"receiver":       "",
}

func main(){
	client, err := ethclient.Dial("ws://" + EthereumConfig["wsHost"] + ":" + EthereumConfig["wsPort"])
	if err != nil {
		log.Fatal(err)
	}
	contractAddress := common.HexToAddress(EthereumConfig["contractAddr"])

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(instapay.InstapayABI)))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			createChannelEvent := struct {
				Id       *big.Int
				Owner    common.Address
				Receiver common.Address
				Deposit  *big.Int
			}{}
			closeChannelEvent := struct {
				Id          *big.Int
				Ownerbal    *big.Int
				Receiverbal *big.Int
			}{}
			ejectEvent := struct {
				Pn              *big.Int
				Registeredstage int
			}{}

			err := contractAbi.Unpack(&createChannelEvent, "EventCreateChannel", vLog.Data)
			if err == nil {
				fmt.Printf("Channel ID       : %d\n", createChannelEvent.Id)
				fmt.Printf("Channel Onwer    : %s\n", createChannelEvent.Owner.Hex())
				fmt.Printf("Channel Receiver : %s\n", createChannelEvent.Receiver.Hex())
				fmt.Printf("Channel Deposit  : %d\n", createChannelEvent.Deposit)
				continue;
			}

			err = contractAbi.Unpack(&closeChannelEvent, "EventCloseChannel", vLog.Data)
			if err == nil {
				fmt.Printf("Channel ID       : %d\n", closeChannelEvent.Id)
				fmt.Printf("Owner Balance    : %d\n", closeChannelEvent.Ownerbal)
				fmt.Printf("Receiver Balance : %d\n", closeChannelEvent.Receiverbal)
				continue;
			}

			err = contractAbi.Unpack(&ejectEvent, "EventEject", vLog.Data)
			if err == nil {
				fmt.Printf("Payment Number   : %d\n", ejectEvent.Pn)
				fmt.Printf("Stage            : %d\n", ejectEvent.Registeredstage)
				continue;
			}
		}
	}
}

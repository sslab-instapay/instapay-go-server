package main

import (
  "github.com/sslab-instapay/instapay-go-server/repository"
  pbClient "github.com/sslab-instapay/instapay-go-server/proto/client"
  "fmt"
  "log"
)

func ReturnArray() []string {
  var a []string
  a = []string{"A", "B"}
  return a
}

func SearchPath(pn int64, amount int64) ([]string, map[string]pbClient.AgreeRequestsMessage)  {
  var a []string
  var channelID1 int64
  var channelID2 int64

  /* composing a */
  a = []string{"0xD03A2CC08755eC7D75887f0997195654b928893e", "0x0b4161ad4f49781a821C308D672E6c669139843C", "0x78902c58006916201F65f52f7834e467877f0500"}

  /* composing w */
  channels, err := repository.GetChannelList()
  if err != nil {
    log.Fatal(err)
  }

  for i := 0; i < len(channels); i++ {
    if channels[i].From == "0xD03A2CC08755eC7D75887f0997195654b928893e" {
      channelID1 = int64(channels[i].ChannelId)
    } else if channels[i].From == "0x0b4161ad4f49781a821C308D672E6c669139843C" {
      channelID2 = int64(channels[i].ChannelId)
    }
  }

  var w map[string]pbClient.AgreeRequestsMessage
  w = make(map[string]pbClient.AgreeRequestsMessage)

  channelID1 = int64(channelID1)
  channelID2 = int64(channelID2)
  amount = int64(amount)
  pn = int64(pn)

  var cps1 []*pbClient.ChannelPayment
  cps1 = append(cps1, &pbClient.ChannelPayment{ChannelId: channelID1, Amount: -amount})
  rqm1 := pbClient.AgreeRequestsMessage{
    PaymentNumber: pn,
    ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cps1},
    Amount: amount}
  w["0xD03A2CC08755eC7D75887f0997195654b928893e"] = rqm1

  var cps2 []*pbClient.ChannelPayment
  cps2 = append(cps2, &pbClient.ChannelPayment{ChannelId: channelID1, Amount: amount})
  cps2 = append(cps2, &pbClient.ChannelPayment{ChannelId: channelID2, Amount: -amount})
  rqm2 := pbClient.AgreeRequestsMessage{
    PaymentNumber: pn,
    ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cps2},
    Amount: amount}
  w["0x0b4161ad4f49781a821C308D672E6c669139843C"] = rqm2

  var cps3 []*pbClient.ChannelPayment
  cps3 = append(cps1, &pbClient.ChannelPayment{ChannelId: channelID2, Amount: amount})
  rqm3 := pbClient.AgreeRequestsMessage{
    PaymentNumber: pn,
    ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cps3},
    Amount: amount}
  w["0x78902c58006916201F65f52f7834e467877f0500"] = rqm3

  return a, w
}

func main() {
  a := ReturnArray()
  fmt.Println(a)

  a, w := SearchPath(8989, 9000)
  fmt.Println(a)
  fmt.Println(w)

  for k, v := range a {
    fmt.Println(v)
    _ = k
  }
}

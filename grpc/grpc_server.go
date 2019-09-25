package grpc

import (
	"net"
	"log"
	"fmt"
	"os"
	"context"
	"strconv"
	"reflect"
	"math/rand"
	"google.golang.org/grpc"
	"github.com/sslab-instapay/instapay-go-server/repository"
	pbServer "github.com/sslab-instapay/instapay-go-server/proto/server"
	pbClient "github.com/sslab-instapay/instapay-go-server/proto/client"
)

type ServerGrpc struct {
}

/* wrapper function */
func WrapperAgreementRequest(pn int64, p []string, w map[string]pbClient.AgreeRequestsMessage) {
	/* remove C's address from p */
	var q []string
	q = p[0:2]

	for _, address := range q {
		info, err := repository.GetClientInfo(address)
		if err != nil {
			log.Fatal(err)
		}

		clientAddr := (*info).IP + ":" + strconv.Itoa((*info).Port)
		conn, err := grpc.Dial(clientAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		client := pbClient.NewClientClient(conn)
		message := w[address]
		_, err = client.AgreementRequest(context.Background(), &message)
		if err != nil {
			log.Println(err)
		}

		_, err = repository.UpdatePaymentAddrsSentAgr(int(pn), address)
		if err != nil {
			log.Println(err)
		}

		paymentData, err := repository.GetPaymentData(int(pn))
		if err != nil {
			log.Fatal(err)
		}

		if reflect.DeepEqual(paymentData.AddrsSentAgr, q) {
			go WrapperUpdateRequest(pn, p, w)
			return
		}
	}
}

func WrapperUpdateRequest(pn int64, p []string, w map[string]pbClient.AgreeRequestsMessage) {
	for _, address := range p {
		info, err := repository.GetClientInfo(address)
		if err != nil {
			log.Fatal(err)
		}

		clientAddr := (*info).IP + ":" + strconv.Itoa((*info).Port)
		conn, err := grpc.Dial(clientAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		/* convert AgreeRequestsMessage to UpdateRequestsMessage */
		rqm := pbClient.UpdateRequestsMessage{
			PaymentNumber:   w[address].PaymentNumber,
			ChannelPayments: w[address].ChannelPayments,
			Amount:          w[address].Amount}

		client := pbClient.NewClientClient(conn)

		// TODO 전부 Fatal 처리는 프로그램이 종료되기 때문에 log.println으로 띄우는게 데모때 나을거같아요.
		_, err = client.UpdateRequest(context.Background(), &rqm)
		if err != nil {
			log.Fatal(err)
		}

		// TODO 파라미터 오류
		_, err = repository.UpdatePaymentAddrsSentUpt(int(pn), address)
		if err != nil {
			log.Fatal(err)
		}

		paymentData, err := repository.GetPaymentData(int(pn))
		if err != nil {
			log.Fatal(err)
		}

		if reflect.DeepEqual(paymentData.AddrsSentUpt, p) {
			go WrapperConfirmPayment(int(pn), p)
			return
		}
	}
}

func WrapperConfirmPayment(pn int, p []string) {
	/* update payment's status */
	_, err := repository.UpdatePaymentStatus(int(pn), "SUCCESS")
	if err != nil {
		log.Fatal(err)
	}

	for _, address := range p {
		info, err := repository.GetClientInfo(address)
		if err != nil {
			log.Fatal(err)
		}

		clientAddr := (*info).IP + ":" + strconv.Itoa((*info).Port)
		conn, err := grpc.Dial(clientAddr, grpc.WithInsecure())
		defer conn.Close()
		if err != nil {
			log.Fatal(err)
		}

		client := pbClient.NewClientClient(conn)
		_, err = client.ConfirmPayment(context.Background(), &pbClient.ConfirmRequestsMessage{PaymentNumber: int64(pn)})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("======== Sent confirm to: " + clientAddr)
	}
}

func SearchPath(pn int64, amount int64) ([]string, map[string]pbClient.AgreeRequestsMessage) {
	var p []string
	var channelID1 int64
	var channelID2 int64

	/* composing p */
	p = []string{"0xD03A2CC08755eC7D75887f0997195654b928893e", "0x0b4161ad4f49781a821C308D672E6c669139843C", "0x78902c58006916201F65f52f7834e467877f0500"}

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
		PaymentNumber:   pn,
		ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cps1},
		Amount:          amount}
	w["0xD03A2CC08755eC7D75887f0997195654b928893e"] = rqm1

	var cps2 []*pbClient.ChannelPayment
	cps2 = append(cps2, &pbClient.ChannelPayment{ChannelId: channelID1, Amount: amount})
	cps2 = append(cps2, &pbClient.ChannelPayment{ChannelId: channelID2, Amount: -amount})
	rqm2 := pbClient.AgreeRequestsMessage{
		PaymentNumber:   pn,
		ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cps2},
		Amount:          amount}
	w["0x0b4161ad4f49781a821C308D672E6c669139843C"] = rqm2

	var cps3 []*pbClient.ChannelPayment
	cps3 = append(cps3, &pbClient.ChannelPayment{ChannelId: channelID2, Amount: amount})
	rqm3 := pbClient.AgreeRequestsMessage{
		PaymentNumber:   pn,
		ChannelPayments: &pbClient.ChannelPayments{ChannelPayments: cps3},
		Amount:          amount}
	w["0x78902c58006916201F65f52f7834e467877f0500"] = rqm3

	return p, w
}

func (s *ServerGrpc) PaymentRequest(ctx context.Context, rq *pbServer.PaymentRequestMessage) (*pbServer.Result, error) {
	from := rq.From
	to := rq.To
	amount := rq.Amount
	pn := rand.Intn(1000)

	p, w := SearchPath(int64(pn), amount)
	_, err := repository.PutPaymentData(pn, from, to, int(amount), p)
	if err != nil {
		log.Println(err)
	}

	go WrapperAgreementRequest(int64(pn), p, w)

	return &pbServer.Result{Result: true}, nil
}

func (s *ServerGrpc) CommunicationInfoRequest(ctx context.Context, address *pbServer.Address) (*pbServer.CommunicationInfo, error) {
	res, err := repository.GetClientInfo(address.Addr)
	if err != nil {
		log.Fatal(err)
	}

	return &pbServer.CommunicationInfo{IPAddress: res.IP, Port: int64(res.Port)}, nil
}

func StartGrpcServer() {
	grpcPort, err := strconv.Atoi(os.Getenv("grpc_port"))
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pbServer.RegisterServerServer(grpcServer, &ServerGrpc{})

	grpcServer.Serve(lis)
}

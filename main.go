package main

import (
	"flag"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/sslab-instapay/instapay-go-server/router"
	"github.com/sslab-instapay/instapay-go-server/config"
	serverGrpc "github.com/sslab-instapay/instapay-go-server/grpc"
)

func StartWebServer() {
	defaultRouter := gin.Default()
	defaultRouter.LoadHTMLGlob("templates/*")

	router.RegisterViewRouter(defaultRouter)
	defaultRouter.Run(":" + os.Getenv("port"))
}

func main() {
	portNum := flag.String("port", "3001", "port number")
	grpcPortNum := flag.String("grpc_port", "50001", "grpc_port number")
	flag.Parse()

	os.Setenv("port", *portNum)
	os.Setenv("grpc_port", *grpcPortNum)

	config.GetContract()
	go serverGrpc.StartGrpcServer()

	StartWebServer()
}

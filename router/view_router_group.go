package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	//"github.com/sslab-instapay/instapay-go-server/config"
	"github.com/sslab-instapay/instapay-go-server/repository"
	//"github.com/sslab-instapay/instapay-go-server/model"
)

func RegisterViewRouter(router *gin.Engine) {

	viewRouter := router.Group("templates")
	{
		viewRouter.GET("clients/list", func(context *gin.Context) {

			clients, err := repository.GetClientList()
			if err != nil {
				log.Fatal(err)
			}

			context.HTML(http.StatusOK, "client.tmpl", gin.H { "clients": clients })
		})

		viewRouter.GET("channels/list", func(context *gin.Context) {

			channels, err := repository.GetChannelList()
			if err != nil {
				log.Fatal(err)
			}

			context.HTML(http.StatusOK, "channel.tmpl", gin.H { "channels": channels })
		})

		viewRouter.GET("histories/list", func(context *gin.Context) {

			histories, err := repository.GetPaymentList()
			if err != nil {
				log.Fatal(err)
			}

			context.HTML(http.StatusOK, "history.tmpl", gin.H { "histories": histories })
		})
	}
}

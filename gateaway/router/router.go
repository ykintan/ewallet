// router/router.go
package router

import (
	"ewallet/gateaway/config"
	"ewallet/gateaway/service"

	"github.com/gin-gonic/gin"
)

func basicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		config.GetBasicAuthUsername(): config.GetBasicAuthPassword(),
	})
}

func SetupRouter(srv *service.Server) *gin.Engine {
	r := gin.Default()
	r.GET("/getUserByID/:userID", srv.GetUserByID)
	r.GET("/getWalletByUserID/:userID", srv.GetWalletByUserID)
	r.GET("/getTransactionByUserID/:userID", srv.GetTransactionByUserID)
	r.GET("/getUserAndBalanceWallet/:userID", srv.GetUserAndBalanceWallet)

	authorized := r.Group("/", basicAuth())
	{
		authorized.POST("/createUser", srv.CreateUser)
		authorized.POST("/transferWallet", srv.TransferWallet)
		authorized.POST("/topUp", srv.TopUp)
	}

	return r
}

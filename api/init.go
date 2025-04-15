package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var port = 10080

func Start() {
	if viper.GetInt("server.port") != 0 {
		port = viper.GetInt("server.port")
	}
	r := gin.Default()
	SetupRouter(r)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
}

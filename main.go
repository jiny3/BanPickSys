package main

import (
	"github.com/jiny3/BanPickSys/api"
	"github.com/jiny3/BanPickSys/pkg"
	"github.com/jiny3/gopkg/hookx"
	"github.com/sirupsen/logrus"
)

func init() {
	hookx.Init(&pkg.InitViper, &pkg.InitLog)
}

func main() {
	logrus.Debug("Starting server...")
	api.Start()
}

package main

import (
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/server/businessLogic"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	NAME = "business"
)

var log *zap.SugaredLogger
var errh *evented.ErrLogger

func main() {
	configure()

	logger, _ := zap.NewDevelopment(zap.AddCaller())
	log = logger.Sugar()

	defer log.Sync()
	log.Infow("Logger Configured")

	errh = &evented.ErrLogger{log}

	server := businessLogic.NewMockBusinessLogic(log, errh)

	log.Info(viper.GetUint("port"))
	port := uint16(viper.GetUint("port"))
	log.Infow("Starting Business Server...", "port", port)
	server.Listen(port)
}

func configure() {
	viper.SetConfigName(NAME)
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("c:/temp/")

	viper.SetEnvPrefix(NAME)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Warn(err)
		} else {
			log.Fatal(err)
		}
	}

}

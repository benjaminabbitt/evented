package main

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	NAME = "int"
)

var log *zap.SugaredLogger
var errh *evented.ErrLogger

func main() {
	setupConfig()

	logger, _ := zap.NewDevelopment(zap.AddCaller())
	log = logger.Sugar()

	errh = &evented.ErrLogger{log}

	log.Info("Starting...")
	target := "localhost:8080" //viper.GetString("commandHandler")
	log.Info(target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	errh.LogIfErr(err, fmt.Sprintf("Error dialing %s", target))
	ch := evented_core.NewCommandHandlerClient(conn)
	log.Info("Client Created...")
	id, err := uuid.NewRandom()
	protoId := evented_proto.UUIDToProto(id)
	pages := []*evented_core.CommandPage{&evented_core.CommandPage{
		Sequence:    0,
		Synchronous: false,
		Command:     nil,
	}}
	commandBook := &evented_core.CommandBook{
		Cover: &evented_core.Cover{
			Domain: "test",
			Root:   &protoId,
		},
		Pages: pages,
	}
	_, _ = ch.Handle(context.Background(), commandBook)
	log.Info("Done!")
}

func setupConfig() {
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

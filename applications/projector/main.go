package main

import (
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/applications/projector/rabbitmq"
	"go.uber.org/zap"
)

/*
Transceiver.  Dequeue from event passing system and translate to GRPC calls
 */
var log *zap.SugaredLogger
var errh *evented.ErrLogger

func main() {
	logger, _ := zap.NewDevelopment(zap.AddCaller())
	log = logger.Sugar()

	defer log.Sync()
	log.Infow("Logger Configured")

	errh = &evented.ErrLogger{Log: log}

	receiver := rabbitmq.RabbitMQReceiver{
		Log: log,
		Errh: errh,
	}
	receiver.Listen()
}

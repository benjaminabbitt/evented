package support

import (
	"fmt"
	"go.uber.org/zap"
	"net"
)

func OpenPort(port uint, log *zap.SugaredLogger) (lis net.Listener, err error) {
	log.Debugw("Opening port", "port", port)
	lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Errorw("Error opening port", port, err)
		return nil, err
	}
	log.Debugw("Listening", "port", port)
	return lis, err
}

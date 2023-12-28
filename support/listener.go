package support

import (
	"fmt"
	"go.uber.org/zap"
	"net"
)

func CreateListener(port uint, log *zap.SugaredLogger) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error(err)
	}
	return lis
}

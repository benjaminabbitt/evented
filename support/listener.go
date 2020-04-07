package support

import (
	"fmt"
	"github.com/benjaminabbitt/evented"
	"net"
)

func CreateListener(port uint16, errh *evented.ErrLogger) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	errh.LogIfErr(err, "Failed to Listen")
	return lis
}

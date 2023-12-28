package network

import "net"

func GetExternalAddrs() (externalAddrs []string) {
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					externalAddrs = addIfNotLoopback(v.IP, externalAddrs)
				case *net.IPAddr:
					externalAddrs = addIfNotLoopback(v.IP, externalAddrs)
				}
			}
		}
	}
	return externalAddrs
}

func addIfNotLoopback(addr net.IP, externalAddrs []string) (rExternalAddrs []string) {
	rExternalAddrs = externalAddrs
	if !addr.IsLoopback() {
		rExternalAddrs = append(rExternalAddrs, addr.String())
	}
	return rExternalAddrs
}

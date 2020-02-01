package transport


type EventTransportSender interface {
	Send(evt TransportEvent)
}


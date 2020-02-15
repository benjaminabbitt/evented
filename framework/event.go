package framework

import "github.com/golang/protobuf/ptypes/any"

type Event struct{
	Id string
	Sequence uint32
	Details any.Any
}


//type CQRSMessage interface {
//	GetId() string
//	GetSequence() uint32
//	GetDetails() interface{}
//}
//
//type EventFromProto struct {
//	page  *evented.EventPage
//	cover *evented.Cover
//}
//
//func NewEventFromProto(page *evented.EventPage, cover *evented.Cover) EventFromProto {
//	return EventFromProto{
//		page:  page,
//		cover: cover,
//	}
//}
//
//func (event EventFromProto) GetId() string {
//	return event.cover.Id
//}
//
//func (event EventFromProto) GetSequence() uint32 {
//	return event.page.Sequence
//}
//
//func (event EventFromProto) GetDetails() interface{} {
//	return event.page.Event
//}
//
//type CommandFromProto struct {
//	page  *evented.CommandPage
//	cover *evented.Cover
//}
//
//func NewCommandsFromProto(book *evented.CommandBook) []CommandFromProto{
//	var commandsFromProto []CommandFromProto
//	for _, page := range book.Pages {
//		commandsFromProto = append(commandsFromProto, NewCommandFromProto(page, book.Cover))
//	}
//	return commandsFromProto
//}
//
//func NewCommandFromProto(page *evented.CommandPage, cover *evented.Cover) CommandFromProto {
//	return CommandFromProto{
//		page:  page,
//		cover: cover,
//	}
//}
//
//func (command CommandFromProto) GetId() string{
//	return command.cover.Id
//}
//func (command CommandFromProto) GetSequence() uint32{
//	return command.page.Sequence
//}
//func (command CommandFromProto) GetDetails() interface{}{
//	return command.page.Command
//}

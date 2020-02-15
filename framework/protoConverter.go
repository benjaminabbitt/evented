package framework

import (
	evented_proto "github.com/benjaminabbitt/evented/proto/core"
	"github.com/thoas/go-funk"
)

func CommandBookToCommand(proto *evented_proto.CommandBook) []Command {
	return funk.Map(proto.Pages, func(page evented_proto.CommandPage) Command {
		return Command{
			Id:       proto.Cover.Id,
			Sequence: page.Sequence,
			Details:  page.Command,
		}
	}).([]Command)
}

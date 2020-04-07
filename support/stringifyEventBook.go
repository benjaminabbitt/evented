package support

import (
	"fmt"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"strings"
)

func StringifyEventBook(eb *evented_core.EventBook) string {
	var pages string
	for _, page := range eb.Pages{
		seq, _ := GetSequence(page)
		pages += fmt.Sprintf("%d,", seq)
	}
	id, _ := evented_proto.ProtoToUUID(*eb.Cover.Root)
	return fmt.Sprintf("%s:%s:%s", eb.Cover.Domain, id.String(), strings.TrimSuffix(pages, ","))
}

package support

import (
	"fmt"
	"github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"strings"
)

func StringifyEventBook(eb *evented.EventBook) string {
	var pages string
	for _, page := range eb.Pages {
		seq, _ := GetSequence(page)
		pages += fmt.Sprintf("%d,", seq)
	}
	id, _ := evented_proto.ProtoToUUID(eb.Cover.Root)
	return fmt.Sprintf("%s:%s:%s", eb.Cover.Domain, id.String(), strings.TrimSuffix(pages, ","))
}

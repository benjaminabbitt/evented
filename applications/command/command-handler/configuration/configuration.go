package configuration

const (
	Domain                  string = "domain"
	Name                    string = "name"
	BusinessUrl             string = "business.url"
	Port                    string = "port"
	SyncSaga                string = "sync.sagas"
	SyncProj                string = "sync.projectors"
	Url                     string = "url"
	SnapshotKind            string = "snapshots.kind"
	SnapshotMongoUrl        string = "snapshots.mongodb.url"
	SnapshotMongoName       string = "snapshots.mongodb.name"
	SnapshotMongoCollection string = "snapshots.mongodb.collection"
	TransportKind           string = "transport"
	TransportRabbitKind     string = "transport.rabbitmq"
	TransportRabbitUrl      string = "transport.rabbitmq.url"
	TransportRabbitExchange string = "transport.rabbitmq.exchange"
	TransportNoOpKind       string = "noop"
	RepoKind                string = "events.kind"
	MemoryKind                     = "memory"
	MongoKind                      = "mongodb"
	RepoMongoUrl                   = "events.mongodb.url"
	RepoMongoName                  = "events.mongodb.name"
	RepoMongoCollection            = "events.mongodb.collection"
	ConsulHost                     = "consulHost"
	ConsulPort                     = "consulPort"
)

//type Configuration struct {
//	support.BasicConfigInit
//	Business struct {
//		Url string
//	}
//	Port   uint
//	Domain string
//	Sync   struct {
//		Sagas []struct {
//			Name string
//			Url  string
//		}
//		Projectors []struct {
//			Name string
//			Url  string
//		}
//	}
//	Snapshots struct {
//		Kind    string
//		Mongodb SnapshotStore
//	}
//	Transport struct {
//		Kind     string
//		Rabbitmq struct {
//			Url      string
//			Exchange string
//		}
//	}
//	Events struct {
//		Kind    string
//		Mongodb struct {
//			Url        string
//			Name       string
//			Collection string
//		}
//	}
//}
//
//type SnapshotStore struct {
//	Url        string
//	Name       string
//	Collection string
//}
//
//func (o Configuration) SnapshotStore() SnapshotStore {
//	return o.Snapshots.Mongodb
//}

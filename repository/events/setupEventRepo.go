package events

import (
	"fmt"
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func SetupEventRepo(log *zap.SugaredLogger, errh *evented.ErrLogger) (repo EventRepository, err error) {
	configurationKey := "eventStore"
	typee := viper.GetString("eventstore.type")
	mongodb := "mongodb"
	if typee == mongodb {
		url := viper.GetString(fmt.Sprintf("%s.%s.url", configurationKey, mongodb))
		dbName := viper.GetString(fmt.Sprintf("%s.%s.database", configurationKey, mongodb))
		collectionName := viper.GetString(fmt.Sprintf("%s.%s.collection", configurationKey, mongodb))
		log.Infow("Using MongoDb for Event Store", "url", url, "dbName", dbName)
		repo = mongo.NewMongoClient(url, dbName, collectionName, log, errh)
	}
	return repo, nil
}

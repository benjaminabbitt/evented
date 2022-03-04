package support

import (
	"context"
	"go.uber.org/zap"
)

type ApplicationContext struct {
	context.Context
	log *zap.SugaredLogger
}

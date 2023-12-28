package support

import "go.uber.org/zap"

var (
	Version        string
	VersionLabel   = "version"
	BuildTime      string
	BuildTimeLabel = "build_support build_time"
)

func LogStartup(log *zap.SugaredLogger, appStartup string) {
	log.Infow("Startup: "+appStartup, VersionLabel, Version, BuildTimeLabel, BuildTime)
}

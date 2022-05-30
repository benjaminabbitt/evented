package support

import (
	"fmt"
	"github.com/dsnet/try"
	viperlib "github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"path"
	"strings"
)

const ConfigMgmtType = "CONFIG_MGMT_TYPE"
const ConfigFormat = "CONFIG_FORMAT"
const ConsulHostType = "CONSUL_HOST"
const ConsulKeyType = "CONSUL_KEY"
const LocalConsul = "localhost:8500"
const ConfigPath = "CONFIG_PATH"

const Port = "PORT"

type ConfigType string

const (
	Consul ConfigType = "consul"
	File   ConfigType = "file"
)

const AppNameType = "APP_NAME"
const AppName = "UnnamedEventedApplication"

func Initialize(log *zap.SugaredLogger, viper *viperlib.Viper) (*viperlib.Viper, error) {
	viper.AutomaticEnv()
	viper.SetDefault(ConfigMgmtType, Consul)
	viper.SetDefault(ConsulHostType, LocalConsul)
	viper.SetDefault(AppNameType, AppName)

	configMgmt := viper.GetString(ConfigMgmtType)
	log.Infow("Configuring.", ConfigMgmtType, configMgmt)
	switch ConfigType(configMgmt) {
	case File:
		viper.AddConfigPath(".")
		fullPath := viper.GetString(ConfigPath)
		dir := path.Dir(fullPath)
		file := path.Base(fullPath)
		log.Infow("parsed configpath", "dir", dir, "file", file)
		viper.AddConfigPath(dir)
		fileExtension := path.Ext(fullPath)
		fileName := strings.Replace(file, fileExtension, "", 1)
		log.Infow("file parsed", "fileName", fileName, "fileExtension", fileExtension)
		explicitConfigType := viper.GetString(ConfigFormat)
		if slices.Contains(viperlib.SupportedExts, explicitConfigType) {
			log.Infow("supported explicit configuration", "explicit configuration type", explicitConfigType)
			viper.SetConfigType(explicitConfigType)
		} else if slices.Contains(viperlib.SupportedExts, fileExtension) {
			viper.SetConfigType(fileExtension)
			log.Infow("supported implicit configuration", "implicit configuration type", fileExtension)
		} else {
			log.Warnf("Configuration type could not be determined.  Please set configuration type explicitly via the flag or environment variable \"%s\" or indirecly via the flag or environment variable \"%s\" with a valid extension.  Supported formats are: %+q",
				ConfigFormat, ConfigPath, viperlib.SupportedExts)
		}
		viper.SetConfigName(fileName)
		try.E(viper.ReadInConfig())
		log.Infow("Configuration set", "configuration", fmt.Sprintf("%+v", viper.AllSettings()))
		log.Infow("Read file.", "Proof", viper.GetString("proof"))
		return viper, nil
	case Consul:
		consulHost := viper.GetString(ConsulHostType)
		consulKey := viper.GetString(ConsulKeyType)
		log.Infow("Attempting to fetch configuration", "provider", configMgmt, "host", consulHost, "key", consulKey)
		try.E(viper.AddRemoteProvider("consul", consulHost, consulKey))
		try.E(viper.ReadRemoteConfig())
		log.Infow("Configuration set", "configuration", fmt.Sprintf("%+v", viper.AllSettings()))
		log.Infow("Read consul.", "Proof", viper.GetString("proof"))
		return viper, nil
	default:
		return viper, fmt.Errorf("configuration provider %s not matched with valid types", configMgmt)
	}
}

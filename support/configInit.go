package support

import (
	"fmt"
	"github.com/dsnet/try"
	viperlib "github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"path/filepath"
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

const (
	yaml string = "yaml"
	json string = "json"
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
		dir, file := filepath.Split(viper.GetString(ConfigPath))
		viper.AddConfigPath(dir)
		f := strings.Split(file, ".")
		fileName := f[0]
		fileExtension := ""
		if len(f) > 1 {
			fileExtension = f[1]
		}
		explicitConfigType := viper.GetString(ConfigFormat)
		if slices.Contains(viperlib.SupportedExts, explicitConfigType) {
			viper.SetConfigType(explicitConfigType)
		} else if slices.Contains(viperlib.SupportedExts, fileExtension) {
			viper.SetConfigType(fileExtension)
		} else {
			log.Warnf("Configuration type could not be determined.  Please set configuration type explicitly via the flag or environment variable \"%s\" or indirecly via the flag or environment variable \"%s\" with a valid extension.  Supported formats are: %+q",
				ConfigFormat, ConfigPath, viperlib.SupportedExts)
		}
		//TODO: improve this to handle multiple periods in file names e.g. foo.backup.yaml
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

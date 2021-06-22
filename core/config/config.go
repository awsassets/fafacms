package config

import (
	"encoding/json"
	"github.com/hunterhug/fafacms/core/util/kv"
	"github.com/hunterhug/fafacms/core/util/mail"
	"github.com/hunterhug/fafacms/core/util/oss"
	"github.com/hunterhug/fafacms/core/util/rdb"
)

var (
	//  Global config!
	FaFaConfig *Config
)

type Config struct {
	DefaultConfig MyConfig       `yaml:"DefaultConfig"` // default config
	OssConfig     oss.Key        `yaml:"OssConfig"`     // oss like aws s3
	DbConfig      rdb.MyDbConfig `yaml:"DbConfig"`      // mysql config
	SessionConfig kv.MyRedisConf `yaml:"SessionConfig"` // redis config for user session
	MailConfig    mail.Sender    `yaml:"MailConfig"`    // email config
}

// Some especial my config
type MyConfig struct {
	WebPort       string `yaml:"WebPort"`
	LogPath       string `yaml:"LogPath"`
	StoragePath   string `yaml:"StoragePath"`
	LogDebug      bool   `yaml:"LogDebug"`
	StorageOss    bool   `yaml:"StorageOss"`
	CloseRegister bool   `yaml:"CloseRegister"`
}

// Let the config struct to json file, just for test
func JsonOutConfig(config Config) (string, error) {
	raw, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	back := string(raw)
	return back, nil
}

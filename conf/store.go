package conf

import (
	"errors"
	"flag"
	"os"
	"path/filepath"

	"github.com/ii64/go-binder/binder"
	"github.com/ii64/go-binder/binder/ext/json"
	"github.com/ii64/go-binder/binder/ext/toml"
	"github.com/ii64/go-binder/binder/ext/yaml"
)

var ConfigPath = getConfigPath()

func _initStore() {
	switch filepath.Ext(ConfigPath) {
	case ".yaml", ".yml":
		binder.LoadConfig = yaml.LoadConfig(ConfigPath)
		binder.SaveConfig = yaml.SaveConfig(ConfigPath, 2)
	case ".toml":
		binder.LoadConfig = toml.LoadConfig(ConfigPath)
		binder.SaveConfig = toml.SaveConfig(ConfigPath, "  ")
	default: // fallback json
		binder.LoadConfig = json.LoadConfig(ConfigPath)
		binder.SaveConfig = json.SaveConfig(ConfigPath, "  ")
	}
	binder.SaveOnClose = false
}

func Init() {
	var err error
	_initStore()
	err = binder.Init()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = binder.Save()
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	flag.Parse()

	binder.In()
}

func getConfigPath() (ret string) {
	ret = os.Getenv("CONFIG_PATH")
	if ret == "" {
		ret = ".dlvman.toml"
	}
	return
}

func Close() error {
	return binder.Close()
}

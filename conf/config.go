package conf

import (
	"github.com/ii64/go-binder/binder"
)

type MultiValue string

type Config struct {
	Addr        string   `argx:"addr" environ:"ADDR" usage:"Listen address"`
	Program     string   `argx:"p" usage:"Debugger program"`
	ProcessArgs []string `argx:"pa" usage:"Debugger program arguments"`
	Debug       bool     `argx:"debug" environ:"DEBUG" usage:"Debug logging"`
}

// dlv dap --check-go-version=false --listen=127.0.0.1:7456 --log=true --log-output=debugger,debuglineerr,gdbwire,lldout,rpc

var frozenDefault = &Config{
	Addr:    "0.0.0.0:7456",
	Program: "dlv",
	ProcessArgs: []string{
		"dap",
		"--check-go-version=false",
		"--log-dest=3",
	},
	Debug: false,
}

var Default = frozenDefault.Clone()

var _ = func() struct{} {
	// init
	binder.BindArgsConf(Default, "conf")
	return struct{}{}
}()

// Validate configuration
func (c *Config) Validate() error {
	return nil
}

// LinkIn called by go-binder.
func (c *Config) LinkIn() {
	err := c.Validate()
	if err != nil {
		panic(err)
	}
}

// Clone is not implementing deep-copy
func (c *Config) Clone() *Config {
	confLocal := *c
	return &confLocal
}

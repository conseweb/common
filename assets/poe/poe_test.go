package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestYaml(t *testing.T) {
	viper.SetConfigName("poe")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	gopath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {
		cfgpath := filepath.Join(p, "src", "github.com/conseweb/common/assets/poe")
		viper.AddConfigPath(cfgpath)
	}
	e := viper.ReadInConfig()
	if e != nil {
		t.Error(e)
	}
	var config []ConfigSystem
	e = viper.UnmarshalKey("system", &config)
	if e != nil {
		t.Error(e)
	}
	t.Log(config)
}
func TestFunc(t *testing.T) {
	sys, e := configSystem("base")
	if e != nil {
		t.Log(e)
	}
	t.Log(sys)

	api := cryptoStrategyMap["default"]("base")
	d, e := api.algorithm([]byte("xiebo"))
	if e != nil {
		t.Log(e)
	}
	t.Log(string(d))
}

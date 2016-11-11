package main

import (
	"os"
	"path/filepath"
	"testing"

	//"github.com/hyperledger/fabric/core/chaincode/shim"
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

//func TestFunc(t *testing.T) {
//	sys, e := configSystem("base")
//	if e != nil {
//		t.Log(e)
//	}
//	api := cryptoStrategyMap["default"](sys)
//	var array []string = []string{"xiebo", "xiebo"}
//	var list []string
//	for _, v := range array {
//		d, e := api.algorithm([]byte(v))
//		if e != nil {
//			t.Log(e)
//		}
//		list = append(list, string(d))
//	}
//	if list[0] == list[1] {
//		t.Log("some")
//	}
//	logger.Info(list[0])
//	logger.Info(list[1])
//}

//func TestInvoke(t *testing.T) {
//	var args []string = []string{"base", "wangchuanjian", "zhaoming"}
//	cc := new(PoeService)
//	stub := shim.NewMockStub("ex05", cc)
//	_, err := stub.MockInvoke("1", "register", args)
//	if err != nil {
//		t.Error("Invoke", args, "failed", err)
//		t.FailNow()
//	}
//}

//func TestQuery(t *testing.T) {
//	var args []string = []string{"base", "wangchuanjian", "zhaoming"}
//	cc := new(PoeService)
//	stub := shim.NewMockStub("ex05", cc)
//	data, err := stub.MockQuery("existence", args)
//	if err != nil {
//		t.Error("Invoke", args, "failed", err)
//		t.FailNow()
//	}
//	t.Log(string(data))
//}

package main

import (
	"github.com/conseweb/common/config"
	"github.com/hyperledger/fabric/flogging"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

// 默认值
const (
	// 默认加密策略主哈希函数默认值
	CRYPTO_STRATEGY_DEFAULT_MAIN_HASH_VALUE = "sha3_256"
	// 默认加密策略副哈希函数默认值
	CRYPTO_STRATEGY_DEFAULT_SUB_HASH_VALUE = "md5"
)

// 日志记录器
var logger = logging.MustGetLogger("poe")

// 初始化函数
func init() {
	// 加载配置文件
	e := config.LoadConfig("", "poe", "github.com/conseweb/common/assets/poe")
	if e != nil {
		logger.Errorf("start PoeService return error: %v, exiting\n", e)
		return
	}
	// 初始化日志
	flogging.LoggingInit("poe")
}

// 根据系统名称获取配置节点
func configSystem(sysName string) (cfgSys *ConfigSystem, e error) {
	var list []*ConfigSystem
	if e = viper.UnmarshalKey("system", &list); e != nil {
		return nil, e
	}
	for _, v := range list {
		if v.Name == sysName {
			cfgSys = v
			break
		}
	}
	if cfgSys == nil {
		return configSystem("base")
	}
	return cfgSys, nil
}

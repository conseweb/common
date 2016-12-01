package main

import (
	"github.com/conseweb/common/config"
	"github.com/hyperledger/fabric/flogging"
	"github.com/spf13/viper"
)

// 默认值
const (
	// 默认加密策略主哈希函数默认值
	CRYPTO_STRATEGY_DEFAULT_MAIN_HASH_VALUE = "sha3_256"
	// 默认加密策略副哈希函数默认值
	CRYPTO_STRATEGY_DEFAULT_SUB_HASH_VALUE = "md5"
)

var (
	cfgmaps map[string]*ConfigSystem
)

// 初始化函数
func init() {
	// 加载配置文件
	e := config.LoadConfig("POE", "poe", "github.com/conseweb/common/assets/poe")
	if e != nil {
		logger.Errorf("start PoeService return error: %v, exiting\n", e)
		return
	}
	// 初始化日志
	flogging.LoggingInit("poe")

	cfgmaps = loadConfigSystem()
}

func loadConfigSystem() map[string]*ConfigSystem {
	var list []*ConfigSystem
	if err := viper.UnmarshalKey("system", &list); err != nil {
		logger.Panic(err)
	}

	cfgs := make(map[string]*ConfigSystem)
	for _, cfg := range list {
		cfgs[cfg.Name] = cfg
	}

	if _, ok := cfgs["base"]; !ok {
		logger.Panic("must has a base config system")
	}

	return cfgs
}

// 根据系统名称获取配置节点
func configSystem(sysName string) *ConfigSystem {
	if cfg, ok := cfgmaps[sysName]; ok {
		return cfg
	}

	return cfgmaps["base"]
}

package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/conseweb/common/crypto"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

type (
	// 加密策略构造函数
	cryptoStrategyFunc func(*ConfigSystem) cryptoStrategy
	// 加密策略接口
	cryptoStrategy interface {
		// 加密算法
		algorithm(d []byte) ([]byte, error)
	}
	// 加密策略默认实现
	cryptoStrategyDefault struct {
		// 主哈希函数名称
		mainHashName string
		// 副哈希函数名称
		subHashName string
	}
)

var (
	// 加密算法集合
	cryptoFuncMap map[string]func() hash.Hash = map[string]func() hash.Hash{
		"ripemd160": ripemd160.New,
		"sha3_224":  sha3.New224,
		"sha3_256":  sha3.New256,
		"sha3_384":  sha3.New384,
		"sha3_512":  sha3.New512,
		"md5":       md5.New,
	}
	// 加密策略集合
	cryptoStrategyMap map[string]cryptoStrategyFunc = map[string]cryptoStrategyFunc{
		"default": NewCryptoStrategyDefault,
	}
)

// 实例化加密策略默认实现
func NewCryptoStrategyDefault(cfg *ConfigSystem) cryptoStrategy {
	cryStrDef := cryptoStrategyDefault{}
	cryStrDef.mainHashName = cfg.CryptoStrategyDefault.MainHash
	cryStrDef.subHashName = cfg.CryptoStrategyDefault.SubHash
	return &cryStrDef
}

// 加密算法
// 主哈希值+副哈希值+长度
func (this *cryptoStrategyDefault) algorithm(d []byte) ([]byte, error) {
	mainHash, subHash := cryptoFuncMap[this.mainHashName](), cryptoFuncMap[this.subHashName]()
	return bytes.NewBufferString(fmt.Sprintf("%x%x%d", crypto.Hash(mainHash, d), crypto.Hash(subHash, d), len(d))).Bytes(), nil
}

// 计算键值哈希
// cfgSys: 业务系统配置
// key: 需要加密的字符串
func hashKey(cfgSys *ConfigSystem, key string) (string, error) {
	var (
		strategy cryptoStrategyFunc
		data     []byte
		e        error
	)
	if strategy = cryptoStrategyMap[cfgSys.CryptoStrategy]; strategy == nil {
		strategy = cryptoStrategyMap["default"]
	}
	if data, e = strategy(cfgSys).algorithm([]byte(key)); e != nil {
		return "", e
	}
	return hex.EncodeToString(data), nil
}

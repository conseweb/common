package main

import (
	"crypto/md5"
	"hash"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
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

// 实例化加密策略默认实现
func NewCryptoStrategyDefault(cfg *ConfigSystem) cryptoStrategy {
	cryStrDef := cryptoStrategyDefault{}
	cryStrDef.mainHashName = cfg.CryptoStrategyDefault.MainHash
	cryStrDef.subHashName = cfg.CryptoStrategyDefault.SubHash
	return &cryStrDef
}

// 加密算法
// 主哈希值+副哈希值+长度
func (this *cryptoStrategyDefault) algorithm(d []byte) (r []byte, e error) {
	mainHash, subHash := cryptoFuncMap[this.mainHashName], cryptoFuncMap[this.subHashName]
	if _, e = mainHash().Write(d); e != nil {
		return nil, e
	}
	if _, e = subHash().Write(d); e != nil {
		return nil, e
	}
	r = []byte(string(mainHash().Sum(nil)) + string(subHash().Sum(nil)) + string(len(d)))
	return r, nil
}

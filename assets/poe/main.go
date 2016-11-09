package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"hash"
	"strings"

	"github.com/conseweb/common/config"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/flogging"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// 日志记录器
var logger = logging.MustGetLogger("poe")

// 哈稀函数标记
const (
	HASH_FUNC_FLAGS_RIPEMD160 = "RIPEMD160"
	HASH_FUNC_FLAGS_SHA3_224  = "SHA3_224"
	HASH_FUNC_FLAGS_SHA3_256  = "SHA3_256"
	HASH_FUNC_FLAGS_SHA3_384  = "SHA3_384"
	HASH_FUNC_FLAGS_SHA3_512  = "SHA3_512"
	HASH_FUNC_FLAGS_MD5       = "MD5"
)

// Proof of Existence Service（存在性证明服务）
type PoeService struct{}

func (this *PoeService) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "deploy" {
		return nil, errors.New("invalid function name, 'deploy' only")
	}
	logger.Info("deploy Lepuscoin successfully")
	return nil, nil
}

func (this *PoeService) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case IF_REGISTER:
		return this.register(stub, args)
	default:
		return nil, errors.New("unsupported operation")
	}
}

func (this *PoeService) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case QF_EXISTENCE:
		return this.existence(stub, args)
	default:
		return nil, errors.New("unsupported operation")
	}
}

const (
	IF_REGISTER  = "register"
	QF_EXISTENCE = "existence"
)

//注册键值
func (this *PoeService) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// 验证非空
	if len(args) == 0 {
		return nil, errors.New("func <register> Parameter is not valid,Cannot be empty or contain null characters")
	}
	//logger.Warningf("<register> parameters :%s", args)
	// 循环存入
	for _, key := range args {
		if len(strings.TrimSpace(key)) > 0 {
			hkey, e := hashKey(key)
			if e != nil {
				return nil, e
			}
			e = stub.PutState(hkey, []byte{1})
			if e != nil {
				return nil, e
			}
		}
	}
	return nil, nil
}

//检索键值是否存在
func (this *PoeService) existence(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var list []QueryResult
	// 验证非空
	if len(args) == 0 {
		return nil, errors.New("func <existence> Parameter is not valid,Cannot be empty or contain null characters")
	}
	//logger.Warningf("<existence> parameters :%s", args)
	for _, key := range args {
		hkey, e := hashKey(key)
		if e != nil {
			return nil, errors.New("func <existence> error:" + e.Error())
		}
		d, e := stub.GetState(hkey)
		if e != nil {
			return nil, errors.New("func <existence> error:" + e.Error())
		}
		if d != nil && len(d) > 0 {
			m := QueryResult{}
			m.Key = key
			m.Exist = true
			list = append(list, m)
		}
	}
	data, e := json.Marshal(&list)
	if e != nil {
		return nil, errors.New("func <existence> error:" + e.Error())
	}
	return data, nil
}

// 获取传入key 的哈希值
// hash1+hash2+len(key)
func hashKey(key string) (string, error) {
	mainHashConf := viper.GetString("hash_func.main_hash")
	mainHash := hashDevice(mainHashConf)
	_, e := mainHash.Write([]byte(key))
	if e != nil {
		return "", e
	}
	subHashConf := viper.GetString("hash_func.sub_hash")
	subHash := hashDevice(subHashConf)
	mainStr := string(mainHash.Sum(nil))
	_, e = subHash.Write([]byte(key))
	if e != nil {
		return "", e
	}
	subStr := string(subHash.Sum(nil))
	return mainStr + subStr + string(len(key)), nil
}

// 根据哈希函数标记构造具体的哈希实例
func hashDevice(flags string) hash.Hash {
	switch flags {
	case HASH_FUNC_FLAGS_RIPEMD160:
		return ripemd160.New()
	case HASH_FUNC_FLAGS_SHA3_224:
		return sha3.New224()
	case HASH_FUNC_FLAGS_SHA3_256:
		return sha3.New256()
	case HASH_FUNC_FLAGS_SHA3_384:
		return sha3.New384()
	case HASH_FUNC_FLAGS_SHA3_512:
		return sha3.New512()
	case HASH_FUNC_FLAGS_MD5:
		return md5.New()
	default:
		return md5.New()
	}
}

func main() {
	e := config.LoadConfig("", "poe", "github.com/conseweb/common/assets/poe")
	if e != nil {
		logger.Errorf("start PoeService return error: %v, exiting\n", e)
		return
	}
	flogging.LoggingInit("poe")
	if e = shim.Start(new(PoeService)); e != nil {
		logger.Errorf("start PoeService return error: %v, exiting\n", e)
	}
}

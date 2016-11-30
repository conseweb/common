package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/spf13/viper"
)

var (
	poeFuncMap map[string]poeFunc = map[string]poeFunc{
		"register":  register,
		"existence": existence,
	}
	inSlice = func(v string, sl []string) bool {
		for _, vv := range sl {
			if vv == v {
				return true
			}
		}
		return false
	}
)

type (
	// Proof of Existence Service（存在性证明服务）
	PoeService struct{}
	poeFunc    func(shim.ChaincodeStubInterface, []string) ([]byte, error)
)

func (this *PoeService) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "deploy" {
		return nil, errors.New("invalid function name, 'deploy' only")
	}
	logger.Debug("deploy Lepuscoin successfully")
	return nil, nil
}

func (this *PoeService) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if fn, ok := poeFuncMap[function]; ok {
		return fn(stub, args)
	} else {
		logger.Debug("func <Invoke> unsupported operation")
		return nil, errors.New("unsupported operation")
	}

}

func (this *PoeService) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if fn, ok := poeFuncMap[function]; ok {
		return fn(stub, args)
	} else {
		logger.Debug("func <Invoke> unsupported operation")
		return nil, errors.New("unsupported operation")
	}

}

// 注册键值
// args : 参数第一个元素为业务系统标记，默认为base
func register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var (
		cfgSys  *ConfigSystem
		sysName string = strings.TrimSpace(args[0])
		hkey    string
		array   []string
		index   int = 0
		e       error
	)
	// 验证非空
	if len(args) == 0 {
		logger.Debug("func <register> Parameter is not valid,Cannot be empty or contain null characters")
		return nil, errors.New("func <register> Parameter is not valid,Cannot be empty or contain null characters")
	}
	if len(sysName) > 0 && inSlice(sysName, strings.Split(viper.GetString("system_list"), ",")) {
		index = 1
	}
	logger.Infof("func <register> index : %v", index)
	if cfgSys, e = configSystem(sysName); e != nil {
		logger.Debugf("func <register> error : %v", e)
		return nil, errors.New("func <register> error:" + e.Error())
	}
	for i := index; i < len(args); i++ {
		if len(strings.TrimSpace(args[i])) == 0 {
			continue
		}
		if hkey, e = hashKey(cfgSys, args[i]); e != nil {
			logger.Debugf("func <register> error : %v", e)
			return nil, errors.New("func <register> error:" + e.Error())
		}
		if e = stub.PutState(hkey, []byte{1}); e != nil {
			logger.Debugf("func <register> error : %v", e)
			return nil, errors.New("func <register> error:" + e.Error())
		}
		array = append(array, hkey)
	}
	if len(array) > 0 {
		if e = stub.SetEvent("invoke_completed", []byte(strings.Join(array, ","))); e != nil {
			logger.Debugf("func <register> error : %v", e)
			return nil, errors.New("func <register> error:" + e.Error())
		}
		logger.Debug("func <register> set event invoke_completed")
	}
	logger.Info("func <register> over")
	return nil, nil
}

// 检索键值是否存在
// args : 参数第一个元素为业务系统标记，默认为base
func existence(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var (
		list    []QueryResult
		cfgSys  *ConfigSystem
		sysName string = strings.TrimSpace(args[0])
		hkey    string
		data    []byte
		index   int = 0
		e       error
	)
	// 验证非空
	if len(args) == 0 {
		logger.Debug("func <existence> Parameter is not valid,Cannot be empty or contain null characters")
		return nil, errors.New("func <existence> Parameter is not valid,Cannot be empty or contain null characters")
	}
	if len(sysName) > 0 && inSlice(sysName, strings.Split(viper.GetString("system_list"), ",")) {
		index = 1
	}
	logger.Debugf("func <existence> index : %v", index)
	if cfgSys, e = configSystem(sysName); e != nil {
		logger.Debugf("func <existence> error : %v", e)
		return nil, e
	}
	for i := index; i < len(args); i++ {
		if len(strings.TrimSpace(args[i])) == 0 {
			continue
		}
		if hkey, e = hashKey(cfgSys, args[i]); e != nil {
			logger.Debugf("func <existence> error : %v", e)
			return nil, errors.New("func <existence> error:" + e.Error())
		}
		if data, e = stub.GetState(hkey); e != nil {
			logger.Debugf("func <existence> error : %v", e)
			return nil, errors.New("func <existence> error:" + e.Error())
		}
		if len(data) > 0 {
			m := QueryResult{}
			m.Key = args[i]
			m.HashKey = hkey
			m.Exist = true
			list = append(list, m)
		}
	}
	if data, e = json.Marshal(&list); e != nil {
		logger.Debugf("func <existence> error : %v", e)
		return nil, errors.New("func <existence> error:" + e.Error())
	}
	logger.Info("func <existence> over")
	return data, nil
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

// 程序入口
func main() {
	if e := shim.Start(new(PoeService)); e != nil {
		logger.Errorf("start PoeService return error: %v, exiting\n", e)
	}
}

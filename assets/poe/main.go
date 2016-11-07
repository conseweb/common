package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("PoeService")

// Proof of Existence Service
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
	var (
		param string = args[0]
		array []string
	)
	// 验证非空
	if len(strings.TrimSpace(param)) == 0 {
		return nil, errors.New("func <register> Parameter is not valid,Cannot be empty or contain null characters")
	}
	array = strings.Split(param, ",")
	// 循环存入
	for _, key := range array {
		m := StoreEntity{}
		m.Key = hashMd5(key)
		m.Length = int64(len(key))
		m.Value = []byte(key)
		data, e := json.Marshal(&m)
		if e != nil {
			logger.Errorf("func <register> error:%s", e.Error())
			break
		}
		e = stub.PutState(m.Key, data)
		if e != nil {
			logger.Errorf("func <register> error:%s", e.Error())
			break
		}
	}
	return nil, nil
}

//检索键值是否存在
func (this *PoeService) existence(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var (
		param string = args[0]
		array []string
		list  []QueryEntity
	)
	// 验证非空
	if len(strings.TrimSpace(param)) == 0 {
		return nil, errors.New("func <existence> Parameter is not valid,Cannot be empty or contain null characters")
	}
	array = strings.Split(param, ",")
	for _, key := range array {
		data, e := stub.GetState(hashMd5(key))
		if e != nil {
			return nil, errors.New("func <existence> error:" + e.Error())
		}
		if data != nil && len(data) > 0 {
			m := QueryEntity{}
			m.Key = key
			m.IsExist = true
			list = append(list, m)
		}
	}
	return nil, nil
}

// MD5 Hash
func hashMd5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	cpiherText := hash.Sum(nil)
	hexText := make([]byte, 32)
	hex.Encode(hexText, cpiherText)
	return string(hexText)
}

func main() {
	if err := shim.Start(new(PoeService)); err != nil {
		logger.Errorf("start Lepuscoin return error: %v, exiting\n", err)
	}
}

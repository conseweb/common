/*
Copyright Mojing Inc. 2016 All Rights Reserved.
Written by mint.zhao.chiu@gmail.com. github.com: https://www.github.com/mintzhao

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"strings"
	"reflect"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type (
	// Proof of Existence Service（存在性证明服务）
	PoeService struct{}
)

var (
	err_unsupported_operation = fmt.Errorf("unsupported operation")
	err_invalid_param         = fmt.Errorf("invalid param")
	filledvalue               = []byte("1")
)

func (this *PoeService) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "deploy" {
		return nil, err_unsupported_operation
	}

	logger.Debug("deploy poe successfully")
	return nil, nil
}

func (this *PoeService) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "register" {
		logger.Debug("unsupported operation %s", function)
		return nil, err_unsupported_operation
	}

	return register(stub, args)
}

func (this *PoeService) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "existence" {
		logger.Debug("unsupported operation %s", function)
		return nil, err_unsupported_operation
	}

	return existence(stub, args)
}

// 注册键值
// args : 参数第一个元素为业务系统标记，默认为base
func register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// 验证非空
	if len(args) <= 1 {
		logger.Warningf("invalid param: %v", args)
		return nil, err_invalid_param
	}

	cfgSys := configSystem(strings.TrimSpace(args[0]))
	results := make(map[string]string)

	for _, doc := range args[1:] {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		hkey, err := hashKey(cfgSys, doc)
		if err != nil {
			logger.Errorf("generate hash key return error: %v", err)
			return nil, err
		}

		if err := stub.PutState(hkey, []byte("1")); err != nil {
			logger.Errorf("put state into blockchain return err: %v", err)
			return nil, err
		}

		results[doc] = hkey
	}

	if len(results) > 0 {
		resultsBytes, err := json.Marshal(results)
		if err != nil {
			logger.Errorf("json.marshal return error: %v", err)
			return nil, err
		}
		if err := stub.SetEvent("invoke_completed", resultsBytes); err != nil {
			logger.Errorf("set chaincode event error: %v", err)
			return nil, err
		}
		logger.Debug("set chaincode event invoke_completed")
	}
	return nil, nil
}

// 检索键值是否存在
// args : 参数第一个元素为业务系统标记，默认为base
func existence(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// 验证非空
	if len(args) <= 1 {
		logger.Warningf("invalid param: %v", args)
		return nil, err_invalid_param
	}

	cfgSys := configSystem(strings.TrimSpace(args[0]))
	results := make(map[string]*QueryResult)

	for _, doc := range args[1:] {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		hkey, err := hashKey(cfgSys, doc)
		if err != nil {
			logger.Errorf("generate hash key return error: %v", err)
			return nil, err
		}

		result := &QueryResult{
			Key:     doc,
			HashKey: hkey,
			Exist:   false,
		}
		val, err := stub.GetState(hkey)
		if err == nil && reflect.DeepEqual(val, filledvalue) {
			result.Exist = true
		}

		results[doc] = result
	}

	resultsBytes, err := json.Marshal(results)
	if err != nil {
		logger.Errorf("json.marshal return error: %v", err)
		return nil, err
	}

	return resultsBytes, nil
}

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
package client

import (
	"bytes"
	"errors"
	"github.com/conseweb/common/fabricGoSDK/client"
	"github.com/conseweb/common/fabricGoSDK/client/chaincode"
	"github.com/conseweb/common/fabricGoSDK/models"
	pb "github.com/conseweb/common/protos"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/op/go-logging"
)

var (
	logger        = logging.MustGetLogger("lepuscoinClient")
	jsonrpc       = "2.0"
	deployMethod  = "deploy"
	invokeMethod  = "invoke"
	chaincodeType = "GOLANG"
)

// LepuscoinClient
type LepuscoinClient struct {
	*client.FabricGoSDK
	chaincodeID string
	requestID   int64
}

type ExampleClient struct {
	*client.FabricGoSDK
	chaincodeID string
	requestID   int64
}

// NewLepuscoinClient return a LepuscoinClient based on given host & basePath
func NewLepuscoinClient(host, basePath string, formats strfmt.Registry) *LepuscoinClient {
	cli := new(LepuscoinClient)
	transport := httptransport.New(host, basePath, []string{"http", "https"})
	cli.FabricGoSDK = client.New(transport, formats)
	cli.requestID = 1

	return cli
}

// NewDefaultLepuscoinClient return a LepuscoinClient using default host 127.0.0.1:7050
func NewDefaultLepuscoinClient() *LepuscoinClient {
	cli := new(LepuscoinClient)
	cli.FabricGoSDK = client.Default
	cli.requestID = 1

	return cli
}

// Deploy deploy lepuscoin chaincode into vp
func (cli *LepuscoinClient) Deploy(name, path, secureContext string, confidentialityLevel models.ConfidentialityLevel) (string, error) {
	opParams := chaincode.NewChaincodeOpParams().WithChaincodeOpPayload(&models.ChaincodeOpPayload{
		ID:      &cli.requestID,
		Jsonrpc: &jsonrpc,
		Method:  &deployMethod,
		Params: &models.ChaincodeSpec{
			ChaincodeID: &models.ChaincodeID{
				Name: name,
				Path: path,
			},
			ConfidentialityLevel: confidentialityLevel,
			SecureContext:        secureContext,
			CtorMsg: &models.ChaincodeInput{
				Args: []string{"deploy"},
			},
			Type: &chaincodeType,
		},
	})

	opOK, err := cli.Chaincode.ChaincodeOp(opParams)
	if err != nil {
		logger.Errorf("deploy lepuscoin return error: %v", err)
		return "", err
	}

	logger.Debugf("Deploy response: %s", opOK.Error())
	deployId := *opOK.Payload.Result.Message
	cli.chaincodeID = deployId
	cli.requestID += 1

	return deployId, nil
}

// InvokeCoinbase
func (cli *LepuscoinClient) InvokeCoinbase(confidentialityLevel models.ConfidentialityLevel, secureContext string, tx *pb.TX) error {
	if cli.chaincodeID == "" {
		return errors.New("no chaincode specified")
	}

	// tx handle
	txBytes, err := tx.Bytes()
	if err != nil {
		logger.Errorf("tx marshal return error: %v", err)
		return err
	}

	opParams := chaincode.NewChaincodeOpParams().WithChaincodeOpPayload(&models.ChaincodeOpPayload{
		ID:      &cli.requestID,
		Jsonrpc: &jsonrpc,
		Method:  &invokeMethod,
		Params: &models.ChaincodeSpec{
			ChaincodeID: &models.ChaincodeID{
				Name: cli.chaincodeID,
			},
			ConfidentialityLevel: confidentialityLevel,
			SecureContext:        secureContext,
			CtorMsg: &models.ChaincodeInput{
				Args: []string{"invoke_coinbase", bytes.NewBuffer(txBytes).String()},
			},
			Type: &chaincodeType,
		},
	})

	opOK, err := cli.Chaincode.ChaincodeOp(opParams)
	if err != nil {
		logger.Errorf("invoke lepuscoin return error: %v", err)
		return err
	}

	logger.Debugf("invoke response: %s", opOK.Error())
	cli.requestID += 1

	return nil
}

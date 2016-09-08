/*
Copyright Mojing Inc. 2016 All Rights Reserved.

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
package coin

import (
	"encoding/base64"
	pb "github.com/conseweb/common/protos"
)

func (coin *Lepuscoin) coinbase(store Store, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgs
	}

	txDataBase64 := args[0]
	txData, err := base64.StdEncoding.DecodeString(txDataBase64)
	if err != nil {
		logger.Errorf("Decoding base64 error: %v\n", err)
		return nil, err
	}

	tx, err := pb.ParseTXBytes(txData)
	if err != nil {
		logger.Errorf("Unmarshal tx bytes error: %v\n", err)
		return nil, err
	}

	// utxo
	utxo := MakeUTXO(store)
	execResult, err := utxo.Execute(tx)
	if err != nil {
		logger.Errorf("execute coinbase tx return error: %v", err)
		return nil, err
	}

	if !execResult.IsCoinbase {
		return nil, ErrMustCoinbase
	}

	// account
	account := MakeAccount(store)
	if err := account.Coinbase(tx); err != nil {
		logger.Errorf("Error execute account model coinbase: %v", err)
		return nil, err
	}

	logger.Debugf("coinbase execute result: %+v", execResult)
	return execResult.Bytes()
}

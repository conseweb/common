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
	"fmt"
)

func (coin *Lepuscoin) coinbase(store Store, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid args")
	}

	// utxo
	utxo := MakeUTXO(store)
	txDataBase64 := args[0]
	txData, err := base64.StdEncoding.DecodeString(txDataBase64)
	if err != nil {
		return nil, fmt.Errorf("Error decoding TX as base64:  %s", err)
	}
	execResult, err := utxo.Execute(txData)
	if err != nil {
		logger.Errorf("execute coinbase tx return error: %v", err)
		return nil, err
	}

	if !execResult.IsCoinbase {
		return nil, fmt.Errorf("the Tx must be a coinbase tx")
	}

	// account
	account := MakeAccount(store)
	if err := account.Coinbase(execResult.SumCurrentOutputs, "xxxxxxxxxxx"); err != nil {
		logger.Errorf("account model execute coinbase return error: %v", err)
		return nil, err
	}

	return nil, nil
}

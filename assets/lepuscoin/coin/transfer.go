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

// 1. 转账是不会产生新货币的
// 2. TODO 每个tx有fee,会发送给区块的生成者,但是怎么给他呢
// 3. 在进行utxo执行前,应该先验证txin是否拥有足够的余额
// 4. 需要验证输入和输出+fee是否一致

func (coin *Lepuscoin) transfer(store Store, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgs
	}

	// parse tx
	txDataBase64 := args[0]
	txData, err := base64.StdEncoding.DecodeString(txDataBase64)
	if err != nil {
		return nil, err
	}

	tx, err := pb.ParseTXBytes(txData)
	if err != nil {
		return nil, err
	}

	// account
	account := MakeAccount(store)
	if err := account.Transfer(tx); err != nil {
		logger.Errorf("Error execute account model transfer: %v", err)
		return nil, err
	}

	// utxo
	utxo := MakeUTXO(store)
	execResult, err := utxo.Execute(tx)
	if err != nil {
		logger.Errorf("execute tx transfer return error: %v", err)
		return nil, err
	}

	if execResult.IsCoinbase {
		return nil, ErrCantCoinbase
	}

	// current outputs must less than prior outputs
	if execResult.SumCurrentOutputs > execResult.SumPriorOutputs {
		return nil, ErrTxOutNotLogic
	}

	// one of transfer main point is in == out, no coin mined, no coin lose
	if execResult.SumCurrentOutputs + tx.Fee != execResult.SumPriorOutputs {
		return nil, ErrTxInOutNotBalance
	}

	return execResult.Bytes()
}

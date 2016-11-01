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
package coin

import (
	"encoding/base64"

	pb "github.com/conseweb/common/assets/lepuscoin/protos"
)

func (coin *Lepuscoin) coinbase(store Store, args []string) ([]byte, error) {
	if len(args) != 1 || args[0] == "" {
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

	// check tx whether is a coinbase tx
	if len(tx.Txin) > 0 {
		return nil, ErrMustCoinbase
	}

	txhash := tx.TxHash()
	execResult := &pb.ExecResult{}

	coinInfo, err := store.GetCoinInfo()
	if err != nil {
		logger.Errorf("Error get coin info: %v", err)
		return nil, err
	}

	// Loop through outputs
	for index, output := range tx.Txout {
		if output.Addr == "" {
			return nil, ErrInvalidLepuscoinTX
		}

		outerAccount, err := store.GetAccount(output.Addr)
		if err != nil {
			logger.Warningf("account[%s] is not existed, creating one...", output.Addr)

			outerAccount = new(pb.Account)
			outerAccount.Addr = output.Addr
			outerAccount.Txouts = make(map[string]*pb.TX_TXOUT)

			coinInfo.AccountTotal += 1
		}
		if outerAccount.Txouts == nil || len(outerAccount.Txouts) == 0 {
			outerAccount.Txouts = make(map[string]*pb.TX_TXOUT)
		}

		currKey := &Key{TxHashAsHex: txhash, TxIndex: uint32(index)}
		if _, ok := outerAccount.Txouts[currKey.String()]; ok {
			return nil, ErrCollisionTxOut
		}

		// store tx out into account
		outerAccount.Txouts[currKey.String()] = output
		outerAccount.Balance += output.Value
		if err := store.PutAccount(outerAccount); err != nil {
			logger.Errorf("Error update account: %v, account info: %+v", err, outerAccount)
			return nil, err
		}
		logger.Debugf("put tx output %s:%v", currKey.String(), output)

		// coin stat
		coinInfo.CoinTotal += output.Value
		coinInfo.TxoutTotal += 1
		execResult.SumCurrentOutputs += output.Value
	}

	if err := store.PutTx(tx); err != nil {
		logger.Errorf("put tx error: %v", err)
		return nil, err
	}
	logger.Debug("put tx into world state")

	// tx total counter
	coinInfo.TxTotal += 1
	if err := store.PutCoinInfo(coinInfo); err != nil {
		logger.Errorf("Error put coin info: %v", err)
		return nil, err
	}

	logger.Debugf("coinbase execute result: %+v", execResult)
	return execResult.Bytes()
}

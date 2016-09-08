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
	pb "github.com/conseweb/common/protos"
	"github.com/golang/protobuf/proto"
)

func (coin *Lepuscoin) queryAddr(store Store, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgs
	}

	addr := args[0]
	queryResult := new(pb.QueryAddrResult)

	// account model
	account := MakeAccount(store)
	a, err := account.QueryAccountByAddr(addr)
	if err != nil {
		logger.Errorf("query account by addr error: %v", err)
		return nil, err
	}
	queryResult.Account = a

	logger.Debugf("query addr combind account: %v", a)

	if a.TxoutKey != "" {
		// utxo
		utxo := MakeUTXO(store)
		out, err := utxo.QueryTxOut(a.TxoutKey)
		if err != nil {
			logger.Errorf("query tx out error: %v", err)
			return nil, err
		}
		queryResult.Txout = out
	}

	logger.Debugf("query addr[%s] result: %+v", addr, queryResult)
	return queryResult.Bytes()
}

func (coin *Lepuscoin) queryAddrs(store Store, args []string) ([]byte, error) {
	results := &pb.QueryAddrResults{
		Results: make([]*pb.QueryAddrResult, 0),
	}

	// account model
	account := MakeAccount(store)
	// utxo
	utxo := MakeUTXO(store)
	for _, arg := range args {
		addr := arg
		queryResult := new(pb.QueryAddrResult)

		a, err := account.QueryAccountByAddr(addr)
		if err != nil {
			logger.Errorf("query account by addr error: %v", err)
			return nil, err
		}
		queryResult.Account = a
		logger.Debugf("query addr[%s] account: %v", addr, a)

		if a.TxoutKey != "" {
			out, err := utxo.QueryTxOut(a.TxoutKey)
			if err != nil {
				logger.Errorf("query tx out error: %v", err)
				return nil, err
			}
			queryResult.Txout = out
		}

		results.Results = append(results.Results, queryResult)
	}

	return proto.Marshal(results)
}

func (coin *Lepuscoin) queryTx(store Store, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgs
	}

	// utxo
	utxo := MakeUTXO(store)
	tx, err := utxo.QueryTx(args[0])
	if err != nil {
		logger.Errorf("utxo query tx return error: %v", err)
		return nil, err
	}

	logger.Debugf("query tx: %+v", tx)
	return tx.Bytes()
}

func (coin *Lepuscoin) queryCoin(store Store, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, ErrInvalidArgs
	}

	info := new(pb.LepuscoinInfo)
	info.CoinTotal = store.GetCoinbase()

	logger.Debugf("query lepuscoin info: %+v", info)
	return proto.Marshal(info)
}

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
package main

import (
	"bytes"
	"fmt"
	pb "github.com/conseweb/common/protos"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

var (
	// invoke function
	invoke_award_miner string = "awardMiner"
)

func (coin *Lepuscoin) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case invoke_award_miner:
		return coin.awardMiner(stub, args)
	default:
		return nil, fmt.Errorf("unsupported function type: %s", function)
	}
}

func (coin *Lepuscoin) awardMiner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("invoke awardMiner")
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid args: %v, required %v args", args, 2)
	}

	account := new(pb.Account)
	addr := args[0]
	award, err := parseCCToTC(args[1])
	if err != nil {
		return nil, err
	}

	account.Addr = addr
	if row, err := stub.GetRow(AccountModelTableName, []shim.Column{
		shim.Column{&shim.Column_String_{String_: addr}},
	}); err != nil || len(row.Columns) == 0 {
		logger.Warningf("get wallet %s info error: %v, maybe not exist, creating one...\n", addr, err)

		// maybe there is no miner's account, create one
		account.Balance = convTCToCC(award)
		account.AvailableBalance = 0.0
		account.FrozenBalance = convTCToCC(award)

		if _, err := stub.InsertRow(AccountModelTableName, shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: addr}},
				&shim.Column{Value: &shim.Column_Uint64{Uint64: uint64(award)}},
				&shim.Column{Value: &shim.Column_Uint64{Uint64: 0}},
				&shim.Column{Value: &shim.Column_Uint64{Uint64: uint64(award)}},
			},
		}); err != nil {
			logger.Errorf("table %s create account: %+v return error: %v", AccountModelTableName, account, err)
			return nil, err
		}
	} else {
		logger.Debugf("get wallet %s info: %+v. Len: %v\n", addr, row, len(row.Columns))
		logger.Debugf("row column 0: %v\n", row.Columns[0].GetString_())
		logger.Debugf("row column 1: %v\n", row.Columns[1].GetUint64())
		logger.Debugf("row column 2: %v\n", row.Columns[2].GetUint64())
		logger.Debugf("row column 3: %v\n", row.Columns[3].GetUint64())

		// miner's account has already build, update
		balance := row.Columns[1].GetUint64() + uint64(award)
		frozenBalance := row.Columns[3].GetUint64() + uint64(award)

		account.Balance = convTCToCC(coinunit(balance))
		account.FrozenBalance = convTCToCC(coinunit(frozenBalance))

		if _, err := stub.ReplaceRow(AccountModelTableName, shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: addr}},
				&shim.Column{Value: &shim.Column_Uint64{Uint64: balance}},
				&shim.Column{Value: &shim.Column_Uint64{Uint64: row.Columns[2].GetUint64()}},
				&shim.Column{Value: &shim.Column_Uint64{Uint64: frozenBalance}},
			},
		}); err != nil {
			logger.Errorf("update account went wrong: %v\n", err)
			return nil, err
		}
	}

	// 2. set count of coin
	coinBytes, err := stub.GetState(CountCoins)
	if err != nil {
		logger.Errorf("get state %s return error: %v\n", CountCoins, err)
		return nil, err
	}
	if len(coinBytes) == 0 {
		coinBytes = []byte("0")
	}
	historyCoins, err := strconv.ParseUint(bytes.NewBuffer(coinBytes).String(), 10, 64)
	if err != nil {
		logger.Errorf("parse string to uint64 return error: %v\n", err)
		return nil, err
	}

	countcoins := uint64(award) + historyCoins
	logger.Debugf("Lepuscoin current count: %d, history count: %d\n", countcoins, historyCoins)
	if err := stub.PutState(CountCoins, bytes.NewBufferString(strconv.FormatUint(countcoins, 10)).Bytes()); err != nil {
		logger.Errorf("update state %s return error: %v", CountCoins, err)
		return nil, err
	}

	logger.Debugf("miner account: %+v\n", account)
	return proto.Marshal(account)
}

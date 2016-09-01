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
	"fmt"
	pb "github.com/conseweb/common/protos"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
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
	}); err != nil {
		// maybe there is no miner's account, create one
		account.Balance = convTCToCC(award)
		account.AvailableBalance = 0.0
		account.FrozenBalance = convTCToCC(award)

		if _, err := stub.InsertRow(AccountModelTableName, shim.Row{
			Columns: []*shim.Column{
				&shim.Column_String_{String_: addr},
				&shim.Column_Uint64{Uint64: uint64(award)},
				&shim.Column_Uint64{Uint64: 0},
				&shim.Column_Uint64{Uint64: uint64(award)},
			},
		}); err != nil {
			return nil, err
		}
	} else {
		// miner's account has already build, update
		balance := row.Columns[1].GetUint64() + uint64(award)
		frozenBalance := row.Columns[3].GetUint64() + uint64(award)

		account.Balance = convTCToCC(coinunit(balance))
		account.FrozenBalance = convTCToCC(coinunit(frozenBalance))

		row.Columns[1] = &shim.Column_Uint64{Uint64: balance}
		row.Columns[3] = &shim.Column_Uint64{Uint64: frozenBalance}
		if _, err := stub.ReplaceRow(AccountModelTableName, row); err != nil {
			return nil, err
		}
	}

	return proto.Marshal(account)
}

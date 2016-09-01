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
	// query functions
	query_account = "queryAccount"
)

func (coin *Lepuscoin) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case query_account:
		return coin.queryAccount(stub, args)
	default:
		return nil, fmt.Errorf("unsupported function type: %s", function)
	}
}

func (coin *Lepuscoin) queryAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid params: %v, required %v args", args, 1)
	}

	account := new(pb.Account)
	addr := args[0]

	// fetch account info from table
	row, err := stub.GetRow(AccountModelTableName, []shim.Column{
		shim.Column{&shim.Column_String_{String_: addr}},
	})
	if err != nil {
		return nil, err
	}

	account.Addr = addr
	account.Balance = convTCToCC(coinunit(row.Columns[1].GetUint64()))
	account.AvailableBalance = convTCToCC(coinunit(row.Columns[2].GetUint64()))
	account.FrozenBalance = convTCToCC(coinunit(row.Columns[3].GetUint64()))

	return proto.Marshal(account)
}

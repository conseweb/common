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
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
	"strconv"
)

var (
	logger = logging.MustGetLogger(CoinName)

	// account table column defines
	AccountModelTableColumns = []*shim.ColumnDefinition{
		// user id column, version 1 uuid
		//&shim.ColumnDefinition{"userID", shim.ColumnDefinition_STRING, true},
		// device id column
		//&shim.ColumnDefinition{"deviceID", shim.ColumnDefinition_STRING, true},
		// hd wallet address column
		&shim.ColumnDefinition{"walletAddr", shim.ColumnDefinition_STRING, true},
		// balance column, available balance + frozen balance
		&shim.ColumnDefinition{"balance", shim.ColumnDefinition_UINT64, false},
		// available balance column
		&shim.ColumnDefinition{"availableBalance", shim.ColumnDefinition_UINT64, false},
		// frozen balance column
		&shim.ColumnDefinition{"frozenBalance", shim.ColumnDefinition_UINT64, false},
	}
)

const (
	// coin name
	CoinName string = "Lepuscoin"
	// account model table name
	AccountModelTableName string = "LepuscoinAccount"
	// world state of how many coins has been mined
	CountCoins string = "LepusCoinCount"
)

type Lepuscoin struct {
}

func (coin *Lepuscoin) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "deploy" {
		return nil, errors.New("invalid function name, 'deploy' only")
	}

	// create account table
	if _, err := stub.GetTable(AccountModelTableName); err == shim.ErrTableNotFound {
		err = stub.CreateTable(AccountModelTableName, AccountModelTableColumns)
		if err != nil {
			logger.Errorf("create %s table return error: %v\n", AccountModelTableName, err)
			return nil, err
		}

		logger.Debugf("successfully create table %s\n", AccountModelTableName)
	}

	// create world state of how many coins has been mined
	if _, err := stub.GetState(CountCoins); err != nil {
		err = stub.PutState(CountCoins, bytes.NewBufferString(strconv.FormatUint(0, 10)).Bytes())
		if err != nil {
			logger.Errorf("create world state of how many coins has been mined return error: %v\n", err)
			return nil, err
		}

		logger.Debugf("successfully set state %s to 0", CountCoins)
	}

	logger.Debugf("bytes.NewBufferString(strconv.FormatUint(0, 10)): %s\n", bytes.NewBufferString(strconv.FormatUint(0, 10)).String())
	logger.Debugf("deploy %s successfully\n", CoinName)
	return nil, nil
}

func main() {
	logging.SetLevel(logging.DEBUG, CoinName)
	err := shim.Start(new(Lepuscoin))
	if err != nil {
		logger.Fatalf("deploy %s return err: %v\n", CoinName, err)
	}
}

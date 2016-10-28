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
	"errors"
	"github.com/hyperledger/fabric/flogging"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var (
	logger = logging.MustGetLogger("lepuscoin")
)

type Lepuscoin struct {
}

// Init deploy Lepuscoin into vp
func (coin *Lepuscoin) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	flogging.LoggingInit("lepuscoin")
	if function != "deploy" {
		return nil, errors.New("invalid function name, 'deploy' only")
	}

	// construct a new store
	store := MakeChaincodeStore(stub)

	// deploy lepuscoin chaincode only need to set coin counter
	if err := store.InitCoinInfo(); err != nil {
		return nil, err
	}

	logger.Info("deploy Lepuscoin successfully")
	return nil, nil
}

// Invoke function
const (
	IF_COINBASE string = "invoke_coinbase"
	IF_TRANSFER string = "invoke_transfer"
)

// Invoke
func (coin *Lepuscoin) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// construct a new store
	store := MakeChaincodeStore(stub)

	switch function {
	case IF_COINBASE:
		return coin.coinbase(store, args)
	case IF_TRANSFER:
		return coin.transfer(store, args)
	default:
		return nil, ErrUnsupportedOperation
	}
}

// Query function
const (
	QF_ADDRS = "query_addrs"
	QF_TX    = "query_tx"
	QF_COIN  = "query_coin"
)

// Query
func (coin *Lepuscoin) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// construct a new store
	store := MakeChaincodeStore(stub)

	switch function {
	case QF_ADDRS:
		return coin.queryAddrs(store, args)
	case QF_TX:
		return coin.queryTx(store, args)
	case QF_COIN:
		return coin.queryCoin(store, args)
	default:
		return nil, ErrUnsupportedOperation
	}
}

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
	"errors"

	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
	"strconv"
)

var (
	logger = logging.MustGetLogger("lepuscoin")
)

func init() {
	logging.SetLevel(logging.DEBUG, "lepuscoin")
}

type Lepuscoin struct {
}

// Init deploy Lepuscoin into vp
func (coin *Lepuscoin) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "deploy" {
		return nil, errors.New("invalid function name, 'deploy' only")
	}

	// construct a new store
	store := MakeChaincodeStore(stub)

	// deploy lepuscoin chaincode only need to set coin counter
	if err := store.AddCoinbase(0); err != nil {
		return nil, err
	}

	logger.Debug("deploy Lepuscoin successfully")
	return nil, nil
}

// Invoke function
const (
	IF_COINBASE string = "invoke_coinbase"
)

// Invoke
func (coin *Lepuscoin) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// construct a new store
	store := MakeChaincodeStore(stub)

	switch function {
	case IF_COINBASE:
		return coin.coinbase(store, args)
	default:
		return nil, fmt.Errorf("unsupported function type: %s", function)
	}
}

// Query function
const (
	QF_TX = "query_tx"
	QF_CB = "query_cb"
)

// Query
func (coin *Lepuscoin) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// construct a new store
	store := MakeChaincodeStore(stub)

	switch function {
	case QF_TX:
		return coin.queryTx(store, args)
	case QF_CB:
		return coin.queryCB(store, args)
	default:
		return nil, fmt.Errorf("unsupported function type: %s", function)
	}
}

func (coin *Lepuscoin) queryTx(store Store, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("invalid args")
	}

	// utxo
	utxo := MakeUTXO(store)
	tx, err := utxo.Query(args[0])
	if err != nil {
		return nil, fmt.Errorf("Error querying for transaction:  %s", err)
	}

	logger.Debugf("query tx return bytes: %s", bytes.NewBuffer(tx).String())
	return tx, nil
}

func (coin *Lepuscoin) queryCB(store Store, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("invalid args")
	}

	cb := store.GetCoinbase()
	return bytes.NewBufferString(strconv.FormatUint(cb, 10)).Bytes(), nil
}

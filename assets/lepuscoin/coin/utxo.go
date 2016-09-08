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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"

	pb "github.com/conseweb/common/protos"
)

// UTXO includes the storage for the chaincode API or an in memory
// store for testing
type UTXO struct {
	Store Store
}

// MakeUTXO constructs a new UTXO with the given store
func MakeUTXO(store Store) *UTXO {
	utxo := new(UTXO)
	utxo.Store = store

	return utxo
}

// GetTransactionHash returns the tx hash
func (u *UTXO) getTXHash(txData []byte) [32]byte {
	fHash := sha256.Sum256(txData)
	return sha256.Sum256(fHash[:])
}

// IsCoinbase returns true if this is a coinbase transaction, false otherwise
func (u *UTXO) isCoinbase(index uint32) bool {
	return index == math.MaxUint32
}

// Execute processes the given transaction and outputs a result
func (u *UTXO) Execute(tx *pb.TX) (*pb.ExecResult, error) {
	newTx := tx.Copy()
	txhash := newTx.TxHash()
	execResult := &pb.ExecResult{}

	// Loop through outputs first
	for index, output := range newTx.Txout {
		currKey := &Key{TxHashAsHex: hex.EncodeToString(txhash), TxIndex: uint32(index)}
		_, ok, err := u.Store.GetTxOut(currKey)
		if err != nil {
			return nil, fmt.Errorf("Error getting state from stores: %s", err)
		}

		if ok == true {
			return nil, fmt.Errorf("COLLISTION detected for key = %v, with output script length: %d", currKey, len(output.Script))
		}

		// Store the output in utxo
		if err := u.Store.PutTxOut(currKey, output); err != nil {
			return nil, err
		}

		logger.Debugf("put tx output %s:%v", currKey.String(), output)
		execResult.SumCurrentOutputs += output.Value
	}

	// Now loop over inputs
	for _, input := range newTx.Txin {
		prevTxHash := input.SourceHash
		prevOutputIx := input.Ix

		if u.isCoinbase(prevOutputIx) {
			execResult.IsCoinbase = true
			logger.Debugf("input[%+v] is coinbase!", input)
		} else {
			keyToPrevOutput := &Key{TxHashAsHex: hex.EncodeToString(prevTxHash), TxIndex: prevOutputIx}
			value, err := u.QueryTxOut(keyToPrevOutput.String())
			if err != nil {
				return nil, fmt.Errorf("Could not find previous transaction output with key = %v, err: %v", keyToPrevOutput, err)
			}

			// verify script

			// verified, now remove prior outputs
			if err := u.Store.DelTxOut(keyToPrevOutput); err != nil {
				return nil, err
			}
			logger.Debugf("delete previous tx out: %s", keyToPrevOutput)
			execResult.SumPriorOutputs += value.Value
		}

		if err := u.Store.PutTx(newTx); err != nil {
			return nil, err
		}
		logger.Debug("put tx into db")
	}

	return execResult, nil
}

// Query search the storage for a given transaction hash
func (u *UTXO) QueryTx(txHashHex string) (*pb.TX, error) {
	tx, _, err := u.Store.GetTx(txHashHex)
	return tx, err
}

// QueryTxOut search txout
func (u *UTXO) QueryTxOut(keyStr string) (*pb.TX_TXOUT, error) {
	key, err := parseKey(keyStr)
	if err != nil {
		return nil, err
	}

	out, _, err := u.Store.GetTxOut(key)
	if err != nil {
		return nil, err
	}

	return out, nil
}

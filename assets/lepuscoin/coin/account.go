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
	"encoding/hex"
	"errors"
	"fmt"
	pb "github.com/conseweb/common/protos"
)

type Account struct {
	store Store
}

// MakeAccount constructs a new account model
func MakeAccount(store Store) *Account {
	account := new(Account)
	account.store = store

	return account
}

// Coinbase give addr some coin and add coin counter
func (a *Account) Coinbase(tx *pb.TX) error {
	txHash := tx.TxHash()
	for idx, txout := range tx.Txout {
		outKey := &Key{TxHashAsHex: hex.EncodeToString(txHash), TxIndex: uint32(idx)}

		if txout.Addr == "" {
			return errors.New("no addr specified")
		}

		// add coin counter
		if err := a.store.AddCoinbase(txout.Value); err != nil {
			return err
		}

		// add account
		account, err := a.store.GetAccount(txout.Addr)
		if err != nil {
			logger.Warningf("Can't get account, creating one...")
			account = new(pb.Account)
		}

		account.Balance += txout.Value
		account.Addr = txout.Addr
		account.TxoutKey = outKey.String()

		if err := a.store.PutAccount(account); err != nil {
			return err
		}
	}

	return nil
}

// Transfer transfer balance from in to out
func (a *Account) Transfer(tx *pb.TX) error {
	for _, ti := range tx.Txin {
		prevTxHash := ti.SourceHash
		prevOutputIx := ti.Ix
		keyToPrevOutput := &Key{TxHashAsHex: hex.EncodeToString(prevTxHash), TxIndex: prevOutputIx}
		value, ok, err := a.store.GetTxOut(keyToPrevOutput)
		if err != nil {
			return fmt.Errorf("Error getting state form store: %v", err)
		}

		if !ok {
			// Previous output not found
			return fmt.Errorf("Could not find previous transaction output with key = %v", keyToPrevOutput)
		}

		account, err := a.store.GetAccount(value.Addr)
		if err != nil {
			return err
		}

		if account.Balance < value.Value {
			return fmt.Errorf("Account %s dont't have enough balance", value.Addr)
		}

		account.Balance -= value.Value
		account.TxoutKey = ""
		if err := a.store.PutAccount(account); err != nil {
			return err
		}
	}

	txHash := tx.TxHash()
	for idx, to := range tx.Txout {
		account, err := a.store.GetAccount(to.Addr)
		if err != nil {
			logger.Warningf("get account doesnt exist, creating one...")
			account = new(pb.Account)
			account.Addr = to.Addr
		}

		outKey := &Key{TxHashAsHex: hex.EncodeToString(txHash), TxIndex: uint32(idx)}
		account.Balance = to.Value
		account.TxoutKey = outKey.String()

		if err := a.store.PutAccount(account); err != nil {
			return err
		}
	}

	return nil
}

// QueryAccountByAddr query account info using addr
func (a *Account) QueryAccountByAddr(addr string) (*pb.Account, error) {
	account, err := a.store.GetAccount(addr)
	if err != nil {
		account = &pb.Account{
			Addr:     addr,
			Balance:  0,
			TxoutKey: "",
		}

		if err := a.store.PutAccount(account); err != nil {
			return nil, err
		}
	}

	return account, nil
}

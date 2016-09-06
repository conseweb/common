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
func (a *Account) Coinbase(supply uint64, addr string) error {
	if addr == "" {
		return errors.New("no addr specified")
	}

	// add coin counter
	if err := a.store.AddCoinbase(supply); err != nil {
		return err
	}

	// add account
	account, err := a.store.GetAccount(addr)
	if err != nil {
		return err
	}

	account.Balance += supply
	return a.store.PutAccount(account)
}

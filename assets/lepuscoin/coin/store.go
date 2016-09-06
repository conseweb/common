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
	"fmt"

	"bytes"
	pb "github.com/conseweb/common/protos"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

const coinbase = "Lepuscoinbase"

// Key represents the key for a transaction in storage. It has both a
// hash and index
type Key struct {
	TxHashAsHex string
	TxIndex     uint32
}

func (k *Key) String() string {
	return fmt.Sprintf("%s:%d", k.TxHashAsHex, k.TxIndex)
}

func generateAccountKey(addr string) string {
	return fmt.Sprintf("account_addr_%s", addr)
}

// Store interface describes the storage used by this chaincode. The interface
// was created so either the state database store can be used or a in memory
// store can be used for unit testing.
type Store interface {
	GetTxOut(*Key) (*pb.TX_TXOUT, bool, error)
	PutTxOut(*Key, *pb.TX_TXOUT) error
	DelTxOut(*Key) error
	GetTx(string) ([]byte, bool, error)
	PutTx(string, []byte) error
	GetCoinbase() uint64
	AddCoinbase(uint64) error
	GetAccount(string) (*pb.Account, error)
	PutAccount(*pb.Account) error
}

// Store struct uses a chaincode stub for state access
type ChaincodeStore struct {
	stub shim.ChaincodeStubInterface
}

// MakeChaincodeStore returns a store for storing keys in the state
func MakeChaincodeStore(stub shim.ChaincodeStubInterface) Store {
	store := &ChaincodeStore{}
	store.stub = stub
	return store
}

// GetTxOut returns the transaction for a given key
func (s *ChaincodeStore) GetTxOut(key *Key) (*pb.TX_TXOUT, bool, error) {
	data, err := s.stub.GetState(key.String())
	if err != nil {
		return nil, false, fmt.Errorf("Error getting state from stub:  %s", err)
	}
	if data == nil || len(data) == 0 {
		return nil, false, nil
	}

	// Value found, unmarshal
	value := &pb.TX_TXOUT{}
	if err := proto.Unmarshal(data, value); err != nil {
		return nil, false, fmt.Errorf("Error unmarshalling value:  %s", err)
	}

	return value, true, nil
}

// DelTxOut deletes the transaction for the given key
func (s *ChaincodeStore) DelTxOut(key *Key) error {
	return s.stub.DelState(key.String())
}

// PutTxOut stores the given transaction and key
func (s *ChaincodeStore) PutTxOut(key *Key, value *pb.TX_TXOUT) error {
	data, err := proto.Marshal(value)
	if err != nil {
		return fmt.Errorf("Error marshalling value to bytes:  %s", err)
	}

	return s.stub.PutState(key.String(), data)
}

// GetTx returns a transaction for the given hash
func (s *ChaincodeStore) GetTx(key string) ([]byte, bool, error) {
	data, err := s.stub.GetState(key)
	if err != nil {
		return nil, false, fmt.Errorf("Error getting state from stub:  %s", err)
	}
	if data == nil || len(data) == 0 {
		return nil, false, nil
	}

	return data, true, nil
}

// PutTx adds a transaction to the state with the hash as a key
func (s *ChaincodeStore) PutTx(key string, value []byte) error {
	return s.stub.PutState(key, value)
}

// GetCoinbase returns monetary supply, based on account model
func (s *ChaincodeStore) GetCoinbase() uint64 {
	data, err := s.stub.GetState(coinbase)
	if err != nil || data == nil || len(data) == 0 {
		logger.Errorf("Error getting state[%v] %v", coinbase, err)
		return 0
	}

	supply, err := strconv.ParseUint(bytes.NewBuffer(data).String(), 10, 64)
	if err != nil {
		logger.Errorf("strconv.ParseUint return error: %v", err)
		return 0
	}

	return supply
}

// AddCoinbase add coinbase into monetary supply counter
func (s *ChaincodeStore) AddCoinbase(add uint64) error {
	supply := s.GetCoinbase() + add

	if err := s.stub.PutState(coinbase, bytes.NewBufferString(strconv.FormatUint(supply, 10)).Bytes()); err != nil {
		logger.Errorf("Error setting state[%s] %v", coinbase, err)
		return err
	}

	return nil
}

// GetAccount returns account from world states
func (s *ChaincodeStore) GetAccount(addr string) (*pb.Account, error) {
	key := generateAccountKey(addr)
	data, err := s.stub.GetState(key)
	if err != nil {
		return nil, err
	}

	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("no account found")
	}

	account := new(pb.Account)
	if err := proto.Unmarshal(data, account); err != nil {
		return nil, err
	}

	return account, nil
}

// PutAccount update or insert account into world states
func (s *ChaincodeStore) PutAccount(account *pb.Account) error {
	key := generateAccountKey(account.Addr)

	aBytes, err := proto.Marshal(account)
	if err != nil {
		return err
	}

	return s.stub.PutState(key, aBytes)
}

// InMemoryStore used for unit testing
type InMemoryStore struct {
	Map        map[*Key]*pb.TX_TXOUT
	TranMap    map[string][]byte
	Accounts   map[string]*pb.Account
	coinsupply uint64
}

// MakeInMemoryStore creates a new in memory store
func MakeInMemoryStore() Store {
	ims := &InMemoryStore{}
	ims.Map = make(map[*Key]*pb.TX_TXOUT)
	ims.TranMap = make(map[string][]byte)
	return ims
}

// GetTxOut returns the transaction for the given key
func (ims *InMemoryStore) GetTxOut(key *Key) (*pb.TX_TXOUT, bool, error) {
	value, ok := ims.Map[key]
	return value, ok, nil
}

// DelTxOut deletes the given key and corresponding transactions
func (ims *InMemoryStore) DelTxOut(key *Key) error {
	delete(ims.Map, key)
	return nil
}

// PutTxOut saves the key and transaction in memory
func (ims *InMemoryStore) PutTxOut(key *Key, value *pb.TX_TXOUT) error {
	ims.Map[key] = value
	return nil
}

// GetTx returns the transaction for the given hash
func (ims *InMemoryStore) GetTx(key string) ([]byte, bool, error) {
	value, ok := ims.TranMap[key]
	return value, ok, nil
}

// PutTx saves the hash and transaction in memory
func (ims *InMemoryStore) PutTx(key string, value []byte) error {
	ims.TranMap[key] = value
	return nil
}

// GetCoinbase returns monetary supply, based on account model
func (s *InMemoryStore) GetCoinbase() uint64 {
	return s.coinsupply
}

// AddCoinbase add coinbase into monetary supply counter
func (s *InMemoryStore) AddCoinbase(add uint64) error {
	supply := s.GetCoinbase() + add
	s.coinsupply = supply

	return nil
}

// GetAccount returns account from world states
func (s *InMemoryStore) GetAccount(addr string) (*pb.Account, error) {
	key := generateAccountKey(addr)
	account, ok := s.Accounts[key]
	if !ok {
		return nil, fmt.Errorf("account not found")
	}

	return account, nil
}

// PutAccount update or insert account into world states
func (s *InMemoryStore) PutAccount(account *pb.Account) error {
	key := generateAccountKey(account.Addr)

	s.Accounts[key] = account

	return nil
}

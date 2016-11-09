package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestRegister(t *testing.T) {
	stub := shim.NewMockStub("ex02", new(PoeService))
	d, e := stub.MockInvoke("1", "register", []string{"xiebo", "wangsiyu", "shanglijun"})
	if e != nil {
		fmt.Println("Invoke,failed", e)
		t.FailNow()
	}
	fmt.Println(string(d))
}
func TestExistence(t *testing.T) {
	stub := shim.NewMockStub("ex02", new(PoeService))
	d, e := stub.MockQuery("existence", []string{"xiebo", "wangsiyu", "shanglijun"})
	if e != nil {
		fmt.Println("Invoke,failed", e)
		t.FailNow()
	}
	fmt.Println(string(d))
}

func TestState(t *testing.T) {
	stub := shim.NewMockStub("ex02", new(PoeService))
	hashKey()
	stub.State["xiebo"]
}

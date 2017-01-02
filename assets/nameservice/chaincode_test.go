package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type kv struct {
	k string
	v string
}

var param []kv = []kv{
	kv{k: "", v: "220.181.112.200"},
	kv{k: "www.baidu.com", v: ""},
	kv{k: "www.baid u.com", v: "220.181.112.200"},
	kv{k: "www.baidu.com", v: "220.18 1.112.200"},
}

var data []kv = []kv{
	kv{k: "www.baidu.com", v: "220.181.112.200"},
	kv{k: "www.wx.com", v: "220.181.112.201"},
	kv{k: "www.qq.com", v: "220.181.112.202"},
	kv{k: "www.taobao.com", v: "220.181.112.203"},
	kv{k: "www.hao123.com", v: "220.181.112.204"},
	kv{k: "www.360.com", v: "220.181.112.205"},
}

func checkInit01(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockInit("1", "init", nil)
	if err != nil {
		t.Logf("init failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func checkInit02(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockInit("2", "deploy", nil)
	if err != nil {
		t.Logf("init failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func Test_Init(t *testing.T) {
	cc := new(NameService)
	stub := shim.NewMockStub("namesrvc01", cc)
	// 空值检测
	checkInit01(stub, t)
	// 成功返回200
	checkInit02(stub, t)
}

func checkAddotoParam(stub *shim.MockStub, t *testing.T) {
	for i, obj := range param {
		_, err := stub.MockInvoke(string(i), "addoto", []string{obj.k, obj.v})
		if err != nil {
			t.Logf("init failed%v: %v", i, err)
		}
	}
}

func checkAddotoInvoke(stub *shim.MockStub, t *testing.T) {
	for i, obj := range data {
		_, err := stub.MockInvoke(string(i), "addoto", []string{obj.k, obj.v})
		if err != nil {
			t.Logf("init failed%v: %v", i, err)
		}
	}
}

func checkAddotoQuery(stub *shim.MockStub, t *testing.T) {
	for i, obj := range data {
		data, err := stub.MockQuery("query", []string{obj.k})
		if err != nil {
			t.Logf("init failed%v: %v", i, err)
		}
		t.Logf("data1 : %v", string(data))
		data, err = stub.MockQuery("query", []string{obj.v})
		if err != nil {
			t.Logf("init failed%v: %v", i, err)
		}
		t.Logf("data2 : %v", string(data))
	}
}

func checkAddotoDel(stub *shim.MockStub, t *testing.T) {
	for i, obj := range data {
		_, err := stub.MockInvoke(string(i), "deloto", []string{obj.k})
		if err != nil {
			t.Logf("init failed%v: %v", i, err)
		}
	}
}

func TestAddoto(t *testing.T) {
	cc := new(NameService)
	stub := shim.NewMockStub("namesrvc02", cc)
	checkAddotoParam(stub, t)
	t.Log("===checkAddotoParam over!")
	checkAddotoInvoke(stub, t)
	t.Log("===checkAddotoInvoke over!")
	checkAddotoQuery(stub, t)
	t.Log("===checkAddotoQuery over!")
	checkAddotoDel(stub, t)
	t.Log("===checkAddotoDel over!")
	checkAddotoQuery(stub, t)
	t.Log("===checkAddotoQuery over!")
}

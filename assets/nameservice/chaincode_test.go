package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit01(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockInit("1", "init", nil)
	if err != nil {
		t.Errorf("init failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func checkInit02(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockInit("2", "deploy", nil)
	if err != nil {
		t.Errorf("init failed %v", err)
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

func checkInvoke01(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockInvoke("1", "add", nil)
	if err != nil {
		t.Errorf("query failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func checkInvoke02(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockInvoke("2", "add", []string{"wwww.baidu.co m:111.206.223.206"})
	if err != nil {
		t.Errorf("query failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func checkInvoke03(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockInvoke("3", "add", []string{"wwww.baidu.com::111.206.223.206"})
	if err != nil {
		t.Errorf("query failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func checkInvoke04(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockInvoke("4", "add", []string{"wwww.baidu.com:111.206.223.206", "wwww.baidu.com:111.206.223.206"})
	if err != nil {
		t.Errorf("query failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func checkInvoke05(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockInvoke("5", "add", []string{"wwww.baidu.com:111.206.223.206"})
	if err != nil {
		t.Errorf("query failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func Test_Invoke(t *testing.T) {
	cc := new(NameService)
	stub := shim.NewMockStub("namesrvc02", cc)
	//空值检测
	checkInvoke01(stub, t)
	//不能包含空字符
	checkInvoke02(stub, t)
	//只能包含一个分隔符':'
	checkInvoke03(stub, t)
	//每次只处理一条数据
	checkInvoke04(stub, t)
	//正确返回输入参数
	checkInvoke05(stub, t)
}

func checkQuery01(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockQuery("query", nil)
	if err != nil {
		t.Errorf("query failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func checkQuery02(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockQuery("query", []string{" "})
	if err != nil {
		t.Errorf("query failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func checkQuery03(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockQuery("query", []string{"wwww.baidu.com", "wwww.baidu.com"})
	if err != nil {
		t.Errorf("query failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func checkQuery04(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockQuery("query", []string{"wwww.baidu.com"})
	if err != nil {
		t.Errorf("query failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func checkQuery05(stub *shim.MockStub, t *testing.T) {
	data, err := stub.MockQuery("query", []string{"111.206.223.206"})
	if err != nil {
		t.Errorf("query failed %v", err)
		return
	}
	t.Logf("data: %s", data)
}

func Test_Query(t *testing.T) {
	cc := new(NameService)
	stub := shim.NewMockStub("namesrvc03", cc)
	//存入数据，存两份，第二份数据 kv 交换存储
	checkInvoke05(stub, t)
	//空值检测
	checkQuery01(stub, t)
	//参数不能为空
	checkQuery02(stub, t)
	//每次只处理一条数据
	checkQuery03(stub, t)
	//get v by k
	checkQuery04(stub, t)
	//get k by v
	checkQuery05(stub, t)
}

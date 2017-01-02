## Name Service（命名服务）
kv 存储， 存储两份，第二份kv 交换存储。满足正反查询的效果

### developer test

run fabric peer using docker

docker-commpose.yml
```
membersrvc:
  image: hyperledger/fabric-membersrvc
  ports:
    - "7054:7054"
  command: membersrvc
vp0:
  image: hyperledger/fabric-peer
  ports:
    - "7050:7050"
    - "7051:7051"
    - "7053:7053"
  environment:
    - CORE_PEER_ADDRESSAUTODETECT=true
    - CORE_VM_ENDPOINT=unix:///var/run/docker.sock
    - CORE_LOGGING_LEVEL=DEBUG
    - CORE_PEER_ID=vp0
    - CORE_PEER_PKI_ECA_PADDR=membersrvc:7054
    - CORE_PEER_PKI_TCA_PADDR=membersrvc:7054
    - CORE_PEER_PKI_TLSCA_PADDR=membersrvc:7054
    - CORE_SECURITY_ENABLED=false
	- CORE_SECURITY_ENROLLID=test_vp0
    - CORE_SECURITY_ENROLLSECRET=MwYpmSRjupbT
  links:
    - membersrvc
  volumes_from:
    - share_data
  command: sh -c "sleep 5; peer node start --peer-chaincodedev"
```
start docker
```
docker-compose up
```
build nameservice and run locally
```
CORE_CHAINCODE_ID_NAME=namesrvc_cc_demo_001 CORE_PEER_ADDRESS=0.0.0.0:7051 ./nameservice
```
```
docker ps -a
docker exec -it xxxx bash
```
### deploy examples

deploy request
```
{
    "jsonrpc": "2.0",
    "method": "deploy",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_001"
        },
        "ctorMsg": {
            "function": "deploy",
            "args": []
        },
        "secureContext": "jim"
    },
    "id": 1
}
```

### transfer examples

#### one to one
- 添加一条数据(一对一)
- 参数规范： 不能包含空字符;每次只处理一条数据;args[0] 为key, args[1] 为value
- 存储会存储两份，第二次存储k和v交换存储
- kv 都不存在或k存在v 不存在时存储，存储后，原v 值将被覆盖， kv交换 逻辑一样

add request
```
{
    "jsonrpc": "2.0",
    "method": "invoke",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_001"
        },
        "ctorMsg": {
            "function": "addoto",
            "args":["wwww.baidu.com:111.206.223.206"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
- 删除一条数据
- 先根据k 删除，再用查出得v 作为新k删除
del request

```
{
    "jsonrpc": "2.0",
    "method": "invoke",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_001"
        },
        "ctorMsg": {
            "function": "deloto",
            "args":["wwww.baidu.com"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
query request
```
{
    "jsonrpc": "2.0",
    "method": "query",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_001"
        },
        "ctorMsg": {
            "function": "query",
            "args":["wwww.baidu.com"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
query response
```
{
    "jsonrpc": "2.0",
    "result": {
        "status": "OK",
        "message":"111.206.223.206"},
    "id": 1
}
```

#### one to many
- 添加一条数据(一对多)
- 参数规范： 不能包含空字符;每次只处理一条数据;args[0] 为key, args[1] 为value
- 存储会存储两份，第二次存储k和v交换存储
- kv 都不存在或k存在v 不存在时存储，存储后，原来的v = 原有v,现有v， kv交换逻辑： 当kv都不存在或者k 存在v 不存在时存储，存储后，原v 值将被覆盖

invoke request 1
```
{
    "jsonrpc": "2.0",
    "method": "invoke",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_002"
        },
        "ctorMsg": {
            "function": "addotm",
            "args":["wwww.baidu.com:111.206.223.206"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
invoke request 2
```
{
    "jsonrpc": "2.0",
    "method": "invoke",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_002"
        },
        "ctorMsg": {
            "function": "addotm",
            "args":["wwww.baidu.com:111.206.223.207"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
query request
```
{
    "jsonrpc": "2.0",
    "method": "query",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_002"
        },
        "ctorMsg": {
            "function": "query",
            "args":["wwww.baidu.com"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
query response
```
{
    "jsonrpc": "2.0",
    "result": {
        "status": "OK",
        "message":"111.206.223.206,111.206.223.207"},
    "id": 1
}
```

#### many to many

- 添加一条数据(多对多)
- 参数规范： 不能包含空字符;每次只处理一条数据;args[0] 为key, args[1] 为value
- 存储会存储两份，第二次存储k和v交换存储
- kv 都不存在或k存在v 不存在时存储，存储后，原来的v = 原有v,现有v， kv交换逻辑一样

invoke request 1
```
{
    "jsonrpc": "2.0",
    "method": "invoke",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_003"
        },
        "ctorMsg": {
            "function": "addmtm",
            "args":["wwww.baidu.com:111.206.223.206"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
invoke request 2
```
{
    "jsonrpc": "2.0",
    "method": "invoke",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_003"
        },
        "ctorMsg": {
            "function": "addmtm",
            "args":["wwww.baidu.com:111.206.223.207"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
invoke request 3
```
{
    "jsonrpc": "2.0",
    "method": "invoke",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_003"
        },
        "ctorMsg": {
            "function": "addmtm",
            "args":["wwww.qq.com:111.206.223.207"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
query request 1
```
{
    "jsonrpc": "2.0",
    "method": "query",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_003"
        },
        "ctorMsg": {
            "function": "query",
            "args":["wwww.baidu.com"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
query response
```
{
    "jsonrpc": "2.0",
    "result": {
        "status": "OK",
        "message":"111.206.223.206,111.206.223.207"},
    "id": 1
}
```
query request 2
```
{
    "jsonrpc": "2.0",
    "method": "query",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "namesrvc_cc_demo_003"
        },
        "ctorMsg": {
            "function": "query",
            "args":["111.206.223.207"]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```
query response
```
{
    "jsonrpc": "2.0",
    "result": {
        "status": "OK",
        "message":"wwww.baidu.com,wwww.qq.com"},
    "id": 1
}
```
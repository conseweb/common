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
build poe and run locally
```
CORE_CHAINCODE_ID_NAME=poe_cc_demo_001 CORE_PEER_ADDRESS=0.0.0.0:7051 ./poe
```
```
docker ps -a
docker exec -it xxxx bash
```
### deploy request

REST
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
            "args": [
            ]
        },
        "secureContext": ""
    },
    "id": 1
}
```

### invoke request

#### invoke_register
REST
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
            "function": "add",
            "args":["wwww.baidu.com:111.206.223.206"]
        },
        "secureContext": ""
    },
    "id": 1
}
```

### query request
#### query_existence

REST
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
        "secureContext": ""
    },
    "id": 1
}
```
Response
```
{
    "jsonrpc": "2.0",
    "result": {
        "status": "OK",
        "message":"111.206.223.206"},
    "id": 1
}
```

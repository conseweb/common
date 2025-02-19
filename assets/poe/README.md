## Proof of Existence Service（存在性证明服务）

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
            "name": "poe_cc_demo_001"
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
            "name": "poe_cc_demo_001"
        },
        "ctorMsg": {
            "function": "register",
            "args":["base","nicsflvowkzdegrh","zrptjuswagkxyivq","bqohtvxernjmwdka","jcrxevpmbyauingl","uhtcvkfpwgaeozbd"]
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
            "name": "poe_cc_demo_001"
        },
        "ctorMsg": {
            "function": "existence",
            "args":["base","nicsflvowkzdegrh","zrptjuswagkxyivq","bqohtvxernjmwdka","jcrxevpmbyauingl","uhtcvkfpwgaeozbd"]
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
        "message":"[{\"key\":\"nicsflvowkzdegrh\",\"hash_key\":\"a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434ad41d8cd98f00b204e9800998ecf8427e10\",\"exist\":true},{\"key\":\"zrptjuswagkxyivq\",\"hash_key\":\"a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434ad41d8cd98f00b204e9800998ecf8427e10\",\"exist\":true},{\"key\":\"bqohtvxernjmwdka\",\"hash_key\":\"a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434ad41d8cd98f00b204e9800998ecf8427e10\",\"exist\":true},{\"key\":\"jcrxevpmbyauingl\",\"hash_key\":\"a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434ad41d8cd98f00b204e9800998ecf8427e10\",\"exist\":true},{\"key\":\"uhtcvkfpwgaeozbd\",\"hash_key\":\"a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434ad41d8cd98f00b204e9800998ecf8427e10\",\"exist\":true}]"},
    "id": 1
}
```

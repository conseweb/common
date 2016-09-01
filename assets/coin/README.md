## coin
Hyperledger fabric chaincode account model digital currency. Maybe the coin will be called Lepuscoin.

### coin unit
Lepuscoin now has 4 units:
+ tinycoin: the smallest unit of Lepuscoin(tc)
+ minicoin: tinycoin * 1000(mc)
+ smallcoin: minicoin * 1000(sc)
+ coin: smallcoin * 1000(cc)

***Notice***chaincode request(invoke/query) args coin unit always be cc, such as "1.999", means 1.999cc
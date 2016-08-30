# common
The common libraries used by Lepus.

## SUB LIBs
+ captcha: used for generate random digits for verify an email address or a mobile phone
+ clientconn: a gRPC client connection generator
+ config: wapper for viper
+ passphrase: A passphrase is a password that is made up of a handful of words. This makes the password easy to remember and easy to type, even on a phone. Note that the passphrase needs to be longer than a traditional totally random password, because the passphrase is made of words not just random garbage. Inspired by [Bitcoin bip-39](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)
+ protos: all the proto buffer files used by Lepus.
+ semaphore: a semaphore system used for limited concurrency worker number
+ snowflake: a distributed unique id generator inspired by Twitter's Snowflake
+ hdwallet: a wallert & address solution inspired by [Bitcoin bip-32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)
// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/conseweb/common/assets/lepuscoin/txscript"
	"github.com/conseweb/common/hdwallet"
	pb "github.com/conseweb/common/protos"
	"math"
)

// This example demonstrates manually creating and signing a redeem transaction.
func ExampleSignTxOutput() {
	masterKey := hdwallet.MasterKey([]byte("123"))
	addr := masterKey.Address()
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	// For this example, create a fake transaction that represents what
	// would ordinarily be the real transaction that is being spent.  It
	// contains a single output that pays to address in the amount of 1 BTC.
	originTx := pb.NewTx()
	txIn := pb.NewTxIn([]byte(""), math.MaxUint32, []byte{txscript.OP_0, txscript.OP_0})
	originTx.AddTxIn(txIn)
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := pb.NewTxOut(100000000, pkScript)
	originTx.AddTxOut(txOut)
	originTxHash := originTx.TxHash()

	// Create the transaction to redeem the fake transaction.
	redeemTx := pb.NewTx()

	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	txIn = pb.NewTxIn(originTxHash, 0, nil)
	redeemTx.AddTxIn(txIn)

	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut = pb.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)

	// Notice that the script database parameter is nil here since it isn't
	// used.  It must be specified when pay-to-script-hash transactions are
	// being signed.
	sigScript, err := txscript.SignTxOutput(redeemTx, 0, originTx.Txout[0].Script, txscript.SigHashAll, privKey, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	redeemTx.Txin[0].Script = sigScript

	// Prove that the transaction has been validly signed by executing the
	// script pair.
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
		txscript.ScriptStrictMultiSig |
		txscript.ScriptDiscourageUpgradableNops
	vm, err := txscript.NewEngine(originTx.Txout[0].Script, redeemTx, 0, flags, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := vm.Execute(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Transaction successfully signed")

	// Output:
	// Transaction successfully signed
}

// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"

	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	pb "github.com/conseweb/common/protos"
)

// RawTxInSignature returns the serialized ECDSA signature for the input idx of
// the given transaction, with hashType appended to it.
func RawTxInSignature(tx *pb.TX, idx int, subScript []byte, hashType SigHashType, key *ecdsa.PrivateKey) ([]byte, error) {
	parsedScript, err := parseScript(subScript)
	if err != nil {
		return nil, fmt.Errorf("cannot parse output script: %v", err)
	}

	hash := calcSignatureHash(parsedScript, hashType, tx, idx)
	signature, err := key.Sign(rand.Reader, hash, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot sign tx input: %s", err)
	}

	return append(signature, byte(hashType)), nil
}

// SignatureScript creates an input signature script for tx to spend BTC sent
// from a previous output to the owner of privKey. tx must include all
// transaction inputs and outputs, however txin scripts are allowed to be filled
// or empty. The returned script is calculated to be used as the idx'th txin
// sigscript for tx. subscript is the PkScript of the previous output being used
// as the idx'th input. privKey is serialized in either a compressed or
// uncompressed format based on compress. This format must match the same format
// used to generate the payment address, or the script validation will fail.
func SignatureScript(tx *pb.TX, idx int, subscript []byte, hashType SigHashType, privKey *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := RawTxInSignature(tx, idx, subscript, hashType, privKey)
	if err != nil {
		return nil, err
	}

	pkData, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return NewScriptBuilder().AddData(sig).AddData(pkData).Script()
}

func p2pkSignatureScript(tx *pb.TX, idx int, subScript []byte, hashType SigHashType, privKey *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := RawTxInSignature(tx, idx, subScript, hashType, privKey)
	if err != nil {
		return nil, err
	}

	return NewScriptBuilder().AddData(sig).Script()
}

func sign(tx *pb.TX, idx int, subScript []byte, hashType SigHashType, privKey *ecdsa.PrivateKey) ([]byte, error) {
	script, err := SignatureScript(tx, idx, subScript, hashType, privKey)
	if err != nil {
		return nil, err
	}

	return script, nil
}

// mergeScripts merges sigScript and prevScript assuming they are both
// partial solutions for pkScript spending output idx of tx. class, addresses
// and nrequired are the result of extracting the addresses from pkscript.
// The return value is the best effort merging of the two scripts. Calling this
// function with addresses, class and nrequired that do not match pkScript is
// an error and results in undefined behaviour.
func mergeScripts(sigScript, prevScript []byte) []byte {

	// It doesn't actually make sense to merge anything other than multiig
	// and scripthash (because it could contain multisig). Everything else
	// has either zero signature, can't be spent, or has a single signature
	// which is either present or not. The other two cases are handled
	// above. In the conflict case here we just assume the longest is
	// correct (this matches behaviour of the reference implementation).
	if len(sigScript) > len(prevScript) {
		return sigScript
	}
	return prevScript
}

// KeyDB is an interface type provided to SignTxOutput, it encapsulates
// any user state required to get the private keys for an address.
type KeyDB interface {
	GetKey(btcutil.Address) (*btcec.PrivateKey, bool, error)
}

// KeyClosure implements KeyDB with a closure.
type KeyClosure func(btcutil.Address) (*btcec.PrivateKey, bool, error)

// GetKey implements KeyDB by returning the result of calling the closure.
func (kc KeyClosure) GetKey(address btcutil.Address) (*btcec.PrivateKey,
	bool, error) {
	return kc(address)
}

// ScriptDB is an interface type provided to SignTxOutput, it encapsulates any
// user state required to get the scripts for an pay-to-script-hash address.
type ScriptDB interface {
	GetScript(btcutil.Address) ([]byte, error)
}

// ScriptClosure implements ScriptDB with a closure.
type ScriptClosure func(btcutil.Address) ([]byte, error)

// GetScript implements ScriptDB by returning the result of calling the closure.
func (sc ScriptClosure) GetScript(address btcutil.Address) ([]byte, error) {
	return sc(address)
}

// SignTxOutput signs output idx of the given tx to resolve the script given in
// pkScript with a signature type of hashType. Any keys required will be
// looked up by calling getKey() with the string of the given address.
// Any pay-to-script-hash signatures will be similarly looked up by calling
// getScript. If previousScript is provided then the results in previousScript
// will be merged in a type-dependent manner with the newly generated.
// signature script.
func SignTxOutput(tx *pb.TX, idx int, pkScript []byte, hashType SigHashType, privKey *ecdsa.PrivateKey, previousScript []byte) ([]byte, error) {

	sigScript, err := sign(tx, idx, pkScript, hashType, privKey)
	if err != nil {
		return nil, err
	}

	// Merge scripts. with any previous data, if any.
	mergedScript := mergeScripts(sigScript, previousScript)
	return mergedScript, nil
}

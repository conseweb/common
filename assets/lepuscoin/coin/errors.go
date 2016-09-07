/*
Copyright Mojing Inc. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package coin

import "errors"

var (
	// ErrInvalidArgs is returned if there are some unused args or not enough args in params
	ErrInvalidArgs = errors.New("invalid args")
	// ErrInvalidTxKey returned if given key is invalid
	ErrInvalidTxKey = errors.New("invalid tx key")
	// ErrUnsupportedOperation returned if invoke or query using unsupported function name
	ErrUnsupportedOperation = errors.New("unsupported operation")
	// ErrBadEncoding
	ErrBadEncoding = errors.New("bad encoding")
	// ErrParseTx
	ErrParseTx = errors.New("error parseing tx bytes into TX")
	// ErrExecuteUTXO
	ErrExecuteUTXO = errors.New("execute UTXO error")
	// ErrMustCoinbase
	ErrMustCoinbase = errors.New("tx must be coinbase")
)

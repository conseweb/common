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
package main

import "strconv"

type coinunit uint64

const (
	tinycoin  coinunit = 1                // tc
	minicoin           = tinycoin * 1000  // mc
	smallcoin          = minicoin * 1000  // sc
	coin               = smallcoin * 1000 // cc
)

// convert string coin into table coin
func parseCCToTC(strcoin string) (coinunit, error) {
	float64coin, err := strconv.ParseFloat(strcoin, 64)
	if err != nil {
		return 0, err
	}

	return coinunit(float64coin * float64(coin)), nil
}

// convert tinycoin into coin
func convTCToCC(tc coinunit) float64 {
	return float64(tc) / float64(coin)
}

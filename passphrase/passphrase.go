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
package passphrase

import (
	"crypto/rand"
	pb "github.com/conseweb/common/protos"
	"math/big"
	"strings"
)

var (
	wordsStore map[pb.PassphraseLanguage][]string
)

func init() {
	wordsStore = make(map[pb.PassphraseLanguage][]string)

	// english
	wordsStore[pb.PassphraseLanguage_English] = strings.Split(words_en, ",")
	// simplified chinese
	wordsStore[pb.PassphraseLanguage_SimplifiedChinese] = strings.Split(words_zh_SC, ",")
	// traditional chinese
	wordsStore[pb.PassphraseLanguage_TraditionalChinese] = strings.Split(words_zh_TC, ",")
}

// Passphrase return phrases, length is len, language is lang
// If lang is not supported, return length is zero string slice
func Passphrase(length int, lang pb.PassphraseLanguage) (phrase []string) {
	words, ok := wordsStore[lang]
	if !ok {
		return
	}

	max := big.NewInt(int64(len(words)))
	for x := 0; x < length; x++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return
		}

		i := n.Int64()
		phrase = append(phrase, words[i])
	}

	return
}

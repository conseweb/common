/*
Copyright Mojing Inc. 2016 All Rights Reserved.
Written by mint.zhao.chiu@gmail.com. github.com: https://www.github.com/mintzhao

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
	pb "github.com/conseweb/common/protos"
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type PassphraseTest struct {
}

var _ = check.Suite(&PassphraseTest{})

func (t *PassphraseTest) TestPassphraseOK(c *check.C) {
	passphrase, err := Passphrase(256, pb.PassphraseLanguage_English)
	c.Check(err, check.IsNil)
	c.Check(len(passphrase), check.Not(check.Equals), 0)
}

func (t *PassphraseTest) BenchmarkPassphrase(c *check.C) {
	for i := 0; i < c.N; i++ {
		Passphrase(256, pb.PassphraseLanguage_SimplifiedChinese)
	}
}

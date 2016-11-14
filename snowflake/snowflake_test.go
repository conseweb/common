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
package snowflake

import (
	"gopkg.in/check.v1"
	"runtime"
	"testing"
	"time"
)

func TestSnowflake(t *testing.T) {
	check.TestingT(t)
}

type SnowflakeTest struct {
	sf        *Snowflake
	startTime int64
	machineID uint64
}

var _ = check.Suite(&SnowflakeTest{})

func (t *SnowflakeTest) SetUpSuite(c *check.C) {
	ts, err := time.Parse("2006/01/02 15:04:05", "2016/08/16 16:35:00")
	c.Check(err, check.IsNil)
	st := &Settings{
		StartTime: ts,
	}

	t.sf = NewSnowflake(st)
	c.Check(t.sf, check.NotNil)

	t.startTime = toSnowflakeTime(st.StartTime)
	ip, _ := lowerPrivateIP()
	t.machineID = ip
}

func (t *SnowflakeTest) TestNextID(c *check.C) {
	_, err := t.sf.NextID(1, 0)
	c.Check(err, check.IsNil)
}

func (t *SnowflakeTest) TestParseRole(c *check.C) {
	id, err := t.sf.NextID(5, 0)
	c.Check(err, check.IsNil)
	c.Check(ParseRole(id), check.Equals, uint64(5))
}

func (t *SnowflakeTest) BenchmarkNextID(c *check.C) {
	for i := 0; i <= c.N; i++ {
		t.sf.NextID(1, 0)
	}
}

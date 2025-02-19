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
package exec

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/op/go-logging"
)

var (
	execLogger = logging.MustGetLogger("exec")
)

type SignalExecuter func() error

// HandleSignal handle os signals
func HandleSignal(fs ...SignalExecuter) {
	execLogger.Info("waitting for exit, press CTRL+C")

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	wg := &sync.WaitGroup{}
	wg.Add(len(fs))

	for {
		s := <-sigChan

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			for _, f := range fs {
				if err := f(); err != nil {
					execLogger.Errorf("exec stop function return error: %v", err)
				}
				wg.Done()
			}

			wg.Wait()
			os.Exit(0)
		}
	}
}

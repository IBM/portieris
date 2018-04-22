// Copyright 2018 IBM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package framework

import (
	"bufio"
	"fmt"
	"os/exec"
)

func (f *Framework) installChart() error {
	cmdName := "helm"
	cmdArgs := []string{"install", f.HelmChart, "-n", f.HelmRelease}
	return f.runCommand(cmdName, cmdArgs)
}

func (f *Framework) deleteChart() error {
	cmdName := "helm"
	cmdArgs := []string{"delete", "--purge", f.HelmRelease}
	return f.runCommand(cmdName, cmdArgs)
}

func (f *Framework) runCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe for helm install: %v", err)
	}
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe for helm install: %v", err)
	}

	commandLog := fmt.Sprintf("%v %v", name, args[0])
	outScanner := bufio.NewScanner(cmdOut)
	go func() {
		for outScanner.Scan() {
			fmt.Printf("%s\n", outScanner.Text())
		}
	}()

	errScanner := bufio.NewScanner(cmdErr)
	go func() {
		for errScanner.Scan() {
			fmt.Printf("%v | %s\n", commandLog, errScanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting helm install: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error waiting for helm install: %v", err)
	}
	return nil
}

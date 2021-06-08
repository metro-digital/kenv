//
// Copyright 2021 METRO Digital GmbH
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
//
// +build unitTests

package env_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

// more details: https://npf.io/2015/06/testing-exec-command/
var testCase string

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	tc := "TEST_CASE=" + testCase
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", tc}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	switch os.Getenv("TEST_CASE") {
	case "case1":
		fmt.Fprintf(os.Stderr, `[GNUPG:] IMPORT_OK 0 E341BE2FA255F0469E952011BD09B4A7A2A2A496
		[GNUPG:] KEY_CONSIDERED E341BE2FA255F0469E952011BD09B4A7A2A2A496 0
		gpg: key BD09B4A7A2A2A496: "bb-example <metro.digital@example.com>" not changed
		[GNUPG:] KEY_CONSIDERED E341BE2FA255F0469E952011BD09B4A7A2A2A496 0
		gpg: key BD09B4A7A2A2A496: secret key imported
		[GNUPG:] IMPORT_OK 16 E341BE2FA255F0469E952011BD09B4A7A2A2A496
		gpg: Total number processed: 1
		gpg:              unchanged: 1
		gpg:       secret keys read: 1
		gpg:  secret keys unchanged: 1
		[GNUPG:] IMPORT_RES 1 0 0 0 1 0 0 0 0 1 0 1 0 0 0`)
		content, _ := ioutil.ReadFile("fixtures/kustomize.output")
		fmt.Fprintf(os.Stdout, string(content))
	case "case2":
		fmt.Fprintf(os.Stderr, "gpg: invalid option")
	}
}

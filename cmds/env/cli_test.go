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
  "io/ioutil"
  "os"
  "os/exec"
  "sort"
  "strings"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/spf13/cobra"

  "github.com/metro-digital/kenv/cmds/env"
  "github.com/metro-digital/kenv/services"
)

var _ = Describe("PrepareRun()", func() {
  var cmd *cobra.Command

  BeforeEach(func() {
    os.Remove("vars")

    cmd = env.Cli()
    cmd.SetArgs([]string{"--input-directory", "waas-config/environments/be-gcw1/pp", "--key", "1234", "--output", "vars"})

    services.ExecCommand = fakeExecCommand
  })

  AfterEach(func() {
    os.Remove("vars")

    services.ExecCommand = exec.Command
  })

  Context("when command is executed and everyting is ok", func() {
    It("creates a dotenv file", func() {
      testCase = "case1"

      cmd.Execute()

      content, _ := ioutil.ReadFile("vars")
      vars := strings.Split(string(content), "\n")
      sort.Strings(vars)

      Expect(len(vars)).To(Equal(9))
      Expect(vars[1]).To(Equal("BACKGROUND_COLOR=green"))
    })
  })

  Context("when command is executed but key cannot be imported", func() {
    It("exits with error", func() {
      testCase = "case2"

      var errMsg string

      func() {
        defer func() {
          r := recover()
          err, _ := r.(error)
          errMsg = err.Error()
        }()
        cmd.Execute()
      }()

      Expect(errMsg).To(Equal("unable to import key"))
    })
  })
})

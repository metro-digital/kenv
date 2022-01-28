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

package services_test

import (
  "os/exec"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "github.com/metro-digital/kenv/services"
)

var _ = Describe("GetVars()", func() {
  var kustomize *services.Kustomize

  BeforeEach(func() {
    services.ExecCommand = fakeExecCommand
    kustomize, _ = services.InitKustomize()
  })

  AfterEach(func() {
    services.ExecCommand = exec.Command
  })

  Context("when function is executed with correct input", func() {
    It("returns extracted environment variables", func() {
      testCase = "case6"

      vars, err := kustomize.GetVars("waas-config/environments/be-gcw1/pp")

      Expect(err).NotTo(HaveOccurred())
      Expect(len(vars)).To(Equal(2))
      Expect(vars[0].Kind).To(Equal("ConfigMap"))
      Expect(vars[0].Data["DEPLOYMENT_DATACENTER"]).To(Equal("be-gcw1"))
      Expect(vars[1].Kind).To(Equal("Secret"))
      Expect(vars[1].Data["BACKGROUND_COLOR"]).To(Equal("green"))
    })
  })

  Context("when function is executed but gpg key was not imported", func() {
    It("returns nothing", func() {
      testCase = "case7"

      vars, err := kustomize.GetVars("waas-config/environments/be-gcw1/pp")

      Expect(err).NotTo(HaveOccurred())
      Expect(len(vars)).To(Equal(0))
    })
  })
})

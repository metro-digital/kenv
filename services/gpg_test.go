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

var _ = Describe("ImportKey()", func() {
	var gpg *services.Gpg

	BeforeEach(func() {
		services.ExecCommand = fakeExecCommand
		gpg, _ = services.InitGpg()
	})

	AfterEach(func() {
		services.ExecCommand = exec.Command
	})

	Context("when function is executed with correct input", func() {
		It("returns valid keyId", func() {
			testCase = "case1"

			keyId, err := gpg.ImportKey("1234")

			Expect(err).NotTo(HaveOccurred())
			Expect(keyId).To(Equal("E341BE2FA255F0469E952011BD09B4A7A2A2A496"))
		})
	})

	Context("when function is executed with correct input, but there is an error", func() {
		It("returns error", func() {
			testCase = "case2"

			_, err := gpg.ImportKey("1234")

			Expect(err).To(HaveOccurred())
		})
	})

	Context("when function is executed with broken input", func() {
		It("returns error", func() {
			testCase = "case3"

			_, err := gpg.ImportKey("1234")

			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("DeleteKey()", func() {
	var gpg *services.Gpg

	BeforeEach(func() {
		services.ExecCommand = fakeExecCommand
		gpg, _ = services.InitGpg()
	})

	AfterEach(func() {
		services.ExecCommand = exec.Command
	})

	Context("when function is executed with correct input", func() {
		It("deletes the key", func() {
			testCase = "case4"

			err := gpg.DeleteKey("E341BE2FA255F0469E952011BD09B4A7A2A2A496")

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when function is executed but key does not exist", func() {
		It("returns error", func() {
			testCase = "case5"

			_, err := gpg.ImportKey("E341BE2FA255F0469E952011BD09B4A7A2A2A496")

			Expect(err).To(HaveOccurred())
		})
	})
})

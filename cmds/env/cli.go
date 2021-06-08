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

package env

import (
	"log"
	"os"

	"github.com/metro-digital/kenv/services"
	"github.com/spf13/cobra"
)

const inputFlag = "input-directory"
const keyFlag = "key"
const outputFlag = "output"

func Cli() *cobra.Command {
	prepare := &cobra.Command{
		Use:   "prepare",
		Short: "Prepare environment variables (dotenv file) based on the kustomize output.",
		Args:  cobra.NoArgs,
		Run:   prepareRun,
	}

	prepare.Flags().StringP(inputFlag, "i", "", "Input directory for stage (e.g. waas-config/environments/be-gcw1/pp)")
	prepare.Flags().StringP(keyFlag, "k", "", "Private GPG key, base64 encoded")
	prepare.Flags().StringP(outputFlag, "o", "", "Output dotenv file")

	return prepare
}

func prepareRun(cmd *cobra.Command, args []string) {
	input, err := cmd.Flags().GetString(inputFlag)
	check(err)

	key, err := cmd.Flags().GetString(keyFlag)
	check(err)

	output, err := cmd.Flags().GetString(outputFlag)
	check(err)

	if input == "" || key == "" || output == "" {
		log.Fatal("All parameters are mandatory")
	}

	gpg, err := services.InitGpg()
	check(err)

	keyID, err := gpg.ImportKey(key)
	check(err)

	kustomize, err := services.InitKustomize()
	check(err)

	docs, err := kustomize.GetVars(input)
	check(err)

	f, err := os.Create(output)
	check(err)

	defer f.Close()

	for _, doc := range docs {
		for k, v := range doc.Data {
			_, err := f.WriteString(k + "=" + v + "\n")
			check(err)
		}
	}

	err = f.Sync()
	check(err)

	// nolint:errcheck
	defer gpg.DeleteKey(keyID)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

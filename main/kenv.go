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

package main

import (
	"github.com/metro-digital/kenv/cmds/env"
	"github.com/spf13/cobra"
)

func initApp() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kenv",
		Short: "kenv is a utility for preparing dotenv file from kustomize config.",
		Args:  cobra.NoArgs,
	}

	rootCmd.AddCommand(env.Cli())

	return rootCmd
}

func main() {
	app := initApp()
	// nolint:errcheck
	app.Execute()
}

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
package services

import (
  "bytes"
  "encoding/base64"
  "errors"
  "fmt"
  "os/exec"

  "go.uber.org/multierr"
  "gopkg.in/yaml.v3"
)

type Vars struct {
  Kind string            `yaml:"kind"`
  Data map[string]string `yaml:"data"`
}

type Kustomize struct {
  Binary string
}

func InitKustomize() (*Kustomize, error) {
  path, err := exec.LookPath("kustomize")

  if err != nil || path == "" {
    return nil, errors.New("kustomize binary not found")
  }

  return &Kustomize{
    Binary: path,
  }, nil
}

func (k *Kustomize) GetVars(input string) ([]Vars, error) {
  args := []string{"build", "--enable-alpha-plugins", input}
  cmd := ExecCommand("kustomize", args...)

  var stdout, stderr bytes.Buffer
  cmd.Stdout = &stdout
  cmd.Stderr = &stderr
  err := cmd.Run()

  allOutput := stdout.String()
  docs := parse([]byte(allOutput))
  if err != nil {
    err = multierr.Combine(
      errors.New(stderr.String()),
      errors.New(fmt.Sprint("kustomize failed to run: ", err)),
    )
  }

  return docs, err
}

func parse(source []byte) []Vars {
  result := []Vars{}

  dec := yaml.NewDecoder(bytes.NewReader(source))
  for {
    var doc Vars
    if dec.Decode(&doc) != nil {
      break
    }

    if doc.Kind == "ConfigMap" {
      result = append(result, doc)
    }

    if doc.Kind == "Secret" {
      tmp := Vars{Kind: doc.Kind, Data: map[string]string{}}
      for k, v := range doc.Data {
        decoded, err := base64.StdEncoding.DecodeString(v)
        if err != nil {
          panic(err)
        }
        tmp.Data[k] = string(decoded)
      }
      result = append(result, tmp)
    }
  }

  return result
}

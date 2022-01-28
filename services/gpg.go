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
  "os"
  "os/exec"
  "strings"
)

type Gpg struct {
  Binary  string
  Homedir string
}

type OutputChunk struct {
  Key  string
  Text string
}

var ExecCommand = exec.Command

func InitGpg() (*Gpg, error) {
  path := os.Getenv("GNUPG_BIN")
  var err error = nil
  if path == "" {
    path, err = exec.LookPath("gpg")
  }
  if err != nil || path == "" {
    return nil, errors.New("gpg binary not found")
  }

  gpg := new(Gpg)
  gpg.Binary = path
  gpg.Homedir = os.Getenv("GNUPGHOME")
  if gpg.Homedir == "" {
    gpg.Homedir = "~/.gnupg"
  }
  return gpg, nil
}

func (gpg *Gpg) ImportKey(key string) (string, error) {
  decoded, err := base64.StdEncoding.DecodeString(key)
  if err != nil {
    return "", err
  }

  chunks, err := gpg.execCommand([]string{"--import"}, string(decoded))
  keyID := ""

  for _, chunk := range chunks {
    if chunk.Key == "IMPORT_OK" {
      keyID = strings.Split(chunk.Text, " ")[1]
      break
    }
  }
  if keyID == "" {
    if err != nil {
      return "", err
    }
    return "", errors.New("unable to import key")
  }

  return keyID, nil
}

func (gpg *Gpg) DeleteKey(keyID string) error {
  args := append([]string{"--batch", "--yes", "--delete-secret-keys"}, keyID)
  chunks, err := gpg.execCommand(args, "")

  for _, chunk := range chunks {
    if chunk.Key == "DELETE_PROBLEM" {
      return errors.New("unable to delete the key")
    }
  }

  return err
}

func (gpg *Gpg) execCommand(commands []string, input string) ([]OutputChunk, error) {
  args := append([]string{
    "--status-fd", "2",
    "--no-tty",
    "--homedir", gpg.Homedir,
  }, commands...)
  cmd := ExecCommand(gpg.Binary, args...)

  if len(input) > 0 {
    cmd.Stdin = strings.NewReader(input)
  }

  var stderr bytes.Buffer
  cmd.Stderr = &stderr
  err := cmd.Run()

  allOutput := stderr.String()
  lines := strings.Split(allOutput, "\n")

  var chunks = []OutputChunk{}
  for _, line := range lines {
    toks := strings.SplitN(line, " ", 3)
    if toks[0] != "[GNUPG:]" {
      continue
    }
    chunk := OutputChunk{toks[1], ""}
    if len(toks) == 3 {
      chunk.Text = toks[2]
    }
    chunks = append(chunks, chunk)
  }

  if err != nil {
    err = errors.New(fmt.Sprint("gpg failed to run: ", err))
  }

  return chunks, err
}

// Copyright 2020 Red Hat, Inc
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
// limitations under the License.)

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/coreos/ignition/config/v2_4"
	"github.com/coreos/ignition/v2/config/shared/errors"
	"github.com/coreos/ignition/v2/config/v3_0"
	"github.com/coreos/ignition/v2/config/v3_1"
	translateTo3_1 "github.com/coreos/ignition/v2/config/v3_1/translate"

	"github.com/coreos/ign-converter/translate/v24tov31"
	"github.com/coreos/ign-converter/translate/v31tov24"
)

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func getMapping(fname string) map[string]string {
	m := map[string]string{}
	if fname == "" {
		return m
	}
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		fail("Error reading %s: %v", fname, err)
	}
	// parse
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			fail("Error parsing line: %q, needs two parts", line)
		}
		m[parts[0]] = parts[1]
	}
	return m
}

func main() {
	var (
		input         string
		output        string
		fsMap         string
		versionFlag   bool
		downtranslate bool
	)
	flag.BoolVar(&versionFlag, "version", false, "print the version and exit")
	flag.StringVar(&input, "input", "", "read from input file instead of stdin")
	flag.StringVar(&fsMap, "fsmap", "", "file containing mapping from filesystem name to path")
	flag.StringVar(&output, "output", "", "write to output file instead of stdout")
	flag.BoolVar(&downtranslate, "downtranslate", false, "translate a spec 3 config down to spec 2")

	flag.Parse()

	if versionFlag {
		fmt.Println("todo: add version")
		os.Exit(0)
	}

	var infile *os.File = os.Stdin
	var outfile *os.File = os.Stdout
	if input != "" {
		var err error
		infile, err = os.Open(input)
		if err != nil {
			fail("failed to open %s: %v\n", input, err)
		}
		defer infile.Close()
	}

	dataIn, err := ioutil.ReadAll(infile)
	if err != nil {
		fail("failed to read %s: %v\n", infile.Name(), err)
	}

	var dataOut []byte
	if downtranslate {
		// translate from 3.x to 2.4
		cfg, rpt, err := v3_1.Parse(dataIn)
		fmt.Fprintf(os.Stderr, "%s", rpt.String())
		if err == errors.ErrUnknownVersion {
			// We can fail unmarshaling if it's an older config. Attempt to parse
			// it as such.
			cfg3_0, rpt, err := v3_0.Parse(dataIn)
			fmt.Fprintf(os.Stderr, "%s", rpt.String())
			if err != nil || rpt.IsFatal() {
				fail("Error parsing spec v3.0 config: %v\n%v", err, rpt)
			}
			cfg = translateTo3_1.Translate(cfg3_0)
		} else if err != nil || rpt.IsFatal() {
			fail("Error parsing spec v3.1 config: %v\n%v", err, rpt)
		}

		newCfg, err := v31tov24.Translate(cfg)
		if err != nil {
			fail("Failed to translate config from 3 to 2: %v", err)
		}
		dataOut, err = json.Marshal(newCfg)
		if err != nil {
			fail("Failed to marshal json: %v", err)
		}
	} else {
		// translate from 2.x to 3.1
		mapping := getMapping(fsMap)

		// parse
		cfg, rpt, err := v2_4.Parse(dataIn)
		fmt.Fprintf(os.Stderr, "%s", rpt.String())
		if err != nil || rpt.IsFatal() {
			fail("Error parsing spec v2 config: %v\n%v", err, rpt)
		}

		newCfg, err := v24tov31.Translate(cfg, mapping)
		if err != nil {
			fail("Failed to translate config from 2 to 3: %v", err)
		}
		dataOut, err = json.Marshal(newCfg)
		if err != nil {
			fail("Failed to marshal json: %v", err)
		}
	}

	if output != "" {
		var err error
		outfile, err = os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fail("failed to open %s: %v\n", output, err)
		}
		defer outfile.Close()
	}

	if _, err := outfile.Write(dataOut); err != nil {
		fail("Failed to write config to %s: %v\n", outfile.Name(), err)
	}
}

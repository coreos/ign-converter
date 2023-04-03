// Copyright 2020 Red Hat, Inc.
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

package v33tov32

import (
	"fmt"
	"testing"

	"github.com/coreos/ign-converter/util"
	types3_2 "github.com/coreos/ignition/v2/config/v3_2/types"
	types3_3 "github.com/coreos/ignition/v2/config/v3_3/types"
	"github.com/stretchr/testify/assert"
)

func TestTranslate3_3to3_23(t *testing.T) {
	type in struct {
		data types3_3.Config
	}

	type out struct {
		data types3_2.Config
	}

	tests := []struct {
		in  in
		out out
		err error
	}{
		{
			in: in{data: types3_3.Config{
				Ignition: types3_3.Ignition{
					Version: "3.3.0",
				},
			}},
			out: out{data: types3_2.Config{
				Ignition: types3_2.Ignition{
					Version: "3.2.0",
				},
			}},
		},
		{
			in:  in{data: util.NonexhaustiveConfig3_3},
			out: out{data: util.NonexhaustiveConfig3_2},
		},
		{
			in: in{data: types3_3.Config{
				Ignition: types3_3.Ignition{
					Version: "3.3.0",
				},
				KernelArguments: types3_3.KernelArguments{
					ShouldExist: []types3_3.KernelArgument{"foo"},
				},
			}},
			out: out{data: types3_2.Config{}},
			err: fmt.Errorf("KernelArguments is not supported on 3.2"),
		},
		{

			in: in{data: types3_3.Config{
				Ignition: types3_3.Ignition{
					Version: "3.3.0",
				},
				KernelArguments: types3_3.KernelArguments{
					ShouldNotExist: []types3_3.KernelArgument{"foo"},
				},
			}},
			out: out{data: types3_2.Config{}},
			err: fmt.Errorf("KernelArguments is not supported on 3.2"),
		},
		{

			in: in{data: types3_3.Config{
				Ignition: types3_3.Ignition{
					Version: "3.3.0",
				},
				Storage: types3_3.Storage{
					Filesystems: []types3_3.Filesystem{
						{
							Device: "/dev/foo",
							Format: util.StrP("util.None"),
						},
					},
				},
			}},
			out: out{data: types3_2.Config{}},
			err: fmt.Errorf("Invalid input config:\nerror at $.storage.filesystems.0.format: invalid filesystem format\n"),
		},
	}
	for i, test := range tests {
		result, err := Translate(test.in.data)
		assert.Equal(t, test.err, err, "#%d: bad error: want %v, got %v", i, test.err, err)
		assert.Equal(t, test.out.data, result, "#%d: bad result, want %v, got %v", i, test.out.data, result)
	}

}

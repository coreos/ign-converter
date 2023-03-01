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

package v32tov31

import (
	"fmt"
	"testing"

	"github.com/coreos/ign-converter/util"
	types3_1 "github.com/coreos/ignition/v2/config/v3_1/types"
	types3_2 "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/stretchr/testify/assert"
)

func TestTranslate3_2to3_1(t *testing.T) {
	type in struct {
		data types3_2.Config
	}

	type out struct {
		data types3_1.Config
	}

	tests := []struct {
		in  in
		out out
		err error
	}{
		{
			in: in{data: types3_2.Config{
				Ignition: types3_2.Ignition{
					Version: "3.2.0",
				},
			}},
			out: out{data: types3_1.Config{
				Ignition: types3_1.Ignition{
					Version: "3.1.0",
				},
			}},
		},
		{
			in:  in{data: util.NonexhaustiveConfig3_2},
			out: out{data: util.NonexhaustiveConfig3_1},
		},
		{
			in: in{data: types3_2.Config{
				Ignition: types3_2.Ignition{
					Version: "3.2.0",
				},
				Storage: types3_2.Storage{
					Luks: []types3_2.Luks{
						{
							Name:   "z",
							Device: util.StrP("/dev/z"),
						},
					},
				},
			},
			},
			out: out{data: types3_1.Config{}},
			err: fmt.Errorf("LUKS is not supported on 3.1"),
		},
		{
			in: in{data: types3_2.Config{
				Ignition: types3_2.Ignition{
					Version: "3.2.0",
				},
				Storage: types3_2.Storage{
					Disks: []types3_2.Disk{
						{
							Device: "/dev/a",
							Partitions: []types3_2.Partition{
								{
									Label:  util.StrP("z"),
									Resize: util.BoolP(true),
								},
							},
						},
					},
				},
			},
			},
			out: out{data: types3_1.Config{}},
			err: fmt.Errorf("Resize in Storage.Disks.Partitions is not supported on 3.1"),
		},
		{
			in: in{data: types3_2.Config{
				Ignition: types3_2.Ignition{
					Version: "3.2.0",
				},
				Passwd: types3_2.Passwd{
					Users: []types3_2.PasswdUser{
						{
							Name:        "z",
							ShouldExist: util.BoolPStrict(false),
						},
					},
				},
			},
			},
			out: out{data: types3_1.Config{}},
			err: fmt.Errorf("ShouldExist in Passwd.Users is not supported on 3.1"),
		},
		{
			in: in{data: types3_2.Config{
				Ignition: types3_2.Ignition{
					Version: "3.2.0",
				},
				Passwd: types3_2.Passwd{
					Groups: []types3_2.PasswdGroup{
						{
							Name:        "z",
							ShouldExist: util.BoolPStrict(false),
						},
					},
				},
			},
			},
			out: out{data: types3_1.Config{}},
			err: fmt.Errorf("ShouldExist in Passwd.Groups is not supported on 3.1"),
		},
	}

	for i, test := range tests {
		result, err := Translate(test.in.data)
		assert.Equal(t, test.err, err, "#%d: bad error: want %v, got %v", i, test.err, err)
		assert.Equal(t, test.out.data, result, "#%d: bad result, want %v, got %v", i, test.out.data, result)
	}
}

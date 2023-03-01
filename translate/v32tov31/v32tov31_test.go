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
	"testing"

	"github.com/coreos/ign-converter/util"
	types3_2 "github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/stretchr/testify/assert"
)

func TestTranslate3_2to3_1(t *testing.T) {
	emptyConfig := types3_2.Config{
		Ignition: types3_2.Ignition{
			Version: "3.2.0",
		},
	}

	_, err := Translate(emptyConfig)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}

	res, err := Translate(util.NonexhaustiveConfig3_2)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, util.NonexhaustiveConfig3_1, res)

	_, err = Translate(types3_2.Config{
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
	})
	assert.Error(t, err)
	_, err = Translate(types3_2.Config{
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
	})
	assert.Error(t, err)
	_, err = Translate(types3_2.Config{
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
	})
	assert.Error(t, err)
	_, err = Translate(types3_2.Config{
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
	})
	assert.Error(t, err)
}

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
	"testing"

	"github.com/coreos/ign-converter/util"

	types3_3 "github.com/coreos/ignition/v2/config/v3_3/types"

	"github.com/stretchr/testify/assert"
)

func TestTranslate3_3to3_2(t *testing.T) {
	emptyConfig := types3_3.Config{
		Ignition: types3_3.Ignition{
			Version: "3.3.0",
		},
	}

	_, err := Translate(emptyConfig)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}

	res, err := Translate(util.NonexhaustiveConfig3_3)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, util.NonexhaustiveConfig3_2, res)

	_, err = Translate(types3_3.Config{
		Ignition: types3_3.Ignition{
			Version: "3.3.0",
		},
		KernelArguments: types3_3.KernelArguments{
			ShouldExist: []types3_3.KernelArgument{"foo"},
		},
	})
	assert.Error(t, err)
	_, err = Translate(types3_3.Config{
		Ignition: types3_3.Ignition{
			Version: "3.3.0",
		},
		KernelArguments: types3_3.KernelArguments{
			ShouldNotExist: []types3_3.KernelArgument{"foo"},
		},
	})
	assert.Error(t, err)
	_, err = Translate(types3_3.Config{
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
	})
	assert.Error(t, err)
}

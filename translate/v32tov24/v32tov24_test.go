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

package v32tov24

import (
	"testing"

	"github.com/coreos/ign-converter/util"
	types3_2 "github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/stretchr/testify/assert"
)

func TestTranslate3_2to2_4(t *testing.T) {
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
	assert.Equal(t, util.ExhaustiveConfig2_4, res)
}

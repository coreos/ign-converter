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

package v23tov30

import (
	"testing"

	"github.com/coreos/ign-converter/util"
	types2_3 "github.com/coreos/ignition/config/v2_3/types"
	"github.com/stretchr/testify/assert"
)

// Configs using _all_ the (undeprecated) fields
type input2_3 struct {
	cfg   types2_3.Config
	fsMap map[string]string
}

func TestCheck2_3(t *testing.T) {
	goodConfigs := []input2_3{
		{
			util.ExhaustiveConfig2_3,
			util.ExhaustiveMap,
		},
	}
	badConfigs := []input2_3{
		{}, // empty config has no version, fails validation
		{
			// need a map for filesystems
			util.ExhaustiveConfig2_3,
			nil,
		},
	}
	for i, e := range goodConfigs {
		if err := Check2_3(e.cfg, e.fsMap); err != nil {
			t.Errorf("Good config test %d: got %v, expected nil", i, err)
		}
	}
	for i, e := range badConfigs {
		if err := Check2_3(e.cfg, e.fsMap); err == nil {
			t.Errorf("Bad config test %d: got ok, expected: %v", i, err)
		}
	}
}

func TestTranslate2_3to3_0(t *testing.T) {
	res, err := Translate(util.ExhaustiveConfig2_3, util.ExhaustiveMap)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, util.ExhaustiveConfig3_0, res)
}

func TestRemoveDuplicateFilesUnitsUsers2_3(t *testing.T) {
	mode := 420
	testDataOld := "data:,old"
	testDataNew := "data:,new"
	testIgn2Config := types2_3.Config{}

	// file test, add a duplicate file and see if the newest one is preserved
	fileOld := types2_3.File{
		Node: types2_3.Node{
			Filesystem: "root", Path: "/etc/testfileconfig",
		},
		FileEmbedded1: types2_3.FileEmbedded1{
			Contents: types2_3.FileContents{
				Source: testDataOld,
			},
			Mode: &mode,
		},
	}
	testIgn2Config.Storage.Files = append(testIgn2Config.Storage.Files, fileOld)

	fileNew := types2_3.File{
		Node: types2_3.Node{
			Filesystem: "root", Path: "/etc/testfileconfig",
		},
		FileEmbedded1: types2_3.FileEmbedded1{
			Contents: types2_3.FileContents{
				Source: testDataNew,
			},
			Mode: &mode,
		},
	}
	testIgn2Config.Storage.Files = append(testIgn2Config.Storage.Files, fileNew)

	// unit test, add three units and three dropins with the same name as follows:
	// unitOne:
	//    contents: old
	//    dropin:
	//        name: one
	//        contents: old
	// unitTwo:
	//    dropin:
	//        name: one
	//        contents: new
	// unitThree:
	//    contents: new
	//    dropin:
	//        name: two
	//        contents: new
	// Which should result in:
	// unitFinal:
	//    contents: new
	//    dropin:
	//      - name: one
	//        contents: new
	//      - name: two
	//        contents: new
	//
	unitName := "testUnit"
	dropinNameOne := "one"
	dropinNameTwo := "two"
	dropinOne := types2_3.SystemdDropin{
		Contents: testDataOld,
		Name:     dropinNameOne,
	}
	dropinTwo := types2_3.SystemdDropin{
		Contents: testDataNew,
		Name:     dropinNameOne,
	}
	dropinThree := types2_3.SystemdDropin{
		Contents: testDataNew,
		Name:     dropinNameTwo,
	}

	unitOne := types2_3.Unit{
		Contents: testDataOld,
		Name:     unitName,
	}
	unitOne.Dropins = append(unitOne.Dropins, dropinOne)
	testIgn2Config.Systemd.Units = append(testIgn2Config.Systemd.Units, unitOne)

	unitTwo := types2_3.Unit{
		Name: unitName,
	}
	unitTwo.Dropins = append(unitTwo.Dropins, dropinTwo)
	testIgn2Config.Systemd.Units = append(testIgn2Config.Systemd.Units, unitTwo)

	unitThree := types2_3.Unit{
		Contents: testDataNew,
		Name:     unitName,
	}
	unitThree.Dropins = append(unitThree.Dropins, dropinThree)
	testIgn2Config.Systemd.Units = append(testIgn2Config.Systemd.Units, unitThree)

	// user test, add a duplicate user and see if it is deduplicated but ssh keys from both are preserved
	userName := "testUser"
	userOne := types2_3.PasswdUser{
		Name: userName,
		SSHAuthorizedKeys: []types2_3.SSHAuthorizedKey{
			"one",
			"two",
		},
	}
	userTwo := types2_3.PasswdUser{
		Name: userName,
		SSHAuthorizedKeys: []types2_3.SSHAuthorizedKey{
			"three",
		},
	}
	userThree := types2_3.PasswdUser{
		Name: "userThree",
		SSHAuthorizedKeys: []types2_3.SSHAuthorizedKey{
			"four",
		},
	}
	testIgn2Config.Passwd.Users = append(testIgn2Config.Passwd.Users, userOne, userTwo, userThree)

	convertedIgn2Config, err := RemoveDuplicateFilesUnitsUsers(testIgn2Config)
	assert.NoError(t, err)

	expectedIgn2Config := types2_3.Config{}
	expectedIgn2Config.Storage.Files = append(expectedIgn2Config.Storage.Files, fileNew)
	unitExpected := types2_3.Unit{
		Contents: testDataNew,
		Name:     unitName,
	}
	unitExpected.Dropins = append(unitExpected.Dropins, dropinThree)
	unitExpected.Dropins = append(unitExpected.Dropins, dropinTwo)
	expectedIgn2Config.Systemd.Units = append(expectedIgn2Config.Systemd.Units, unitExpected)
	expectedMergedUser := types2_3.PasswdUser{
		Name: userName,
		SSHAuthorizedKeys: []types2_3.SSHAuthorizedKey{
			"three",
			"one",
			"two",
		},
	}
	expectedIgn2Config.Passwd.Users = append(expectedIgn2Config.Passwd.Users, userThree, expectedMergedUser)

	assert.Equal(t, expectedIgn2Config, convertedIgn2Config)
}

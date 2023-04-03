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

package v24tov31

import (
	"testing"

	"github.com/coreos/ign-converter/util"
	types2_4 "github.com/coreos/ignition/config/v2_4/types"
	"github.com/stretchr/testify/assert"
)

type input2_4 struct {
	cfg   types2_4.Config
	fsMap map[string]string
}

func TestCheck2_4(t *testing.T) {
	goodConfigs := []input2_4{
		{
			util.ExhaustiveConfig2_4,
			util.ExhaustiveMap,
		},
	}
	badConfigs := []input2_4{
		{}, // empty config has no version, fails validation
		{
			// need a map for filesystems
			util.ExhaustiveConfig2_4,
			nil,
		},
	}
	for i, e := range goodConfigs {
		if err := Check2_4(e.cfg, e.fsMap); err != nil {
			t.Errorf("Good config test %d: got %v, expected nil", i, err)
		}
	}
	for i, e := range badConfigs {
		if err := Check2_4(e.cfg, e.fsMap); err == nil {
			t.Errorf("Bad config test %d: got ok, expected: %v", i, err)
		}
	}
}
func TestTranslate2_4to3_1(t *testing.T) {
	res, err := Translate(util.ExhaustiveConfig2_4, util.ExhaustiveMap)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, util.NonexhaustiveConfig3_1, res)
}
func TestRemoveDuplicateFilesUnitsUsers2_4(t *testing.T) {
	mode := 420
	testDataOld := "data:,old"
	testDataNew := "data:,new"
	testIgn2Config := types2_4.Config{}

	// file test, add a duplicate file and see if the newest one is preserved
	fileOld := types2_4.File{
		Node: types2_4.Node{
			Filesystem: "root", Path: "/etc/testfileconfig",
		},
		FileEmbedded1: types2_4.FileEmbedded1{
			Contents: types2_4.FileContents{
				Source: testDataOld,
			},
			Mode: &mode,
		},
	}
	testIgn2Config.Storage.Files = append(testIgn2Config.Storage.Files, fileOld)

	fileNew := types2_4.File{
		Node: types2_4.Node{
			Filesystem: "root", Path: "/etc/testfileconfig",
		},
		FileEmbedded1: types2_4.FileEmbedded1{
			Contents: types2_4.FileContents{
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
	dropinOne := types2_4.SystemdDropin{
		Contents: testDataOld,
		Name:     dropinNameOne,
	}
	dropinTwo := types2_4.SystemdDropin{
		Contents: testDataNew,
		Name:     dropinNameOne,
	}
	dropinThree := types2_4.SystemdDropin{
		Contents: testDataNew,
		Name:     dropinNameTwo,
	}

	unitOne := types2_4.Unit{
		Contents: testDataOld,
		Name:     unitName,
	}
	unitOne.Dropins = append(unitOne.Dropins, dropinOne)
	testIgn2Config.Systemd.Units = append(testIgn2Config.Systemd.Units, unitOne)

	unitTwo := types2_4.Unit{
		Name: unitName,
	}
	unitTwo.Dropins = append(unitTwo.Dropins, dropinTwo)
	testIgn2Config.Systemd.Units = append(testIgn2Config.Systemd.Units, unitTwo)

	unitThree := types2_4.Unit{
		Contents: testDataNew,
		Name:     unitName,
	}
	unitThree.Dropins = append(unitThree.Dropins, dropinThree)
	testIgn2Config.Systemd.Units = append(testIgn2Config.Systemd.Units, unitThree)

	// user test, add a duplicate user and see if it is deduplicated but ssh keys from both are preserved
	userName := "testUser"
	userOne := types2_4.PasswdUser{
		Name: userName,
		SSHAuthorizedKeys: []types2_4.SSHAuthorizedKey{
			"one",
			"two",
		},
	}
	userTwo := types2_4.PasswdUser{
		Name: userName,
		SSHAuthorizedKeys: []types2_4.SSHAuthorizedKey{
			"three",
		},
	}
	userThree := types2_4.PasswdUser{
		Name: "userThree",
		SSHAuthorizedKeys: []types2_4.SSHAuthorizedKey{
			"four",
		},
	}
	testIgn2Config.Passwd.Users = append(testIgn2Config.Passwd.Users, userOne, userTwo, userThree)

	convertedIgn2Config, err := RemoveDuplicateFilesUnitsUsers(testIgn2Config)
	assert.NoError(t, err)

	expectedIgn2Config := types2_4.Config{}
	expectedIgn2Config.Storage.Files = append(expectedIgn2Config.Storage.Files, fileNew)
	unitExpected := types2_4.Unit{
		Contents: testDataNew,
		Name:     unitName,
	}
	unitExpected.Dropins = append(unitExpected.Dropins, dropinThree)
	unitExpected.Dropins = append(unitExpected.Dropins, dropinTwo)
	expectedIgn2Config.Systemd.Units = append(expectedIgn2Config.Systemd.Units, unitExpected)
	expectedMergedUser := types2_4.PasswdUser{
		Name: userName,
		SSHAuthorizedKeys: []types2_4.SSHAuthorizedKey{
			"three",
			"one",
			"two",
		},
	}
	expectedIgn2Config.Passwd.Users = append(expectedIgn2Config.Passwd.Users, userThree, expectedMergedUser)
	assert.Equal(t, expectedIgn2Config, convertedIgn2Config)
}

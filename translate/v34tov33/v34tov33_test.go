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

package v34tov33

import (
	"testing"

	"github.com/coreos/ign-converter/util"
	types3_4 "github.com/coreos/ignition/v2/config/v3_4/types"

	"github.com/stretchr/testify/assert"
)

func TestTranslate3_4to3_3(t *testing.T) {
	emptyConfig := types3_4.Config{
		Ignition: types3_4.Ignition{
			Version: "3.4.0",
		},
	}

	_, err := Translate(emptyConfig)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}

	res, err := Translate(util.NonexhaustiveConfig3_4)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, util.DowntranslateConfig3_3, res)

	_, err = Translate(types3_4.Config{
		Ignition: types3_4.Ignition{
			Version: "3.4.0",
		},
		Storage: types3_4.Storage{
			Luks: []types3_4.Luks{
				{
					Device: util.StrP("/dev/sda"),
					Clevis: types3_4.Clevis{
						Tang: []types3_4.Tang{
							{
								URL:           "http://example.com",
								Thumbprint:    util.StrP("z"),
								Advertisement: util.StrP("{\"payload\": \"eyJrZXlzIjogW3siYWxnIjogIkVTNTEyIiwgImt0eSI6ICJFQyIsICJjcnYiOiAiUC01MjEiLCAieCI6ICJBRGFNajJmazNob21CWTF5WElSQ21uRk92cmUzOFZjdHMwTnNHeDZ6RWNxdEVXcjh5ekhUMkhfa2hjNGpSa19FQWFLdjNrd2RjZ05sOTBLcGhfMGYyQ190IiwgInkiOiAiQUZ2d1UyeGJ5T1RydWo0V1NtcVlqN2wtcUVTZmhWakdCNTI1Q2d6d0NoZUZRRTBvb1o3STYyamt3NkRKQ05yS3VPUDRsSEhicm8tYXhoUk9MSXNJVExvNCIsICJrZXlfb3BzIjogWyJ2ZXJpZnkiXX0sIHsiYWxnIjogIkVDTVIiLCAia3R5IjogIkVDIiwgImNydiI6ICJQLTUyMSIsICJ4IjogIkFOZDVYcTFvZklUbTdNWG16OUY0VVRSYmRNZFNIMl9XNXczTDVWZ0w3b3hwdmpyM0hkLXNLNUVqd3A1V2swMnJMb3NXVUJjYkZyZEhjZFJTTVJoZlVFTFIiLCAieSI6ICJBRVVaVlVZWkFBY2hVcmdoX3poaTV3SUUzeTEycGwzeWhqUk5LcGpSdW9tUFhKaDhRaFhXRmRWZEtMUlEwX1lwUjNOMjNSUk1pU1lvWlg0Qm42QnlrQVBMIiwgImtleV9vcHMiOiBbImRlcml2ZUtleSJdfV19\", \"protected\": \"eyJhbGciOiJFUzUxMiIsImN0eSI6Imp3ay1zZXQranNvbiJ9\", \"signature\": \"APHfSyVzLwELwG0pMJyIP74gWvhHUvDtv0SESBxA2uOdSXq76IdWHW2xvCZDdlNan8pnqUvEedPZjf_vdKBw9MTXAPMkRxVnu64HepKwlrzzm_zG2R4CHpoCOsGgjH9-acYxg-Vha63oMojv3_bV0VHg-NbzNLaxietgYplstvcNIwkv\"}"),
							},
						},
					},
				},
			},
		},
	})
	assert.Error(t, err)

	_, err = Translate(types3_4.Config{
		Ignition: types3_4.Ignition{
			Version: "3.4.0",
		},
		Storage: types3_4.Storage{
			Luks: []types3_4.Luks{
				{
					Device:  util.StrP("/dev/sda"),
					Discard: util.BoolP(true),
				},
			},
		},
	})
	assert.Error(t, err)

	_, err = Translate(types3_4.Config{
		Ignition: types3_4.Ignition{
			Version: "3.4.0",
		},
		Storage: types3_4.Storage{
			Luks: []types3_4.Luks{
				{
					Device:      util.StrP("/dev/sda"),
					OpenOptions: []types3_4.OpenOption{"foo"},
				},
			},
		},
	})
	assert.Error(t, err)

	_, err = Translate(types3_4.Config{
		Ignition: types3_4.Ignition{
			Version: "3.4.0",
		},
		Storage: types3_4.Storage{
			Files: []types3_4.File{
				{
					Node: types3_4.Node{
						Path: "/root/file.txt",
					},

					FileEmbedded1: types3_4.FileEmbedded1{
						Mode: util.IntP(01777),
					},
				},
			},
		},
	})

	assert.Error(t, err)

	_, err = Translate(types3_4.Config{
		Ignition: types3_4.Ignition{
			Version: "3.4.0",
		},
		Storage: types3_4.Storage{
			Directories: []types3_4.Directory{
				{
					Node: types3_4.Node{
						Path:      "/rootdir",
						Overwrite: util.BoolP(true),
						User: types3_4.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_4.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_4.DirectoryEmbedded1{
						Mode: util.IntP(01777),
					},
				},
			},
		},
	})

	assert.Error(t, err)

	_, err = Translate(types3_4.Config{
		Ignition: types3_4.Ignition{
			Version: "3.4.0",
		},
		Storage: types3_4.Storage{
			Files: []types3_4.File{
				{
					Node: types3_4.Node{
						Path: "/path",
					},
					FileEmbedded1: types3_4.FileEmbedded1{
						Contents: types3_4.Resource{
							Source: util.StrP("arn:aws:s3:us-west-1:123456789012:accesspoint/test/object/some/path"),
						},
					},
				},
			},
		},
	})
	assert.Error(t, err)
}

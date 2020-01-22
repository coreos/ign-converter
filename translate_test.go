// Copyright 2019 Red Hat, Inc.
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

package ignconverter

import (
	"testing"

	old2_2 "github.com/coreos/ignition/config/v2_2/types"
	old "github.com/coreos/ignition/config/v2_3/types"
	"github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/stretchr/testify/assert"
)

type input struct {
	cfg   old.Config
	fsMap map[string]string
}

// Config using _all_ the (undeprecated) fields
var (
	aSha512Hash = "sha512-c6100de5624cfb3c109909948ecb8d703bbddcd3725b8bd43dcf2cee6d2f5dc990a757575e0306a8e8eea354bcd7cfac354da911719766225668fe5430477fa8"
	aUUID       = "9d6e42cd-dcef-4177-b4c6-2a0c979e3d82"

	exhaustiveConfig = old.Config{
		Ignition: old.Ignition{
			Version: "2.3.0",
			Config: old.IgnitionConfig{
				Append: []old.ConfigReference{
					{
						Source: "https://example.com",
						Verification: old.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: &old.ConfigReference{
					Source: "https://example.com",
					Verification: old.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: old.Timeouts{
				HTTPResponseHeaders: intP(5),
				HTTPTotal:           intP(10),
			},
			Security: old.Security{
				TLS: old.TLS{
					CertificateAuthorities: []old.CaReference{
						{
							Source: "https://example.com",
							Verification: old.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			// Proxy is unsupported
		},
		Storage: old.Storage{
			Disks: []old.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: true,
					Partitions: []old.Partition{
						{
							Label:              strP("var"),
							Number:             1,
							SizeMiB:            intP(5000),
							StartMiB:           intP(2048),
							TypeGUID:           aUUID,
							GUID:               aUUID,
							WipePartitionEntry: true,
							ShouldExist:        boolP(true),
						},
					},
				},
			},
			Raid: []old.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []old.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  1,
					Options: []old.RaidOption{"foobar"},
				},
			},
			Filesystems: []old.Filesystem{
				{
					Name: "var",
					Mount: &old.Mount{
						Device:         "/dev/disk/by-partlabel/var",
						Format:         "xfs",
						WipeFilesystem: true,
						Label:          strP("var"),
						UUID:           &aUUID,
						Options:        []old.MountOption{"rw"},
					},
				},
			},
			Files: []old.File{
				{
					Node: old.Node{
						Filesystem: "var",
						Path:       "/varfile",
						Overwrite:  boolPStrict(false),
						User: &old.NodeUser{
							ID: intP(1000),
						},
						Group: &old.NodeGroup{
							Name: "groupname",
						},
					},
					FileEmbedded1: old.FileEmbedded1{
						Append: true,
						Mode:   intP(420),
						Contents: old.FileContents{
							Compression: "gzip",
							Source:      "https://example.com",
							Verification: old.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
				{
					Node: old.Node{
						Filesystem: "root",
						Path:       "/empty",
					},
					FileEmbedded1: old.FileEmbedded1{
						Mode: intP(420),
					},
				},
			},
			Directories: []old.Directory{
				{
					Node: old.Node{
						Filesystem: "root",
						Path:       "/rootdir",
						Overwrite:  boolP(true),
						User: &old.NodeUser{
							ID: intP(1000),
						},
						Group: &old.NodeGroup{
							Name: "groupname",
						},
					},
					DirectoryEmbedded1: old.DirectoryEmbedded1{
						Mode: intP(420),
					},
				},
			},
			Links: []old.Link{
				{
					Node: old.Node{
						Filesystem: "root",
						Path:       "/rootlink",
						Overwrite:  boolP(true),
						User: &old.NodeUser{
							ID: intP(1000),
						},
						Group: &old.NodeGroup{
							Name: "groupname",
						},
					},
					LinkEmbedded1: old.LinkEmbedded1{
						Hard:   false,
						Target: "/foobar",
					},
				},
			},
		},
	}
	exhaustiveMap     = map[string]string{"var": "/var"}
	exhaustiveConfig3 = types.Config{
		Ignition: types.Ignition{
			Version: "3.0.0",
			Config: types.IgnitionConfig{
				Merge: []types.ConfigReference{
					{
						Source: strP("https://example.com"),
						Verification: types.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types.ConfigReference{
					Source: strP("https://example.com"),
					Verification: types.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types.Timeouts{
				HTTPResponseHeaders: intP(5),
				HTTPTotal:           intP(10),
			},
			Security: types.Security{
				TLS: types.TLS{
					CertificateAuthorities: []types.CaReference{
						{
							Source: "https://example.com",
							Verification: types.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			// Proxy is unsupported
		},
		Storage: types.Storage{
			Disks: []types.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: boolP(true),
					Partitions: []types.Partition{
						{
							Label:              strP("var"),
							Number:             1,
							SizeMiB:            intP(5000),
							StartMiB:           intP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: boolP(true),
							ShouldExist:        boolP(true),
						},
					},
				},
			},
			Raid: []types.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  intP(1),
					Options: []types.RaidOption{"foobar"},
				},
			},
			Filesystems: []types.Filesystem{
				{
					Path:           strP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         strP("xfs"),
					WipeFilesystem: boolP(true),
					Label:          strP("var"),
					UUID:           &aUUID,
					Options:        []types.FilesystemOption{"rw"},
				},
			},
			Files: []types.File{
				{
					Node: types.Node{
						Path:      "/var/varfile",
						Overwrite: boolPStrict(false),
						User: types.NodeUser{
							ID: intP(1000),
						},
						Group: types.NodeGroup{
							Name: strP("groupname"),
						},
					},
					FileEmbedded1: types.FileEmbedded1{
						Mode: intP(420),
						Append: []types.FileContents{
							{
								Compression: strP("gzip"),
								Source:      strP("https://example.com"),
								Verification: types.Verification{
									Hash: &aSha512Hash,
								},
							},
						},
					},
				},
				{
					Node: types.Node{
						Path:      "/empty",
						Overwrite: boolPStrict(true),
					},
					FileEmbedded1: types.FileEmbedded1{
						Mode: intP(420),
						Contents: types.FileContents{
							Source: strPStrict(""),
						},
					},
				},
			},
			Directories: []types.Directory{
				{
					Node: types.Node{
						Path:      "/rootdir",
						Overwrite: boolP(true),
						User: types.NodeUser{
							ID: intP(1000),
						},
						Group: types.NodeGroup{
							Name: strP("groupname"),
						},
					},
					DirectoryEmbedded1: types.DirectoryEmbedded1{
						Mode: intP(420),
					},
				},
			},
			Links: []types.Link{
				{
					Node: types.Node{
						Path:      "/rootlink",
						Overwrite: boolP(true),
						User: types.NodeUser{
							ID: intP(1000),
						},
						Group: types.NodeGroup{
							Name: strP("groupname"),
						},
					},
					LinkEmbedded1: types.LinkEmbedded1{
						Hard:   boolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	// Variables with slight modifications for downtranslation
	exhaustiveConfigDowntranslate = old2_2.Config{
		Ignition: old2_2.Ignition{
			Version: "2.2.0",
			Config: old2_2.IgnitionConfig{
				Append: []old2_2.ConfigReference{
					{
						Source: "https://example.com",
						Verification: old2_2.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: &old2_2.ConfigReference{
					Source: "https://example.com",
					Verification: old2_2.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: old2_2.Timeouts{
				HTTPResponseHeaders: intP(5),
				HTTPTotal:           intP(10),
			},
			Security: old2_2.Security{
				TLS: old2_2.TLS{
					CertificateAuthorities: []old2_2.CaReference{
						{
							Source: "https://example.com",
							Verification: old2_2.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			// Proxy is unsupported
		},
		Storage: old2_2.Storage{
			Disks: []old2_2.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: true,
					Partitions: []old2_2.Partition{
						{
							Label:    "var",
							Number:   1,
							TypeGUID: aUUID,
							GUID:     aUUID,
						},
					},
				},
			},
			Raid: []old2_2.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []old2_2.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  1,
					Options: []old2_2.RaidOption{"foobar"},
				},
			},
			Filesystems: []old2_2.Filesystem{
				{
					Name: "/var",
					Mount: &old2_2.Mount{
						Device:         "/dev/disk/by-partlabel/var",
						Format:         "xfs",
						WipeFilesystem: true,
						Label:          strP("var"),
						UUID:           &aUUID,
						Options:        []old2_2.MountOption{"rw"},
					},
				},
			},
			Files: []old2_2.File{
				{
					Node: old2_2.Node{
						Filesystem: "/var",
						Path:       "/varfile",
						Overwrite:  boolPStrict(false),
						User: &old2_2.NodeUser{
							ID: intP(1000),
						},
						Group: &old2_2.NodeGroup{
							Name: "groupname",
						},
					},
					FileEmbedded1: old2_2.FileEmbedded1{
						Append: true,
						Mode:   intP(420),
						Contents: old2_2.FileContents{
							Compression: "gzip",
							Source:      "https://example.com",
							Verification: old2_2.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
				{
					Node: old2_2.Node{
						Filesystem: "root",
						Path:       "/empty",
						Overwrite:  boolPStrict(false),
					},
					FileEmbedded1: old2_2.FileEmbedded1{
						Mode: intP(420),
					},
				},
			},
			Directories: []old2_2.Directory{
				{
					Node: old2_2.Node{
						Filesystem: "root",
						Path:       "/rootdir",
						Overwrite:  boolP(true),
						User: &old2_2.NodeUser{
							ID: intP(1000),
						},
						Group: &old2_2.NodeGroup{
							Name: "groupname",
						},
					},
					DirectoryEmbedded1: old2_2.DirectoryEmbedded1{
						Mode: intP(420),
					},
				},
			},
			Links: []old2_2.Link{
				{
					Node: old2_2.Node{
						Filesystem: "root",
						Path:       "/rootlink",
						Overwrite:  boolP(true),
						User: &old2_2.NodeUser{
							ID: intP(1000),
						},
						Group: &old2_2.NodeGroup{
							Name: "groupname",
						},
					},
					LinkEmbedded1: old2_2.LinkEmbedded1{
						Hard:   false,
						Target: "/foobar",
					},
				},
			},
		},
	}

	exhaustiveConfig3Downtranslate = types.Config{
		Ignition: types.Ignition{
			Version: "3.0.0",
			Config: types.IgnitionConfig{
				Merge: []types.ConfigReference{
					{
						Source: strP("https://example.com"),
						Verification: types.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types.ConfigReference{
					Source: strP("https://example.com"),
					Verification: types.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types.Timeouts{
				HTTPResponseHeaders: intP(5),
				HTTPTotal:           intP(10),
			},
			Security: types.Security{
				TLS: types.TLS{
					CertificateAuthorities: []types.CaReference{
						{
							Source: "https://example.com",
							Verification: types.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			// Proxy is unsupported
		},
		Storage: types.Storage{
			Disks: []types.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: boolP(true),
					Partitions: []types.Partition{
						{
							Label:              strP("var"),
							Number:             1,
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: boolP(true),
							ShouldExist:        boolP(true),
						},
					},
				},
			},
			Raid: []types.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  intP(1),
					Options: []types.RaidOption{"foobar"},
				},
			},
			Filesystems: []types.Filesystem{
				{
					Path:           strP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         strP("xfs"),
					WipeFilesystem: boolP(true),
					Label:          strP("var"),
					UUID:           &aUUID,
					Options:        []types.FilesystemOption{"rw"},
				},
			},
			Files: []types.File{
				{
					Node: types.Node{
						Path:      "/var/varfile",
						Overwrite: boolPStrict(false),
						User: types.NodeUser{
							ID: intP(1000),
						},
						Group: types.NodeGroup{
							Name: strP("groupname"),
						},
					},
					FileEmbedded1: types.FileEmbedded1{
						Mode: intP(420),
						Append: []types.FileContents{
							{
								Compression: strP("gzip"),
								Source:      strP("https://example.com"),
								Verification: types.Verification{
									Hash: &aSha512Hash,
								},
							},
						},
					},
				},
				{
					Node: types.Node{
						Path: "/empty",
					},
					FileEmbedded1: types.FileEmbedded1{
						Mode: intP(420),
						Contents: types.FileContents{
							Source: strPStrict(""),
						},
					},
				},
			},
			Directories: []types.Directory{
				{
					Node: types.Node{
						Path:      "/rootdir",
						Overwrite: boolP(true),
						User: types.NodeUser{
							ID: intP(1000),
						},
						Group: types.NodeGroup{
							Name: strP("groupname"),
						},
					},
					DirectoryEmbedded1: types.DirectoryEmbedded1{
						Mode: intP(420),
					},
				},
			},
			Links: []types.Link{
				{
					Node: types.Node{
						Path:      "/rootlink",
						Overwrite: boolP(true),
						User: types.NodeUser{
							ID: intP(1000),
						},
						Group: types.NodeGroup{
							Name: strP("groupname"),
						},
					},
					LinkEmbedded1: types.LinkEmbedded1{
						Hard:   boolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}
)

func TestCheck(t *testing.T) {
	goodConfigs := []input{
		{
			exhaustiveConfig,
			exhaustiveMap,
		},
	}
	badConfigs := []input{
		{}, // empty config has no version, fails validation
		{
			// need a map for filesystems
			exhaustiveConfig,
			nil,
		},
	}
	for i, e := range goodConfigs {
		if err := Check(e.cfg, e.fsMap); err != nil {
			t.Errorf("Good config test %d: got %v, expected nil", i, err)
		}
	}
	for i, e := range badConfigs {
		if err := Check(e.cfg, e.fsMap); err == nil {
			t.Errorf("Bad config test %d: got ok, expected: %v", i, err)
		}
	}
}

func Test2To3(t *testing.T) {
	res, err := Translate(exhaustiveConfig, exhaustiveMap)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, exhaustiveConfig3, res, "Configs Differed")
}

func Test3To2(t *testing.T) {
	res, err := Translate3to2(exhaustiveConfig3Downtranslate)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, exhaustiveConfigDowntranslate, res, "Configs Differed")
}

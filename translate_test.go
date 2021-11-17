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

package ignconverter

import (
	"fmt"
	"testing"

	types2_2 "github.com/coreos/ignition/config/v2_2/types"
	types2_3 "github.com/coreos/ignition/config/v2_3/types"
	types2_4 "github.com/coreos/ignition/config/v2_4/types"
	types3_0 "github.com/coreos/ignition/v2/config/v3_0/types"
	types3_1 "github.com/coreos/ignition/v2/config/v3_1/types"
	types3_2 "github.com/coreos/ignition/v2/config/v3_2/types"
	types3_3 "github.com/coreos/ignition/v2/config/v3_3/types"

	"github.com/stretchr/testify/assert"

	"github.com/coreos/ign-converter/translate/v23tov30"
	"github.com/coreos/ign-converter/translate/v24tov31"
	"github.com/coreos/ign-converter/translate/v30tov22"
	"github.com/coreos/ign-converter/translate/v31tov22"
	"github.com/coreos/ign-converter/translate/v31tov24"
	"github.com/coreos/ign-converter/translate/v32tov22"
	"github.com/coreos/ign-converter/translate/v32tov24"
	"github.com/coreos/ign-converter/translate/v32tov31"
	"github.com/coreos/ign-converter/translate/v33tov32"
	"github.com/coreos/ign-converter/util"
)

// Configs using _all_ the (undeprecated) fields
var (
	aSha512Hash   = "sha512-c6100de5624cfb3c109909948ecb8d703bbddcd3725b8bd43dcf2cee6d2f5dc990a757575e0306a8e8eea354bcd7cfac354da911719766225668fe5430477fa8"
	aUUID         = "9d6e42cd-dcef-4177-b4c6-2a0c979e3d82"
	exhaustiveMap = map[string]string{
		"var":  "/var",
		"/var": "/var",
	}

	exhaustiveConfig2_3 = types2_3.Config{
		Ignition: types2_3.Ignition{
			Version: "2.3.0",
			Config: types2_3.IgnitionConfig{
				Append: []types2_3.ConfigReference{
					{
						Source: "https://example.com",
						Verification: types2_3.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: &types2_3.ConfigReference{
					Source: "https://example.com",
					Verification: types2_3.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types2_3.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types2_3.Security{
				TLS: types2_3.TLS{
					CertificateAuthorities: []types2_3.CaReference{
						{
							Source: "https://example.com",
							Verification: types2_3.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			// Proxy is unsupported
		},
		Storage: types2_3.Storage{
			Disks: []types2_3.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: true,
					Partitions: []types2_3.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							SizeMiB:            util.IntP(5000),
							StartMiB:           util.IntP(2048),
							TypeGUID:           aUUID,
							GUID:               aUUID,
							WipePartitionEntry: true,
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types2_3.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types2_3.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  1,
					Options: []types2_3.RaidOption{"foobar"},
				},
			},
			Filesystems: []types2_3.Filesystem{
				{
					Name: "/var",
					Mount: &types2_3.Mount{
						Device:         "/dev/disk/by-partlabel/var",
						Format:         "xfs",
						WipeFilesystem: true,
						Label:          util.StrP("var"),
						UUID:           &aUUID,
						Options:        []types2_3.MountOption{"rw"},
					},
				},
			},
			Files: []types2_3.File{
				{
					Node: types2_3.Node{
						Filesystem: "/var",
						Path:       "/varfile",
						Overwrite:  util.BoolPStrict(false),
						User: &types2_3.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_3.NodeGroup{
							Name: "groupname",
						},
					},
					FileEmbedded1: types2_3.FileEmbedded1{
						Append: true,
						Mode:   util.IntP(420),
						Contents: types2_3.FileContents{
							Compression: "gzip",
							Source:      "https://example.com",
							Verification: types2_3.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
				{
					Node: types2_3.Node{
						Filesystem: "root",
						Path:       "/empty",
					},
					FileEmbedded1: types2_3.FileEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Directories: []types2_3.Directory{
				{
					Node: types2_3.Node{
						Filesystem: "root",
						Path:       "/rootdir",
						Overwrite:  util.BoolP(true),
						User: &types2_3.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_3.NodeGroup{
							Name: "groupname",
						},
					},
					DirectoryEmbedded1: types2_3.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types2_3.Link{
				{
					Node: types2_3.Node{
						Filesystem: "root",
						Path:       "/rootlink",
						Overwrite:  util.BoolP(true),
						User: &types2_3.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_3.NodeGroup{
							Name: "groupname",
						},
					},
					LinkEmbedded1: types2_3.LinkEmbedded1{
						Hard:   false,
						Target: "/foobar",
					},
				},
			},
		},
	}

	wrongDeprecatedConfig2_4 = types2_4.Config{
		Ignition: types2_4.Ignition{
			Version: "2.4.0",
			Config: types2_4.IgnitionConfig{
				Append: []types2_4.ConfigReference{
					{
						Source: "https://example.com",
						Verification: types2_4.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: &types2_4.ConfigReference{
					Source: "https://example.com",
					Verification: types2_4.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types2_4.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types2_4.Security{
				TLS: types2_4.TLS{
					CertificateAuthorities: []types2_4.CaReference{
						{
							Source: "https://example.com",
							Verification: types2_4.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types2_4.Proxy{
				HTTPProxy:  "https://proxy.example.net/",
				HTTPSProxy: "https://secure.proxy.example.net/",
				NoProxy: []types2_4.NoProxyItem{
					"www.example.net",
					"www.example2.net",
				},
			},
		},
		Storage: types2_4.Storage{
			Disks: []types2_4.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: true,
					Partitions: []types2_4.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							SizeMiB:            util.IntP(5000),
							StartMiB:           util.IntP(2048),
							TypeGUID:           aUUID,
							GUID:               aUUID,
							WipePartitionEntry: true,
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types2_4.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types2_4.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  1,
					Options: []types2_4.RaidOption{"foobar"},
				},
			},
			Filesystems: []types2_4.Filesystem{
				{
					Name: "/var",
					Mount: &types2_4.Mount{
						Device: "/dev/disk/by-partlabel/var",
						Format: "xfs",
						Create: &types2_4.Create{
							Force: true,
							Options: []types2_4.CreateOption{
								"--labl=ROOT",
								types2_4.CreateOption(fmt.Sprintf("--uuid=%s", aUUID)),
							},
						},
						UUID: &aUUID,
					},
				},
			},
			Files: []types2_4.File{
				{
					Node: types2_4.Node{
						Filesystem: "/var",
						Path:       "/varfile",
						Overwrite:  util.BoolPStrict(false),
						User: &types2_4.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					FileEmbedded1: types2_4.FileEmbedded1{
						Append: true,
						Mode:   util.IntP(420),
						Contents: types2_4.FileContents{
							Compression: "gzip",
							Source:      "https://example.com",
							Verification: types2_4.Verification{
								Hash: &aSha512Hash,
							},
							HTTPHeaders: types2_4.HTTPHeaders{
								types2_4.HTTPHeader{
									Name:  "Authorization",
									Value: "Basic YWxhZGRpbjpvcGVuc2VzYW1l",
								},
								types2_4.HTTPHeader{
									Name:  "User-Agent",
									Value: "Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)",
								},
							},
						},
					},
				},
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/empty",
						Overwrite:  util.BoolPStrict(false),
					},
					FileEmbedded1: types2_4.FileEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Directories: []types2_4.Directory{
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/rootdir",
						Overwrite:  util.BoolP(true),
						User: &types2_4.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					DirectoryEmbedded1: types2_4.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types2_4.Link{
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/rootlink",
						Overwrite:  util.BoolP(true),
						User: &types2_4.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					LinkEmbedded1: types2_4.LinkEmbedded1{
						Hard:   false,
						Target: "/foobar",
					},
				},
			},
		},
	}

	deprecatedConfig2_4 = types2_4.Config{
		Ignition: types2_4.Ignition{
			Version: "2.4.0",
			Config: types2_4.IgnitionConfig{
				Append: []types2_4.ConfigReference{
					{
						Source: "https://example.com",
						Verification: types2_4.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: &types2_4.ConfigReference{
					Source: "https://example.com",
					Verification: types2_4.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types2_4.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types2_4.Security{
				TLS: types2_4.TLS{
					CertificateAuthorities: []types2_4.CaReference{
						{
							Source: "https://example.com",
							Verification: types2_4.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types2_4.Proxy{
				HTTPProxy:  "https://proxy.example.net/",
				HTTPSProxy: "https://secure.proxy.example.net/",
				NoProxy: []types2_4.NoProxyItem{
					"www.example.net",
					"www.example2.net",
				},
			},
		},
		Storage: types2_4.Storage{
			Disks: []types2_4.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: true,
					Partitions: []types2_4.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							SizeMiB:            util.IntP(5000),
							StartMiB:           util.IntP(2048),
							TypeGUID:           aUUID,
							GUID:               aUUID,
							WipePartitionEntry: true,
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types2_4.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types2_4.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  1,
					Options: []types2_4.RaidOption{"foobar"},
				},
			},
			Filesystems: []types2_4.Filesystem{
				{
					Name: "/var",
					Mount: &types2_4.Mount{
						Device: "/dev/disk/by-partlabel/var",
						Format: "xfs",
						Label:  util.StrP("var"),
						UUID:   &aUUID,
						Create: &types2_4.Create{
							Force: true,
							Options: []types2_4.CreateOption{
								"--label=var",
								types2_4.CreateOption(fmt.Sprintf("--uuid=%s", aUUID)),
							},
						},
					},
				},
			},
			Files: []types2_4.File{
				{
					Node: types2_4.Node{
						Filesystem: "/var",
						Path:       "/varfile",
						Overwrite:  util.BoolPStrict(false),
						User: &types2_4.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					FileEmbedded1: types2_4.FileEmbedded1{
						Append: true,
						Mode:   util.IntP(420),
						Contents: types2_4.FileContents{
							Compression: "gzip",
							Source:      "https://example.com",
							Verification: types2_4.Verification{
								Hash: &aSha512Hash,
							},
							HTTPHeaders: types2_4.HTTPHeaders{
								types2_4.HTTPHeader{
									Name:  "Authorization",
									Value: "Basic YWxhZGRpbjpvcGVuc2VzYW1l",
								},
								types2_4.HTTPHeader{
									Name:  "User-Agent",
									Value: "Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)",
								},
							},
						},
					},
				},
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/empty",
						Overwrite:  util.BoolPStrict(false),
					},
					FileEmbedded1: types2_4.FileEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Directories: []types2_4.Directory{
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/rootdir",
						Overwrite:  util.BoolP(true),
						User: &types2_4.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					DirectoryEmbedded1: types2_4.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types2_4.Link{
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/rootlink",
						Overwrite:  util.BoolP(true),
						User: &types2_4.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					LinkEmbedded1: types2_4.LinkEmbedded1{
						Hard:   false,
						Target: "/foobar",
					},
				},
			},
		},
	}

	badDeprecatedConfig2_4 = types2_4.Config{
		Ignition: types2_4.Ignition{
			Version: "2.4.0",
			Config: types2_4.IgnitionConfig{
				Append: []types2_4.ConfigReference{
					{
						Source: "https://example.com",
						Verification: types2_4.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: &types2_4.ConfigReference{
					Source: "https://example.com",
					Verification: types2_4.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types2_4.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types2_4.Security{
				TLS: types2_4.TLS{
					CertificateAuthorities: []types2_4.CaReference{
						{
							Source: "https://example.com",
							Verification: types2_4.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types2_4.Proxy{
				HTTPProxy:  "https://proxy.example.net/",
				HTTPSProxy: "https://secure.proxy.example.net/",
				NoProxy: []types2_4.NoProxyItem{
					"www.example.net",
					"www.example2.net",
				},
			},
		},
		Storage: types2_4.Storage{
			Filesystems: []types2_4.Filesystem{
				{
					Name: "/var",
					Mount: &types2_4.Mount{
						Device: "/dev/disk/by-partlabel/var",
						Format: "xfs",
						Label:  util.StrP("var"),
						UUID:   &aUUID,
						Create: &types2_4.Create{
							Force: false,
							Options: []types2_4.CreateOption{
								"--label=var",
								types2_4.CreateOption(fmt.Sprintf("--uuid=%s", aUUID)),
							},
						},
					},
				},
			},
		},
	}

	exhaustiveConfig2_4 = types2_4.Config{
		Ignition: types2_4.Ignition{
			Version: "2.4.0",
			Config: types2_4.IgnitionConfig{
				Append: []types2_4.ConfigReference{
					{
						Source: "https://example.com",
						Verification: types2_4.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: &types2_4.ConfigReference{
					Source: "https://example.com",
					Verification: types2_4.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types2_4.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types2_4.Security{
				TLS: types2_4.TLS{
					CertificateAuthorities: []types2_4.CaReference{
						{
							Source: "https://example.com",
							Verification: types2_4.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types2_4.Proxy{
				HTTPProxy:  "https://proxy.example.net/",
				HTTPSProxy: "https://secure.proxy.example.net/",
				NoProxy: []types2_4.NoProxyItem{
					"www.example.net",
					"www.example2.net",
				},
			},
		},
		Storage: types2_4.Storage{
			Disks: []types2_4.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: true,
					Partitions: []types2_4.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							SizeMiB:            util.IntP(5000),
							StartMiB:           util.IntP(2048),
							TypeGUID:           aUUID,
							GUID:               aUUID,
							WipePartitionEntry: true,
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types2_4.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types2_4.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  1,
					Options: []types2_4.RaidOption{"foobar"},
				},
			},
			Filesystems: []types2_4.Filesystem{
				{
					Name: "/var",
					Mount: &types2_4.Mount{
						Device:         "/dev/disk/by-partlabel/var",
						Format:         "xfs",
						WipeFilesystem: true,
						Label:          util.StrP("var"),
						UUID:           &aUUID,
						Options:        []types2_4.MountOption{"rw"},
					},
				},
			},
			Files: []types2_4.File{
				{
					Node: types2_4.Node{
						Filesystem: "/var",
						Path:       "/varfile",
						Overwrite:  util.BoolPStrict(false),
						User: &types2_4.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					FileEmbedded1: types2_4.FileEmbedded1{
						Append: true,
						Mode:   util.IntP(420),
						Contents: types2_4.FileContents{
							Compression: "gzip",
							Source:      "https://example.com",
							Verification: types2_4.Verification{
								Hash: &aSha512Hash,
							},
							HTTPHeaders: types2_4.HTTPHeaders{
								types2_4.HTTPHeader{
									Name:  "Authorization",
									Value: "Basic YWxhZGRpbjpvcGVuc2VzYW1l",
								},
								types2_4.HTTPHeader{
									Name:  "User-Agent",
									Value: "Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)",
								},
							},
						},
					},
				},
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/empty",
						Overwrite:  util.BoolPStrict(false),
					},
					FileEmbedded1: types2_4.FileEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Directories: []types2_4.Directory{
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/rootdir",
						Overwrite:  util.BoolP(true),
						User: &types2_4.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					DirectoryEmbedded1: types2_4.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types2_4.Link{
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/rootlink",
						Overwrite:  util.BoolP(true),
						User: &types2_4.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					LinkEmbedded1: types2_4.LinkEmbedded1{
						Hard:   false,
						Target: "/foobar",
					},
				},
			},
		},
	}

	exhaustiveConfig3_0 = types3_0.Config{
		Ignition: types3_0.Ignition{
			Version: "3.0.0",
			Config: types3_0.IgnitionConfig{
				Merge: []types3_0.ConfigReference{
					{
						Source: util.StrP("https://example.com"),
						Verification: types3_0.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_0.ConfigReference{
					Source: util.StrP("https://example.com"),
					Verification: types3_0.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_0.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types3_0.Security{
				TLS: types3_0.TLS{
					CertificateAuthorities: []types3_0.CaReference{
						{
							Source: "https://example.com",
							Verification: types3_0.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			// Proxy is unsupported
		},
		Storage: types3_0.Storage{
			Disks: []types3_0.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: util.BoolP(true),
					Partitions: []types3_0.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							SizeMiB:            util.IntP(5000),
							StartMiB:           util.IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: util.BoolP(true),
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types3_0.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_0.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  util.IntP(1),
					Options: []types3_0.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_0.Filesystem{
				{
					Path:           util.StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         util.StrP("xfs"),
					WipeFilesystem: util.BoolP(true),
					Label:          util.StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_0.FilesystemOption{"rw"},
				},
			},
			Files: []types3_0.File{
				{
					Node: types3_0.Node{
						Path:      "/var/varfile",
						Overwrite: util.BoolPStrict(false),
						User: types3_0.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					FileEmbedded1: types3_0.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_0.FileContents{
							{
								Compression: util.StrP("gzip"),
								Source:      util.StrP("https://example.com"),
								Verification: types3_0.Verification{
									Hash: &aSha512Hash,
								},
							},
						},
					},
				},
				{
					Node: types3_0.Node{
						Path:      "/empty",
						Overwrite: util.BoolPStrict(true),
					},
					FileEmbedded1: types3_0.FileEmbedded1{
						Mode: util.IntP(420),
						Contents: types3_0.FileContents{
							Source: util.StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_0.Directory{
				{
					Node: types3_0.Node{
						Path:      "/rootdir",
						Overwrite: util.BoolP(true),
						User: types3_0.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_0.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types3_0.Link{
				{
					Node: types3_0.Node{
						Path:      "/rootlink",
						Overwrite: util.BoolP(true),
						User: types3_0.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_0.LinkEmbedded1{
						Hard:   util.BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	exhaustiveConfig2_2 = types2_2.Config{
		Ignition: types2_2.Ignition{
			Version: "2.2.0",
			Config: types2_2.IgnitionConfig{
				Append: []types2_2.ConfigReference{
					{
						Source: "https://example.com",
						Verification: types2_2.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: &types2_2.ConfigReference{
					Source: "https://example.com",
					Verification: types2_2.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types2_2.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types2_2.Security{
				TLS: types2_2.TLS{
					CertificateAuthorities: []types2_2.CaReference{
						{
							Source: "https://example.com",
							Verification: types2_2.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			// Proxy is unsupported
		},
		Storage: types2_2.Storage{
			Disks: []types2_2.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: true,
					Partitions: []types2_2.Partition{
						{
							Label:    "var",
							Number:   1,
							TypeGUID: aUUID,
							GUID:     aUUID,
						},
					},
				},
			},
			Raid: []types2_2.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types2_2.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  1,
					Options: []types2_2.RaidOption{"foobar"},
				},
			},
			Filesystems: []types2_2.Filesystem{
				{
					Name: "/var",
					Mount: &types2_2.Mount{
						Device:         "/dev/disk/by-partlabel/var",
						Format:         "xfs",
						WipeFilesystem: true,
						Label:          util.StrP("var"),
						UUID:           &aUUID,
						Options:        []types2_2.MountOption{"rw"},
					},
				},
			},
			Files: []types2_2.File{
				{
					Node: types2_2.Node{
						Filesystem: "/var",
						Path:       "/varfile",
						Overwrite:  util.BoolPStrict(false),
						User: &types2_2.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_2.NodeGroup{
							Name: "groupname",
						},
					},
					FileEmbedded1: types2_2.FileEmbedded1{
						Append: true,
						Mode:   util.IntP(420),
						Contents: types2_2.FileContents{
							Compression: "gzip",
							Source:      "https://example.com",
							Verification: types2_2.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
				{
					Node: types2_2.Node{
						Filesystem: "root",
						Path:       "/etc/motd",
						Overwrite:  util.BoolPStrict(false),
					},
					FileEmbedded1: types2_2.FileEmbedded1{
						Append: true,
						Mode:   util.IntP(420),
						Contents: types2_2.FileContents{
							Source: "data:text/plain;base64,Zm9vCg==",
						},
					},
				},
				{
					Node: types2_2.Node{
						Filesystem: "root",
						Path:       "/empty",
						Overwrite:  util.BoolPStrict(false),
					},
					FileEmbedded1: types2_2.FileEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Directories: []types2_2.Directory{
				{
					Node: types2_2.Node{
						Filesystem: "root",
						Path:       "/rootdir",
						Overwrite:  util.BoolP(true),
						User: &types2_2.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_2.NodeGroup{
							Name: "groupname",
						},
					},
					DirectoryEmbedded1: types2_2.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types2_2.Link{
				{
					Node: types2_2.Node{
						Filesystem: "root",
						Path:       "/rootlink",
						Overwrite:  util.BoolP(true),
						User: &types2_2.NodeUser{
							ID: util.IntP(1000),
						},
						Group: &types2_2.NodeGroup{
							Name: "groupname",
						},
					},
					LinkEmbedded1: types2_2.LinkEmbedded1{
						Hard:   false,
						Target: "/foobar",
					},
				},
			},
		},
	}

	downtranslateConfig3_0 = types3_0.Config{
		Ignition: types3_0.Ignition{
			Version: "3.0.0",
			Config: types3_0.IgnitionConfig{
				Merge: []types3_0.ConfigReference{
					{
						Source: util.StrP("https://example.com"),
						Verification: types3_0.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_0.ConfigReference{
					Source: util.StrP("https://example.com"),
					Verification: types3_0.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_0.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types3_0.Security{
				TLS: types3_0.TLS{
					CertificateAuthorities: []types3_0.CaReference{
						{
							Source: "https://example.com",
							Verification: types3_0.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			// Proxy is unsupported
		},
		Storage: types3_0.Storage{
			Disks: []types3_0.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: util.BoolP(true),
					Partitions: []types3_0.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: util.BoolP(true),
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types3_0.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_0.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  util.IntP(1),
					Options: []types3_0.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_0.Filesystem{
				{
					Path:           util.StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         util.StrP("xfs"),
					WipeFilesystem: util.BoolP(true),
					Label:          util.StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_0.FilesystemOption{"rw"},
				},
			},
			Files: []types3_0.File{
				{
					Node: types3_0.Node{
						Path:      "/var/varfile",
						Overwrite: util.BoolPStrict(false),
						User: types3_0.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					FileEmbedded1: types3_0.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_0.FileContents{
							{
								Compression: util.StrP("gzip"),
								Source:      util.StrP("https://example.com"),
								Verification: types3_0.Verification{
									Hash: &aSha512Hash,
								},
							},
						},
					},
				},
				{
					Node: types3_0.Node{
						Path: "/etc/motd",
						// Test default append with overwrite unset
					},
					FileEmbedded1: types3_0.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_0.FileContents{
							{
								Source: util.StrP("data:text/plain;base64,Zm9vCg=="),
							},
						},
					},
				},
				{
					Node: types3_0.Node{
						Path: "/empty",
					},
					FileEmbedded1: types3_0.FileEmbedded1{
						Mode: util.IntP(420),
						Contents: types3_0.FileContents{
							Source: util.StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_0.Directory{
				{
					Node: types3_0.Node{
						Path:      "/rootdir",
						Overwrite: util.BoolP(true),
						User: types3_0.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_0.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types3_0.Link{
				{
					Node: types3_0.Node{
						Path:      "/rootlink",
						Overwrite: util.BoolP(true),
						User: types3_0.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_0.LinkEmbedded1{
						Hard:   util.BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	config3_1WithNoFSOptions = types3_1.Config{
		Ignition: types3_1.Ignition{
			Version: "3.1.0",
			Config: types3_1.IgnitionConfig{
				Merge: []types3_1.Resource{
					{
						Source: util.StrP("https://example.com"),
						Verification: types3_1.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_1.Resource{
					Source: util.StrP("https://example.com"),
					Verification: types3_1.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_1.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types3_1.Security{
				TLS: types3_1.TLS{
					CertificateAuthorities: []types3_1.Resource{
						{
							Source: util.StrP("https://example.com"),
							Verification: types3_1.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types3_1.Proxy{
				HTTPProxy:  util.StrP("https://proxy.example.net/"),
				HTTPSProxy: util.StrP("https://secure.proxy.example.net/"),
				NoProxy: []types3_1.NoProxyItem{
					"www.example.net",
					"www.example2.net",
				},
			},
		},
		Storage: types3_1.Storage{
			Disks: []types3_1.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: util.BoolP(true),
					Partitions: []types3_1.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							SizeMiB:            util.IntP(5000),
							StartMiB:           util.IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: util.BoolP(true),
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types3_1.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_1.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  util.IntP(1),
					Options: []types3_1.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_1.Filesystem{
				{
					Path:           util.StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         util.StrP("xfs"),
					WipeFilesystem: util.BoolP(true),
					Label:          util.StrP("var"),
					UUID:           &aUUID,
					Options: []types3_1.FilesystemOption{
						types3_1.FilesystemOption("--label=var"),
						types3_1.FilesystemOption("--uuid=9d6e42cd-dcef-4177-b4c6-2a0c979e3d82"),
					},
				},
			},
			Files: []types3_1.File{
				{
					Node: types3_1.Node{
						Path:      "/var/varfile",
						Overwrite: util.BoolPStrict(false),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_1.Resource{
							{
								Compression: util.StrP("gzip"),
								Source:      util.StrP("https://example.com"),
								Verification: types3_1.Verification{
									Hash: &aSha512Hash,
								},
								HTTPHeaders: types3_1.HTTPHeaders{
									types3_1.HTTPHeader{
										Name:  "Authorization",
										Value: util.StrP("Basic YWxhZGRpbjpvcGVuc2VzYW1l"),
									},
									types3_1.HTTPHeader{
										Name:  "User-Agent",
										Value: util.StrP("Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
									},
								},
							},
						},
					},
				},
				{
					Node: types3_1.Node{
						Path:      "/empty",
						Overwrite: util.BoolPStrict(false),
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: util.IntP(420),
						Contents: types3_1.Resource{
							Source: util.StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_1.Directory{
				{
					Node: types3_1.Node{
						Path:      "/rootdir",
						Overwrite: util.BoolP(true),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_1.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types3_1.Link{
				{
					Node: types3_1.Node{
						Path:      "/rootlink",
						Overwrite: util.BoolP(true),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_1.LinkEmbedded1{
						Hard:   util.BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	config3_1WithNoFSOptionsAndNoLabel = types3_1.Config{
		Ignition: types3_1.Ignition{
			Version: "3.1.0",
			Config: types3_1.IgnitionConfig{
				Merge: []types3_1.Resource{
					{
						Source: util.StrP("https://example.com"),
						Verification: types3_1.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_1.Resource{
					Source: util.StrP("https://example.com"),
					Verification: types3_1.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_1.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types3_1.Security{
				TLS: types3_1.TLS{
					CertificateAuthorities: []types3_1.Resource{
						{
							Source: util.StrP("https://example.com"),
							Verification: types3_1.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types3_1.Proxy{
				HTTPProxy:  util.StrP("https://proxy.example.net/"),
				HTTPSProxy: util.StrP("https://secure.proxy.example.net/"),
				NoProxy: []types3_1.NoProxyItem{
					"www.example.net",
					"www.example2.net",
				},
			},
		},
		Storage: types3_1.Storage{
			Disks: []types3_1.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: util.BoolP(true),
					Partitions: []types3_1.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							SizeMiB:            util.IntP(5000),
							StartMiB:           util.IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: util.BoolP(true),
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types3_1.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_1.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  util.IntP(1),
					Options: []types3_1.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_1.Filesystem{
				{
					Path:           util.StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         util.StrP("xfs"),
					WipeFilesystem: util.BoolP(true),
					UUID:           &aUUID,
					Options: []types3_1.FilesystemOption{
						types3_1.FilesystemOption("--labl=ROOT"),
						types3_1.FilesystemOption("--uuid=9d6e42cd-dcef-4177-b4c6-2a0c979e3d82"),
					},
				},
			},
			Files: []types3_1.File{
				{
					Node: types3_1.Node{
						Path:      "/var/varfile",
						Overwrite: util.BoolPStrict(false),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_1.Resource{
							{
								Compression: util.StrP("gzip"),
								Source:      util.StrP("https://example.com"),
								Verification: types3_1.Verification{
									Hash: &aSha512Hash,
								},
								HTTPHeaders: types3_1.HTTPHeaders{
									types3_1.HTTPHeader{
										Name:  "Authorization",
										Value: util.StrP("Basic YWxhZGRpbjpvcGVuc2VzYW1l"),
									},
									types3_1.HTTPHeader{
										Name:  "User-Agent",
										Value: util.StrP("Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
									},
								},
							},
						},
					},
				},
				{
					Node: types3_1.Node{
						Path:      "/empty",
						Overwrite: util.BoolPStrict(false),
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: util.IntP(420),
						Contents: types3_1.Resource{
							Source: util.StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_1.Directory{
				{
					Node: types3_1.Node{
						Path:      "/rootdir",
						Overwrite: util.BoolP(true),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_1.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types3_1.Link{
				{
					Node: types3_1.Node{
						Path:      "/rootlink",
						Overwrite: util.BoolP(true),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_1.LinkEmbedded1{
						Hard:   util.BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	nonexhaustiveConfig3_1 = types3_1.Config{
		Ignition: types3_1.Ignition{
			Version: "3.1.0",
			Config: types3_1.IgnitionConfig{
				Merge: []types3_1.Resource{
					{
						Source: util.StrP("https://example.com"),
						Verification: types3_1.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_1.Resource{
					Source: util.StrP("https://example.com"),
					Verification: types3_1.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_1.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types3_1.Security{
				TLS: types3_1.TLS{
					CertificateAuthorities: []types3_1.Resource{
						{
							Source: util.StrP("https://example.com"),
							Verification: types3_1.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types3_1.Proxy{
				HTTPProxy:  util.StrP("https://proxy.example.net/"),
				HTTPSProxy: util.StrP("https://secure.proxy.example.net/"),
				NoProxy: []types3_1.NoProxyItem{
					"www.example.net",
					"www.example2.net",
				},
			},
		},
		Storage: types3_1.Storage{
			Disks: []types3_1.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: util.BoolP(true),
					Partitions: []types3_1.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							SizeMiB:            util.IntP(5000),
							StartMiB:           util.IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: util.BoolP(true),
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types3_1.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_1.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  util.IntP(1),
					Options: []types3_1.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_1.Filesystem{
				{
					Path:           util.StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         util.StrP("xfs"),
					WipeFilesystem: util.BoolP(true),
					Label:          util.StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_1.FilesystemOption{"rw"},
				},
			},
			Files: []types3_1.File{
				{
					Node: types3_1.Node{
						Path:      "/var/varfile",
						Overwrite: util.BoolPStrict(false),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_1.Resource{
							{
								Compression: util.StrP("gzip"),
								Source:      util.StrP("https://example.com"),
								Verification: types3_1.Verification{
									Hash: &aSha512Hash,
								},
								HTTPHeaders: types3_1.HTTPHeaders{
									types3_1.HTTPHeader{
										Name:  "Authorization",
										Value: util.StrP("Basic YWxhZGRpbjpvcGVuc2VzYW1l"),
									},
									types3_1.HTTPHeader{
										Name:  "User-Agent",
										Value: util.StrP("Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
									},
								},
							},
						},
					},
				},
				{
					Node: types3_1.Node{
						Path:      "/empty",
						Overwrite: util.BoolPStrict(false),
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: util.IntP(420),
						Contents: types3_1.Resource{
							Source: util.StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_1.Directory{
				{
					Node: types3_1.Node{
						Path:      "/rootdir",
						Overwrite: util.BoolP(true),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_1.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types3_1.Link{
				{
					Node: types3_1.Node{
						Path:      "/rootlink",
						Overwrite: util.BoolP(true),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_1.LinkEmbedded1{
						Hard:   util.BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	downtranslateConfig3_1 = types3_1.Config{
		Ignition: types3_1.Ignition{
			Version: "3.1.0",
			Config: types3_1.IgnitionConfig{
				Merge: []types3_1.Resource{
					{
						Source: util.StrP("https://example.com"),
						Verification: types3_1.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_1.Resource{
					Source: util.StrP("https://example.com"),
					Verification: types3_1.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_1.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types3_1.Security{
				TLS: types3_1.TLS{
					CertificateAuthorities: []types3_1.Resource{
						{
							Source: util.StrP("https://example.com"),
							Verification: types3_1.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			// Proxy is unsupported
		},
		Storage: types3_1.Storage{
			Disks: []types3_1.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: util.BoolP(true),
					Partitions: []types3_1.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: util.BoolP(true),
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types3_1.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_1.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  util.IntP(1),
					Options: []types3_1.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_1.Filesystem{
				{
					Path:           util.StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         util.StrP("xfs"),
					WipeFilesystem: util.BoolP(true),
					Label:          util.StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_1.FilesystemOption{"rw"},
				},
			},
			Files: []types3_1.File{
				{
					Node: types3_1.Node{
						Path:      "/var/varfile",
						Overwrite: util.BoolPStrict(false),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_1.Resource{
							{
								Compression: util.StrP("gzip"),
								Source:      util.StrP("https://example.com"),
								Verification: types3_1.Verification{
									Hash: &aSha512Hash,
								},
							},
						},
					},
				},
				{
					Node: types3_1.Node{
						Path: "/etc/motd",
						// Test default append with overwrite unset
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_1.Resource{
							{
								Source: util.StrP("data:text/plain;base64,Zm9vCg=="),
							},
						},
					},
				},
				{
					Node: types3_1.Node{
						Path: "/empty",
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: util.IntP(420),
						Contents: types3_1.Resource{
							Source: util.StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_1.Directory{
				{
					Node: types3_1.Node{
						Path:      "/rootdir",
						Overwrite: util.BoolP(true),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_1.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types3_1.Link{
				{
					Node: types3_1.Node{
						Path:      "/rootlink",
						Overwrite: util.BoolP(true),
						User: types3_1.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_1.LinkEmbedded1{
						Hard:   util.BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	nonexhaustiveConfig3_2 = types3_2.Config{
		Ignition: types3_2.Ignition{
			Version: "3.2.0",
			Config: types3_2.IgnitionConfig{
				Merge: []types3_2.Resource{
					{
						Source: util.StrP("https://example.com"),
						Verification: types3_2.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_2.Resource{
					Source: util.StrP("https://example.com"),
					Verification: types3_2.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_2.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types3_2.Security{
				TLS: types3_2.TLS{
					CertificateAuthorities: []types3_2.Resource{
						{
							Source: util.StrP("https://example.com"),
							Verification: types3_2.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types3_2.Proxy{
				HTTPProxy:  util.StrP("https://proxy.example.net/"),
				HTTPSProxy: util.StrP("https://secure.proxy.example.net/"),
				NoProxy: []types3_2.NoProxyItem{
					"www.example.net",
					"www.example2.net",
				},
			},
		},
		Storage: types3_2.Storage{
			Disks: []types3_2.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: util.BoolP(true),
					Partitions: []types3_2.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							SizeMiB:            util.IntP(5000),
							StartMiB:           util.IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: util.BoolP(true),
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types3_2.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_2.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  util.IntP(1),
					Options: []types3_2.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_2.Filesystem{
				{
					Path:           util.StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         util.StrP("xfs"),
					WipeFilesystem: util.BoolP(true),
					Label:          util.StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_2.FilesystemOption{"rw"},
				},
			},
			Files: []types3_2.File{
				{
					Node: types3_2.Node{
						Path:      "/var/varfile",
						Overwrite: util.BoolPStrict(false),
						User: types3_2.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					FileEmbedded1: types3_2.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_2.Resource{
							{
								Compression: util.StrP("gzip"),
								Source:      util.StrP("https://example.com"),
								Verification: types3_2.Verification{
									Hash: &aSha512Hash,
								},
								HTTPHeaders: types3_2.HTTPHeaders{
									types3_2.HTTPHeader{
										Name:  "Authorization",
										Value: util.StrP("Basic YWxhZGRpbjpvcGVuc2VzYW1l"),
									},
									types3_2.HTTPHeader{
										Name:  "User-Agent",
										Value: util.StrP("Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
									},
								},
							},
						},
					},
				},
				{
					Node: types3_2.Node{
						Path:      "/empty",
						Overwrite: util.BoolPStrict(false),
					},
					FileEmbedded1: types3_2.FileEmbedded1{
						Mode: util.IntP(420),
						Contents: types3_2.Resource{
							Source: util.StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_2.Directory{
				{
					Node: types3_2.Node{
						Path:      "/rootdir",
						Overwrite: util.BoolP(true),
						User: types3_2.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_2.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types3_2.Link{
				{
					Node: types3_2.Node{
						Path:      "/rootlink",
						Overwrite: util.BoolP(true),
						User: types3_2.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_2.LinkEmbedded1{
						Hard:   util.BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	nonexhaustiveConfig3_3 = types3_3.Config{
		Ignition: types3_3.Ignition{
			Version: "3.3.0",
			Config: types3_3.IgnitionConfig{
				Merge: []types3_3.Resource{
					{
						Source: util.StrP("https://example.com"),
						Verification: types3_3.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_3.Resource{
					Source: util.StrP("https://example.com"),
					Verification: types3_3.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_3.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types3_3.Security{
				TLS: types3_3.TLS{
					CertificateAuthorities: []types3_3.Resource{
						{
							Source: util.StrP("https://example.com"),
							Verification: types3_3.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types3_3.Proxy{
				HTTPProxy:  util.StrP("https://proxy.example.net/"),
				HTTPSProxy: util.StrP("https://secure.proxy.example.net/"),
				NoProxy: []types3_3.NoProxyItem{
					"www.example.net",
					"www.example2.net",
				},
			},
		},
		Storage: types3_3.Storage{
			Disks: []types3_3.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: util.BoolP(true),
					Partitions: []types3_3.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							SizeMiB:            util.IntP(5000),
							StartMiB:           util.IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: util.BoolP(true),
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types3_3.Raid{
				{
					Name:    "array",
					Level:   util.StrP("raid10"),
					Devices: []types3_3.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  util.IntP(1),
					Options: []types3_3.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_3.Filesystem{
				{
					Path:           util.StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         util.StrP("xfs"),
					WipeFilesystem: util.BoolP(true),
					Label:          util.StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_3.FilesystemOption{"rw"},
				},
			},
			Files: []types3_3.File{
				{
					Node: types3_3.Node{
						Path:      "/var/varfile",
						Overwrite: util.BoolPStrict(false),
						User: types3_3.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_3.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					FileEmbedded1: types3_3.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_3.Resource{
							{
								Compression: util.StrP("gzip"),
								Source:      util.StrP("https://example.com"),
								Verification: types3_3.Verification{
									Hash: &aSha512Hash,
								},
								HTTPHeaders: types3_3.HTTPHeaders{
									types3_3.HTTPHeader{
										Name:  "Authorization",
										Value: util.StrP("Basic YWxhZGRpbjpvcGVuc2VzYW1l"),
									},
									types3_3.HTTPHeader{
										Name:  "User-Agent",
										Value: util.StrP("Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
									},
								},
							},
						},
					},
				},
				{
					Node: types3_3.Node{
						Path:      "/empty",
						Overwrite: util.BoolPStrict(false),
					},
					FileEmbedded1: types3_3.FileEmbedded1{
						Mode: util.IntP(420),
						Contents: types3_3.Resource{
							Source: util.StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_3.Directory{
				{
					Node: types3_3.Node{
						Path:      "/rootdir",
						Overwrite: util.BoolP(true),
						User: types3_3.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_3.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_3.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types3_3.Link{
				{
					Node: types3_3.Node{
						Path:      "/rootlink",
						Overwrite: util.BoolP(true),
						User: types3_3.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_3.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_3.LinkEmbedded1{
						Hard:   util.BoolP(false),
						Target: util.StrP("/foobar"),
					},
				},
			},
		},
	}

	downtranslateConfig3_2 = types3_2.Config{
		Ignition: types3_2.Ignition{
			Version: "3.2.0",
			Config: types3_2.IgnitionConfig{
				Merge: []types3_2.Resource{
					{
						Source: util.StrP("https://example.com"),
						Verification: types3_2.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_2.Resource{
					Source: util.StrP("https://example.com"),
					Verification: types3_2.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_2.Timeouts{
				HTTPResponseHeaders: util.IntP(5),
				HTTPTotal:           util.IntP(10),
			},
			Security: types3_2.Security{
				TLS: types3_2.TLS{
					CertificateAuthorities: []types3_2.Resource{
						{
							Source: util.StrP("https://example.com"),
							Verification: types3_2.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			// Proxy is unsupported
		},
		Storage: types3_2.Storage{
			Disks: []types3_2.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: util.BoolP(true),
					Partitions: []types3_2.Partition{
						{
							Label:              util.StrP("var"),
							Number:             1,
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: util.BoolP(true),
							ShouldExist:        util.BoolP(true),
						},
					},
				},
			},
			Raid: []types3_2.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_2.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  util.IntP(1),
					Options: []types3_2.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_2.Filesystem{
				{
					Path:           util.StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         util.StrP("xfs"),
					WipeFilesystem: util.BoolP(true),
					Label:          util.StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_2.FilesystemOption{"rw"},
				},
			},
			Files: []types3_2.File{
				{
					Node: types3_2.Node{
						Path:      "/var/varfile",
						Overwrite: util.BoolPStrict(false),
						User: types3_2.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					FileEmbedded1: types3_2.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_2.Resource{
							{
								Compression: util.StrP("gzip"),
								Source:      util.StrP("https://example.com"),
								Verification: types3_2.Verification{
									Hash: &aSha512Hash,
								},
							},
						},
					},
				},
				{
					Node: types3_2.Node{
						Path: "/etc/motd",
						// Test default append with overwrite unset
					},
					FileEmbedded1: types3_2.FileEmbedded1{
						Mode: util.IntP(420),
						Append: []types3_2.Resource{
							{
								Source: util.StrP("data:text/plain;base64,Zm9vCg=="),
							},
						},
					},
				},
				{
					Node: types3_2.Node{
						Path: "/empty",
					},
					FileEmbedded1: types3_2.FileEmbedded1{
						Mode: util.IntP(420),
						Contents: types3_2.Resource{
							Source: util.StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_2.Directory{
				{
					Node: types3_2.Node{
						Path:      "/rootdir",
						Overwrite: util.BoolP(true),
						User: types3_2.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_2.DirectoryEmbedded1{
						Mode: util.IntP(420),
					},
				},
			},
			Links: []types3_2.Link{
				{
					Node: types3_2.Node{
						Path:      "/rootlink",
						Overwrite: util.BoolP(true),
						User: types3_2.NodeUser{
							ID: util.IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: util.StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_2.LinkEmbedded1{
						Hard:   util.BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}
)

type input2_3 struct {
	cfg   types2_3.Config
	fsMap map[string]string
}

func TestCheck2_3(t *testing.T) {
	goodConfigs := []input2_3{
		{
			exhaustiveConfig2_3,
			exhaustiveMap,
		},
	}
	badConfigs := []input2_3{
		{}, // empty config has no version, fails validation
		{
			// need a map for filesystems
			exhaustiveConfig2_3,
			nil,
		},
	}
	for i, e := range goodConfigs {
		if err := v23tov30.Check2_3(e.cfg, e.fsMap); err != nil {
			t.Errorf("Good config test %d: got %v, expected nil", i, err)
		}
	}
	for i, e := range badConfigs {
		if err := v23tov30.Check2_3(e.cfg, e.fsMap); err == nil {
			t.Errorf("Bad config test %d: got ok, expected: %v", i, err)
		}
	}
}

type input2_4 struct {
	cfg   types2_4.Config
	fsMap map[string]string
}

func TestCheck2_4(t *testing.T) {
	goodConfigs := []input2_4{
		{
			exhaustiveConfig2_4,
			exhaustiveMap,
		},
	}
	badConfigs := []input2_4{
		{}, // empty config has no version, fails validation
		{
			// need a map for filesystems
			exhaustiveConfig2_4,
			nil,
		},
		{
			// use `mount.create` with `mount.create.force` set to false.
			badDeprecatedConfig2_4,
			exhaustiveMap,
		},
	}
	for i, e := range goodConfigs {
		if err := v24tov31.Check2_4(e.cfg, e.fsMap); err != nil {
			t.Errorf("Good config test %d: got %v, expected nil", i, err)
		}
	}
	for i, e := range badConfigs {
		if err := v24tov31.Check2_4(e.cfg, e.fsMap); err == nil {
			t.Errorf("Bad config test %d: got ok, expected: %v", i, err)
		}
	}
}

func TestTranslate2_3to3_0(t *testing.T) {
	res, err := v23tov30.Translate(exhaustiveConfig2_3, exhaustiveMap)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, exhaustiveConfig3_0, res)
}

func TestTranslate2_4to3_1(t *testing.T) {
	res, err := v24tov31.Translate(exhaustiveConfig2_4, exhaustiveMap)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, nonexhaustiveConfig3_1, res)
}

func TestTranslateDeprecated2_4to3_1(t *testing.T) {
	res, err := v24tov31.Translate(deprecatedConfig2_4, exhaustiveMap)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, config3_1WithNoFSOptions, res)
}

func TestTranslateWrongDeprecated2_4to3_1(t *testing.T) {
	res, err := v24tov31.Translate(wrongDeprecatedConfig2_4, exhaustiveMap)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, config3_1WithNoFSOptionsAndNoLabel, res)
}

func TestTranslate3_0to2_2(t *testing.T) {
	emptyConfig := types3_0.Config{
		Ignition: types3_0.Ignition{
			Version: "3.0.0",
		},
	}

	_, err := v30tov22.Translate(emptyConfig)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}

	res, err := v30tov22.Translate(downtranslateConfig3_0)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, exhaustiveConfig2_2, res)
}

func TestTranslate3_1to2_2(t *testing.T) {
	emptyConfig := types3_1.Config{
		Ignition: types3_1.Ignition{
			Version: "3.1.0",
		},
	}

	_, err := v31tov22.Translate(emptyConfig)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}

	res, err := v31tov22.Translate(downtranslateConfig3_1)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, exhaustiveConfig2_2, res)
}

func TestTranslate3_1to2_4(t *testing.T) {
	emptyConfig := types3_1.Config{
		Ignition: types3_1.Ignition{
			Version: "3.1.0",
		},
	}

	_, err := v31tov24.Translate(emptyConfig)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}

	res, err := v31tov24.Translate(nonexhaustiveConfig3_1)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, exhaustiveConfig2_4, res)
}

func TestTranslate3_2to2_2(t *testing.T) {
	emptyConfig := types3_2.Config{
		Ignition: types3_2.Ignition{
			Version: "3.2.0",
		},
	}

	_, err := v32tov22.Translate(emptyConfig)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}

	res, err := v32tov22.Translate(downtranslateConfig3_2)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, exhaustiveConfig2_2, res)
}

func TestTranslate3_2to2_4(t *testing.T) {
	emptyConfig := types3_2.Config{
		Ignition: types3_2.Ignition{
			Version: "3.2.0",
		},
	}

	_, err := v32tov24.Translate(emptyConfig)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}

	res, err := v32tov24.Translate(nonexhaustiveConfig3_2)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, exhaustiveConfig2_4, res)
}

func TestTranslate3_2to3_1(t *testing.T) {
	emptyConfig := types3_2.Config{
		Ignition: types3_2.Ignition{
			Version: "3.2.0",
		},
	}

	_, err := v32tov31.Translate(emptyConfig)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}

	res, err := v32tov31.Translate(nonexhaustiveConfig3_2)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, nonexhaustiveConfig3_1, res)

	_, err = v32tov31.Translate(types3_2.Config{
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
	_, err = v32tov31.Translate(types3_2.Config{
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
	_, err = v32tov31.Translate(types3_2.Config{
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
	_, err = v32tov31.Translate(types3_2.Config{
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

func TestTranslate3_3to3_2(t *testing.T) {
	emptyConfig := types3_3.Config{
		Ignition: types3_3.Ignition{
			Version: "3.3.0",
		},
	}

	_, err := v33tov32.Translate(emptyConfig)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}

	res, err := v33tov32.Translate(nonexhaustiveConfig3_3)
	if err != nil {
		t.Fatalf("Failed translation: %v", err)
	}
	assert.Equal(t, nonexhaustiveConfig3_2, res)

	_, err = v33tov32.Translate(types3_3.Config{
		Ignition: types3_3.Ignition{
			Version: "3.3.0",
		},
		KernelArguments: types3_3.KernelArguments{
			ShouldExist: []types3_3.KernelArgument{"foo"},
		},
	})
	assert.Error(t, err)
	_, err = v33tov32.Translate(types3_3.Config{
		Ignition: types3_3.Ignition{
			Version: "3.3.0",
		},
		KernelArguments: types3_3.KernelArguments{
			ShouldNotExist: []types3_3.KernelArgument{"foo"},
		},
	})
	assert.Error(t, err)
	_, err = v33tov32.Translate(types3_3.Config{
		Ignition: types3_3.Ignition{
			Version: "3.3.0",
		},
		Storage: types3_3.Storage{
			Filesystems: []types3_3.Filesystem{
				{
					Device: "/dev/foo",
					Format: util.StrP("none"),
				},
			},
		},
	})
	assert.Error(t, err)
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

	convertedIgn2Config, err := v23tov30.RemoveDuplicateFilesUnitsUsers(testIgn2Config)
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

	convertedIgn2Config, err := v24tov31.RemoveDuplicateFilesUnitsUsers(testIgn2Config)
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

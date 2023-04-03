package util

import (
	types2_2 "github.com/coreos/ignition/config/v2_2/types"
	types2_3 "github.com/coreos/ignition/config/v2_3/types"
	types2_4 "github.com/coreos/ignition/config/v2_4/types"
	types3_0 "github.com/coreos/ignition/v2/config/v3_0/types"
	types3_1 "github.com/coreos/ignition/v2/config/v3_1/types"
	types3_2 "github.com/coreos/ignition/v2/config/v3_2/types"
	types3_3 "github.com/coreos/ignition/v2/config/v3_3/types"
	types3_4 "github.com/coreos/ignition/v2/config/v3_4/types"
)

// Configs using _all_ the (undeprecated) fields
var (
	aSha512Hash   = "sha512-c6100de5624cfb3c109909948ecb8d703bbddcd3725b8bd43dcf2cee6d2f5dc990a757575e0306a8e8eea354bcd7cfac354da911719766225668fe5430477fa8"
	aUUID         = "9d6e42cd-dcef-4177-b4c6-2a0c979e3d82"
	ExhaustiveMap = map[string]string{
		"var":  "/var",
		"/var": "/var",
	}

	ExhaustiveConfig2_3 = types2_3.Config{
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
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
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
							Label:              StrP("var"),
							Number:             1,
							SizeMiB:            IntP(5000),
							StartMiB:           IntP(2048),
							TypeGUID:           aUUID,
							GUID:               aUUID,
							WipePartitionEntry: true,
							ShouldExist:        BoolP(true),
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
						Label:          StrP("var"),
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
						Overwrite:  BoolPStrict(false),
						User: &types2_3.NodeUser{
							ID: IntP(1000),
						},
						Group: &types2_3.NodeGroup{
							Name: "groupname",
						},
					},
					FileEmbedded1: types2_3.FileEmbedded1{
						Append: true,
						Mode:   IntP(420),
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
						Mode: IntP(420),
					},
				},
			},
			Directories: []types2_3.Directory{
				{
					Node: types2_3.Node{
						Filesystem: "root",
						Path:       "/rootdir",
						Overwrite:  BoolP(true),
						User: &types2_3.NodeUser{
							ID: IntP(1000),
						},
						Group: &types2_3.NodeGroup{
							Name: "groupname",
						},
					},
					DirectoryEmbedded1: types2_3.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types2_3.Link{
				{
					Node: types2_3.Node{
						Filesystem: "root",
						Path:       "/rootlink",
						Overwrite:  BoolP(true),
						User: &types2_3.NodeUser{
							ID: IntP(1000),
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

	ExhaustiveConfig2_4 = types2_4.Config{
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
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
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
							Label:              StrP("var"),
							Number:             1,
							SizeMiB:            IntP(5000),
							StartMiB:           IntP(2048),
							TypeGUID:           aUUID,
							GUID:               aUUID,
							WipePartitionEntry: true,
							ShouldExist:        BoolP(true),
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
						Label:          StrP("var"),
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
						Overwrite:  BoolPStrict(false),
						User: &types2_4.NodeUser{
							ID: IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					FileEmbedded1: types2_4.FileEmbedded1{
						Append: true,
						Mode:   IntP(420),
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
						Overwrite:  BoolPStrict(false),
					},
					FileEmbedded1: types2_4.FileEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Directories: []types2_4.Directory{
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/rootdir",
						Overwrite:  BoolP(true),
						User: &types2_4.NodeUser{
							ID: IntP(1000),
						},
						Group: &types2_4.NodeGroup{
							Name: "groupname",
						},
					},
					DirectoryEmbedded1: types2_4.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types2_4.Link{
				{
					Node: types2_4.Node{
						Filesystem: "root",
						Path:       "/rootlink",
						Overwrite:  BoolP(true),
						User: &types2_4.NodeUser{
							ID: IntP(1000),
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

	ExhaustiveConfig3_0 = types3_0.Config{
		Ignition: types3_0.Ignition{
			Version: "3.0.0",
			Config: types3_0.IgnitionConfig{
				Merge: []types3_0.ConfigReference{
					{
						Source: StrP("https://example.com"),
						Verification: types3_0.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_0.ConfigReference{
					Source: StrP("https://example.com"),
					Verification: types3_0.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_0.Timeouts{
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
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
					WipeTable: BoolP(true),
					Partitions: []types3_0.Partition{
						{
							Label:              StrP("var"),
							Number:             1,
							SizeMiB:            IntP(5000),
							StartMiB:           IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: BoolP(true),
							ShouldExist:        BoolP(true),
						},
					},
				},
			},
			Raid: []types3_0.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_0.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  IntP(1),
					Options: []types3_0.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_0.Filesystem{
				{
					Path:           StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         StrP("xfs"),
					WipeFilesystem: BoolP(true),
					Label:          StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_0.FilesystemOption{"rw"},
				},
			},
			Files: []types3_0.File{
				{
					Node: types3_0.Node{
						Path:      "/var/varfile",
						Overwrite: BoolPStrict(false),
						User: types3_0.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					FileEmbedded1: types3_0.FileEmbedded1{
						Mode: IntP(420),
						Append: []types3_0.FileContents{
							{
								Compression: StrP("gzip"),
								Source:      StrP("https://example.com"),
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
						Overwrite: BoolPStrict(true),
					},
					FileEmbedded1: types3_0.FileEmbedded1{
						Mode: IntP(420),
						Contents: types3_0.FileContents{
							Source: StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_0.Directory{
				{
					Node: types3_0.Node{
						Path:      "/rootdir",
						Overwrite: BoolP(true),
						User: types3_0.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_0.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types3_0.Link{
				{
					Node: types3_0.Node{
						Path:      "/rootlink",
						Overwrite: BoolP(true),
						User: types3_0.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_0.LinkEmbedded1{
						Hard:   BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	ExhaustiveConfig2_2 = types2_2.Config{
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
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
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
						Label:          StrP("var"),
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
						Overwrite:  BoolPStrict(false),
						User: &types2_2.NodeUser{
							ID: IntP(1000),
						},
						Group: &types2_2.NodeGroup{
							Name: "groupname",
						},
					},
					FileEmbedded1: types2_2.FileEmbedded1{
						Append: true,
						Mode:   IntP(420),
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
						Overwrite:  BoolPStrict(false),
					},
					FileEmbedded1: types2_2.FileEmbedded1{
						Append: true,
						Mode:   IntP(420),
						Contents: types2_2.FileContents{
							Source: "data:text/plain;base64,Zm9vCg==",
						},
					},
				},
				{
					Node: types2_2.Node{
						Filesystem: "root",
						Path:       "/empty",
						Overwrite:  BoolPStrict(false),
					},
					FileEmbedded1: types2_2.FileEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Directories: []types2_2.Directory{
				{
					Node: types2_2.Node{
						Filesystem: "root",
						Path:       "/rootdir",
						Overwrite:  BoolP(true),
						User: &types2_2.NodeUser{
							ID: IntP(1000),
						},
						Group: &types2_2.NodeGroup{
							Name: "groupname",
						},
					},
					DirectoryEmbedded1: types2_2.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types2_2.Link{
				{
					Node: types2_2.Node{
						Filesystem: "root",
						Path:       "/rootlink",
						Overwrite:  BoolP(true),
						User: &types2_2.NodeUser{
							ID: IntP(1000),
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

	DowntranslateConfig3_0 = types3_0.Config{
		Ignition: types3_0.Ignition{
			Version: "3.0.0",
			Config: types3_0.IgnitionConfig{
				Merge: []types3_0.ConfigReference{
					{
						Source: StrP("https://example.com"),
						Verification: types3_0.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_0.ConfigReference{
					Source: StrP("https://example.com"),
					Verification: types3_0.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_0.Timeouts{
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
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
					WipeTable: BoolP(true),
					Partitions: []types3_0.Partition{
						{
							Label:              StrP("var"),
							Number:             1,
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: BoolP(true),
							ShouldExist:        BoolP(true),
						},
					},
				},
			},
			Raid: []types3_0.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_0.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  IntP(1),
					Options: []types3_0.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_0.Filesystem{
				{
					Path:           StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         StrP("xfs"),
					WipeFilesystem: BoolP(true),
					Label:          StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_0.FilesystemOption{"rw"},
				},
			},
			Files: []types3_0.File{
				{
					Node: types3_0.Node{
						Path:      "/var/varfile",
						Overwrite: BoolPStrict(false),
						User: types3_0.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					FileEmbedded1: types3_0.FileEmbedded1{
						Mode: IntP(420),
						Append: []types3_0.FileContents{
							{
								Compression: StrP("gzip"),
								Source:      StrP("https://example.com"),
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
						Mode: IntP(420),
						Append: []types3_0.FileContents{
							{
								Source: StrP("data:text/plain;base64,Zm9vCg=="),
							},
						},
					},
				},
				{
					Node: types3_0.Node{
						Path: "/empty",
					},
					FileEmbedded1: types3_0.FileEmbedded1{
						Mode: IntP(420),
						Contents: types3_0.FileContents{
							Source: StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_0.Directory{
				{
					Node: types3_0.Node{
						Path:      "/rootdir",
						Overwrite: BoolP(true),
						User: types3_0.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_0.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types3_0.Link{
				{
					Node: types3_0.Node{
						Path:      "/rootlink",
						Overwrite: BoolP(true),
						User: types3_0.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_0.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_0.LinkEmbedded1{
						Hard:   BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	NonexhaustiveConfig3_1 = types3_1.Config{
		Ignition: types3_1.Ignition{
			Version: "3.1.0",
			Config: types3_1.IgnitionConfig{
				Merge: []types3_1.Resource{
					{
						Source: StrP("https://example.com"),
						Verification: types3_1.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_1.Resource{
					Source: StrP("https://example.com"),
					Verification: types3_1.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_1.Timeouts{
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
			},
			Security: types3_1.Security{
				TLS: types3_1.TLS{
					CertificateAuthorities: []types3_1.Resource{
						{
							Source: StrP("https://example.com"),
							Verification: types3_1.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types3_1.Proxy{
				HTTPProxy:  StrP("https://proxy.example.net/"),
				HTTPSProxy: StrP("https://secure.proxy.example.net/"),
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
					WipeTable: BoolP(true),
					Partitions: []types3_1.Partition{
						{
							Label:              StrP("var"),
							Number:             1,
							SizeMiB:            IntP(5000),
							StartMiB:           IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: BoolP(true),
							ShouldExist:        BoolP(true),
						},
					},
				},
			},
			Raid: []types3_1.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_1.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  IntP(1),
					Options: []types3_1.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_1.Filesystem{
				{
					Path:           StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         StrP("xfs"),
					WipeFilesystem: BoolP(true),
					Label:          StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_1.FilesystemOption{"rw"},
				},
			},
			Files: []types3_1.File{
				{
					Node: types3_1.Node{
						Path:      "/var/varfile",
						Overwrite: BoolPStrict(false),
						User: types3_1.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: IntP(420),
						Append: []types3_1.Resource{
							{
								Compression: StrP("gzip"),
								Source:      StrP("https://example.com"),
								Verification: types3_1.Verification{
									Hash: &aSha512Hash,
								},
								HTTPHeaders: types3_1.HTTPHeaders{
									types3_1.HTTPHeader{
										Name:  "Authorization",
										Value: StrP("Basic YWxhZGRpbjpvcGVuc2VzYW1l"),
									},
									types3_1.HTTPHeader{
										Name:  "User-Agent",
										Value: StrP("Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
									},
								},
							},
						},
					},
				},
				{
					Node: types3_1.Node{
						Path:      "/empty",
						Overwrite: BoolPStrict(false),
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: IntP(420),
						Contents: types3_1.Resource{
							Source: StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_1.Directory{
				{
					Node: types3_1.Node{
						Path:      "/rootdir",
						Overwrite: BoolP(true),
						User: types3_1.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_1.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types3_1.Link{
				{
					Node: types3_1.Node{
						Path:      "/rootlink",
						Overwrite: BoolP(true),
						User: types3_1.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_1.LinkEmbedded1{
						Hard:   BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	DowntranslateConfig3_1 = types3_1.Config{
		Ignition: types3_1.Ignition{
			Version: "3.1.0",
			Config: types3_1.IgnitionConfig{
				Merge: []types3_1.Resource{
					{
						Source: StrP("https://example.com"),
						Verification: types3_1.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_1.Resource{
					Source: StrP("https://example.com"),
					Verification: types3_1.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_1.Timeouts{
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
			},
			Security: types3_1.Security{
				TLS: types3_1.TLS{
					CertificateAuthorities: []types3_1.Resource{
						{
							Source: StrP("https://example.com"),
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
					WipeTable: BoolP(true),
					Partitions: []types3_1.Partition{
						{
							Label:              StrP("var"),
							Number:             1,
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: BoolP(true),
							ShouldExist:        BoolP(true),
						},
					},
				},
			},
			Raid: []types3_1.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_1.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  IntP(1),
					Options: []types3_1.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_1.Filesystem{
				{
					Path:           StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         StrP("xfs"),
					WipeFilesystem: BoolP(true),
					Label:          StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_1.FilesystemOption{"rw"},
				},
			},
			Files: []types3_1.File{
				{
					Node: types3_1.Node{
						Path:      "/var/varfile",
						Overwrite: BoolPStrict(false),
						User: types3_1.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: IntP(420),
						Append: []types3_1.Resource{
							{
								Compression: StrP("gzip"),
								Source:      StrP("https://example.com"),
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
						Mode: IntP(420),
						Append: []types3_1.Resource{
							{
								Source: StrP("data:text/plain;base64,Zm9vCg=="),
							},
						},
					},
				},
				{
					Node: types3_1.Node{
						Path: "/empty",
					},
					FileEmbedded1: types3_1.FileEmbedded1{
						Mode: IntP(420),
						Contents: types3_1.Resource{
							Source: StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_1.Directory{
				{
					Node: types3_1.Node{
						Path:      "/rootdir",
						Overwrite: BoolP(true),
						User: types3_1.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_1.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types3_1.Link{
				{
					Node: types3_1.Node{
						Path:      "/rootlink",
						Overwrite: BoolP(true),
						User: types3_1.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_1.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_1.LinkEmbedded1{
						Hard:   BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	NonexhaustiveConfig3_2 = types3_2.Config{
		Ignition: types3_2.Ignition{
			Version: "3.2.0",
			Config: types3_2.IgnitionConfig{
				Merge: []types3_2.Resource{
					{
						Source: StrP("https://example.com"),
						Verification: types3_2.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_2.Resource{
					Source: StrP("https://example.com"),
					Verification: types3_2.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_2.Timeouts{
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
			},
			Security: types3_2.Security{
				TLS: types3_2.TLS{
					CertificateAuthorities: []types3_2.Resource{
						{
							Source: StrP("https://example.com"),
							Verification: types3_2.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types3_2.Proxy{
				HTTPProxy:  StrP("https://proxy.example.net/"),
				HTTPSProxy: StrP("https://secure.proxy.example.net/"),
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
					WipeTable: BoolP(true),
					Partitions: []types3_2.Partition{
						{
							Label:              StrP("var"),
							Number:             1,
							SizeMiB:            IntP(5000),
							StartMiB:           IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: BoolP(true),
							ShouldExist:        BoolP(true),
						},
					},
				},
			},
			Raid: []types3_2.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_2.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  IntP(1),
					Options: []types3_2.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_2.Filesystem{
				{
					Path:           StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         StrP("xfs"),
					WipeFilesystem: BoolP(true),
					Label:          StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_2.FilesystemOption{"rw"},
				},
			},
			Files: []types3_2.File{
				{
					Node: types3_2.Node{
						Path:      "/var/varfile",
						Overwrite: BoolPStrict(false),
						User: types3_2.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					FileEmbedded1: types3_2.FileEmbedded1{
						Mode: IntP(420),
						Append: []types3_2.Resource{
							{
								Compression: StrP("gzip"),
								Source:      StrP("https://example.com"),
								Verification: types3_2.Verification{
									Hash: &aSha512Hash,
								},
								HTTPHeaders: types3_2.HTTPHeaders{
									types3_2.HTTPHeader{
										Name:  "Authorization",
										Value: StrP("Basic YWxhZGRpbjpvcGVuc2VzYW1l"),
									},
									types3_2.HTTPHeader{
										Name:  "User-Agent",
										Value: StrP("Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
									},
								},
							},
						},
					},
				},
				{
					Node: types3_2.Node{
						Path:      "/empty",
						Overwrite: BoolPStrict(false),
					},
					FileEmbedded1: types3_2.FileEmbedded1{
						Mode: IntP(420),
						Contents: types3_2.Resource{
							Source: StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_2.Directory{
				{
					Node: types3_2.Node{
						Path:      "/rootdir",
						Overwrite: BoolP(true),
						User: types3_2.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_2.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types3_2.Link{
				{
					Node: types3_2.Node{
						Path:      "/rootlink",
						Overwrite: BoolP(true),
						User: types3_2.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_2.LinkEmbedded1{
						Hard:   BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	NonexhaustiveConfig3_3 = types3_3.Config{
		Ignition: types3_3.Ignition{
			Version: "3.3.0",
			Config: types3_3.IgnitionConfig{
				Merge: []types3_3.Resource{
					{
						Source: StrP("https://example.com"),
						Verification: types3_3.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_3.Resource{
					Source: StrP("https://example.com"),
					Verification: types3_3.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_3.Timeouts{
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
			},
			Security: types3_3.Security{
				TLS: types3_3.TLS{
					CertificateAuthorities: []types3_3.Resource{
						{
							Source: StrP("https://example.com"),
							Verification: types3_3.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types3_3.Proxy{
				HTTPProxy:  StrP("https://proxy.example.net/"),
				HTTPSProxy: StrP("https://secure.proxy.example.net/"),
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
					WipeTable: BoolP(true),
					Partitions: []types3_3.Partition{
						{
							Label:              StrP("var"),
							Number:             1,
							SizeMiB:            IntP(5000),
							StartMiB:           IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: BoolP(true),
							ShouldExist:        BoolP(true),
						},
					},
				},
			},
			Raid: []types3_3.Raid{
				{
					Name:    "array",
					Level:   StrP("raid10"),
					Devices: []types3_3.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  IntP(1),
					Options: []types3_3.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_3.Filesystem{
				{
					Path:           StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         StrP("xfs"),
					WipeFilesystem: BoolP(true),
					Label:          StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_3.FilesystemOption{"rw"},
				},
			},
			Files: []types3_3.File{
				{
					Node: types3_3.Node{
						Path:      "/var/varfile",
						Overwrite: BoolPStrict(false),
						User: types3_3.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_3.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					FileEmbedded1: types3_3.FileEmbedded1{
						Mode: IntP(420),
						Append: []types3_3.Resource{
							{
								Compression: StrP("gzip"),
								Source:      StrP("https://example.com"),
								Verification: types3_3.Verification{
									Hash: &aSha512Hash,
								},
								HTTPHeaders: types3_3.HTTPHeaders{
									types3_3.HTTPHeader{
										Name:  "Authorization",
										Value: StrP("Basic YWxhZGRpbjpvcGVuc2VzYW1l"),
									},
									types3_3.HTTPHeader{
										Name:  "User-Agent",
										Value: StrP("Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
									},
								},
							},
						},
					},
				},
				{
					Node: types3_3.Node{
						Path:      "/empty",
						Overwrite: BoolPStrict(false),
					},
					FileEmbedded1: types3_3.FileEmbedded1{
						Mode: IntP(420),
						Contents: types3_3.Resource{
							Source: StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_3.Directory{
				{
					Node: types3_3.Node{
						Path:      "/rootdir",
						Overwrite: BoolP(true),
						User: types3_3.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_3.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_3.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types3_3.Link{
				{
					Node: types3_3.Node{
						Path:      "/rootlink",
						Overwrite: BoolP(true),
						User: types3_3.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_3.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_3.LinkEmbedded1{
						Hard:   BoolP(false),
						Target: StrP("/foobar"),
					},
				},
			},
		},
	}

	DowntranslateConfig3_2 = types3_2.Config{
		Ignition: types3_2.Ignition{
			Version: "3.2.0",
			Config: types3_2.IgnitionConfig{
				Merge: []types3_2.Resource{
					{
						Source: StrP("https://example.com"),
						Verification: types3_2.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_2.Resource{
					Source: StrP("https://example.com"),
					Verification: types3_2.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_2.Timeouts{
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
			},
			Security: types3_2.Security{
				TLS: types3_2.TLS{
					CertificateAuthorities: []types3_2.Resource{
						{
							Source: StrP("https://example.com"),
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
					WipeTable: BoolP(true),
					Partitions: []types3_2.Partition{
						{
							Label:              StrP("var"),
							Number:             1,
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: BoolP(true),
							ShouldExist:        BoolP(true),
						},
					},
				},
			},
			Raid: []types3_2.Raid{
				{
					Name:    "array",
					Level:   "raid10",
					Devices: []types3_2.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  IntP(1),
					Options: []types3_2.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_2.Filesystem{
				{
					Path:           StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         StrP("xfs"),
					WipeFilesystem: BoolP(true),
					Label:          StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_2.FilesystemOption{"rw"},
				},
			},
			Files: []types3_2.File{
				{
					Node: types3_2.Node{
						Path:      "/var/varfile",
						Overwrite: BoolPStrict(false),
						User: types3_2.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					FileEmbedded1: types3_2.FileEmbedded1{
						Mode: IntP(420),
						Append: []types3_2.Resource{
							{
								Compression: StrP("gzip"),
								Source:      StrP("https://example.com"),
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
						Mode: IntP(420),
						Append: []types3_2.Resource{
							{
								Source: StrP("data:text/plain;base64,Zm9vCg=="),
							},
						},
					},
				},
				{
					Node: types3_2.Node{
						Path: "/empty",
					},
					FileEmbedded1: types3_2.FileEmbedded1{
						Mode: IntP(420),
						Contents: types3_2.Resource{
							Source: StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_2.Directory{
				{
					Node: types3_2.Node{
						Path:      "/rootdir",
						Overwrite: BoolP(true),
						User: types3_2.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_2.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types3_2.Link{
				{
					Node: types3_2.Node{
						Path:      "/rootlink",
						Overwrite: BoolP(true),
						User: types3_2.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_2.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_2.LinkEmbedded1{
						Hard:   BoolP(false),
						Target: "/foobar",
					},
				},
			},
		},
	}

	NonexhaustiveConfig3_4 = types3_4.Config{
		Ignition: types3_4.Ignition{
			Version: "3.4.0",
			Config: types3_4.IgnitionConfig{
				Merge: []types3_4.Resource{
					{
						Source: StrP("https://example.com"),
						Verification: types3_4.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_4.Resource{
					Source: StrP("https://example.com"),
					Verification: types3_4.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_4.Timeouts{
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
			},
			Security: types3_4.Security{
				TLS: types3_4.TLS{
					CertificateAuthorities: []types3_4.Resource{
						{
							Source: StrP("https://example.com"),
							Verification: types3_4.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types3_4.Proxy{
				HTTPProxy:  StrP("https://proxy.example.net/"),
				HTTPSProxy: StrP("https://secure.proxy.example.net/"),
				NoProxy: []types3_4.NoProxyItem{
					"www.example.net",
					"www.example2.net",
				},
			},
		},
		Storage: types3_4.Storage{
			Disks: []types3_4.Disk{
				{
					Device:    "/dev/sda",
					WipeTable: BoolP(true),
					Partitions: []types3_4.Partition{
						{
							Label:              StrP("var"),
							Number:             1,
							SizeMiB:            IntP(5000),
							StartMiB:           IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: BoolP(true),
							ShouldExist:        BoolP(true),
						},
					},
				},
			},
			Raid: []types3_4.Raid{
				{
					Name:    "array",
					Level:   StrP("raid10"),
					Devices: []types3_4.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  IntP(1),
					Options: []types3_4.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_4.Filesystem{
				{
					Path:           StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         StrP("xfs"),
					WipeFilesystem: BoolP(true),
					Label:          StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_4.FilesystemOption{"rw"},
				},
			},
			Files: []types3_4.File{
				{
					Node: types3_4.Node{
						Path:      "/var/varfile",
						Overwrite: BoolPStrict(false),
						User: types3_4.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_4.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					FileEmbedded1: types3_4.FileEmbedded1{
						Mode: IntP(420),
						Append: []types3_4.Resource{
							{
								Compression: StrP("gzip"),
								Source:      StrP("https://example.com"),
								Verification: types3_4.Verification{
									Hash: &aSha512Hash,
								},
								HTTPHeaders: types3_4.HTTPHeaders{
									types3_4.HTTPHeader{
										Name:  "Authorization",
										Value: StrP("Basic YWxhZGRpbjpvcGVuc2VzYW1l"),
									},
									types3_4.HTTPHeader{
										Name:  "User-Agent",
										Value: StrP("Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
									},
								},
							},
						},
					},
				},
				{
					Node: types3_4.Node{
						Path:      "/empty",
						Overwrite: BoolPStrict(false),
					},
					FileEmbedded1: types3_4.FileEmbedded1{
						Mode: IntP(420),
						Contents: types3_4.Resource{
							Source: StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_4.Directory{
				{
					Node: types3_4.Node{
						Path:      "/rootdir",
						Overwrite: BoolP(true),
						User: types3_4.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_4.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_4.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types3_4.Link{
				{
					Node: types3_4.Node{
						Path:      "/rootlink",
						Overwrite: BoolP(true),
						User: types3_4.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_4.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_4.LinkEmbedded1{
						Hard:   BoolP(false),
						Target: StrP("/foobar"),
					},
				},
			},
		},
	}

	DowntranslateConfig3_3 = types3_3.Config{
		Ignition: types3_3.Ignition{
			Version: "3.3.0",
			Config: types3_3.IgnitionConfig{
				Merge: []types3_3.Resource{
					{
						Source: StrP("https://example.com"),
						Verification: types3_3.Verification{
							Hash: &aSha512Hash,
						},
					},
				},
				Replace: types3_3.Resource{
					Source: StrP("https://example.com"),
					Verification: types3_3.Verification{
						Hash: &aSha512Hash,
					},
				},
			},
			Timeouts: types3_3.Timeouts{
				HTTPResponseHeaders: IntP(5),
				HTTPTotal:           IntP(10),
			},
			Security: types3_3.Security{
				TLS: types3_3.TLS{
					CertificateAuthorities: []types3_3.Resource{
						{
							Source: StrP("https://example.com"),
							Verification: types3_3.Verification{
								Hash: &aSha512Hash,
							},
						},
					},
				},
			},
			Proxy: types3_3.Proxy{
				HTTPProxy:  StrP("https://proxy.example.net/"),
				HTTPSProxy: StrP("https://secure.proxy.example.net/"),
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
					WipeTable: BoolP(true),
					Partitions: []types3_3.Partition{
						{
							Label:              StrP("var"),
							Number:             1,
							SizeMiB:            IntP(5000),
							StartMiB:           IntP(2048),
							TypeGUID:           &aUUID,
							GUID:               &aUUID,
							WipePartitionEntry: BoolP(true),
							ShouldExist:        BoolP(true),
						},
					},
				},
			},
			Raid: []types3_3.Raid{
				{
					Name:    "array",
					Level:   StrP("raid10"),
					Devices: []types3_3.Device{"/dev/sdb", "/dev/sdc"},
					Spares:  IntP(1),
					Options: []types3_3.RaidOption{"foobar"},
				},
			},
			Filesystems: []types3_3.Filesystem{
				{
					Path:           StrP("/var"),
					Device:         "/dev/disk/by-partlabel/var",
					Format:         StrP("xfs"),
					WipeFilesystem: BoolP(true),
					Label:          StrP("var"),
					UUID:           &aUUID,
					Options:        []types3_3.FilesystemOption{"rw"},
				},
			},
			Files: []types3_3.File{
				{
					Node: types3_3.Node{
						Path:      "/var/varfile",
						Overwrite: BoolPStrict(false),
						User: types3_3.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_3.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					FileEmbedded1: types3_3.FileEmbedded1{
						Mode: IntP(420),
						Append: []types3_3.Resource{
							{
								Compression: StrP("gzip"),
								Source:      StrP("https://example.com"),
								Verification: types3_3.Verification{
									Hash: &aSha512Hash,
								},
								HTTPHeaders: types3_3.HTTPHeaders{
									types3_3.HTTPHeader{
										Name:  "Authorization",
										Value: StrP("Basic YWxhZGRpbjpvcGVuc2VzYW1l"),
									},
									types3_3.HTTPHeader{
										Name:  "User-Agent",
										Value: StrP("Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
									},
								},
							},
						},
					},
				},
				{
					Node: types3_3.Node{
						Path:      "/empty",
						Overwrite: BoolPStrict(false),
					},
					FileEmbedded1: types3_3.FileEmbedded1{
						Mode: IntP(420),
						Contents: types3_3.Resource{
							Source: StrPStrict(""),
						},
					},
				},
			},
			Directories: []types3_3.Directory{
				{
					Node: types3_3.Node{
						Path:      "/rootdir",
						Overwrite: BoolP(true),
						User: types3_3.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_3.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					DirectoryEmbedded1: types3_3.DirectoryEmbedded1{
						Mode: IntP(420),
					},
				},
			},
			Links: []types3_3.Link{
				{
					Node: types3_3.Node{
						Path:      "/rootlink",
						Overwrite: BoolP(true),
						User: types3_3.NodeUser{
							ID: IntP(1000),
						},
						Group: types3_3.NodeGroup{
							Name: StrP("groupname"),
						},
					},
					LinkEmbedded1: types3_3.LinkEmbedded1{
						Hard:   BoolP(false),
						Target: StrP("/foobar"),
					},
				},
			},
		},
	}
)

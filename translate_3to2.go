package ign2to3

import (
	"fmt"
	"strings"

	old "github.com/coreos/ignition/config/v2_2/types"
	new "github.com/coreos/ignition/v2/config/v3_0/types"
	newValidate "github.com/coreos/ignition/v2/config/validate"
)

func Translate3to2(cfg new.Config) (old.Config, error) {
	rpt := newValidate.ValidateWithContext(cfg, nil)
	if rpt.IsFatal() {
		return old.Config{}, fmt.Errorf("Invalid input config:\n%s", rpt.String())
	}
	// fsMap is a mapping of filesystems populated via the v3 config, to be
	// used for v2 files sections. The naming of each section will be uniquely
	// named by the path
	fsList := generateFsList(cfg.Storage.Filesystems)

	res := old.Config{
		// Ignition section
		Ignition: old.Ignition{
			Version: "2.2.0",
			Config: old.IgnitionConfig{
				Replace: translateCfgRef3to2(cfg.Ignition.Config.Replace),
				Append:  translateCfgRefs3to2(cfg.Ignition.Config.Merge),
			},
			Security: old.Security{
				TLS: old.TLS{
					CertificateAuthorities: translateCAs3to2(cfg.Ignition.Security.TLS.CertificateAuthorities),
				},
			},
			Timeouts: old.Timeouts{
				HTTPResponseHeaders: cfg.Ignition.Timeouts.HTTPResponseHeaders,
				HTTPTotal:           cfg.Ignition.Timeouts.HTTPTotal,
			},
		},
		// Passwd section
		Passwd: old.Passwd{
			Users:  translateUsers3to2(cfg.Passwd.Users),
			Groups: translateGroups3to2(cfg.Passwd.Groups),
		},
		Systemd: old.Systemd{
			Units: translateUnits3to2(cfg.Systemd.Units),
		},
		Storage: old.Storage{
			Disks:       translateDisks3to2(cfg.Storage.Disks),
			Raid:        translateRaid3to2(cfg.Storage.Raid),
			Filesystems: translateFilesystems3to2(cfg.Storage.Filesystems),
			Files:       translateFiles3to2(cfg.Storage.Files, fsList),
			Directories: translateDirectories3to2(cfg.Storage.Directories, fsList),
			Links:       translateLinks3to2(cfg.Storage.Links, fsList),
		},
	}
	return res, nil
}

func generateFsList(fss []new.Filesystem) (ret []string) {
	for _, f := range fss {
		ret = append(ret, *f.Path)
	}
	return
}

func translateCfgRef3to2(ref new.ConfigReference) (ret *old.ConfigReference) {
	if ref.Source == nil {
		return
	}
	ret = &old.ConfigReference{}
	ret.Source = strV(ref.Source)
	ret.Verification.Hash = ref.Verification.Hash
	return
}

func translateCfgRefs3to2(refs []new.ConfigReference) (ret []old.ConfigReference) {
	for _, ref := range refs {
		ret = append(ret, *translateCfgRef3to2(ref))
	}
	return
}

func translateCAs3to2(refs []new.CaReference) (ret []old.CaReference) {
	for _, ref := range refs {
		ret = append(ret, old.CaReference{
			Source: ref.Source,
			Verification: old.Verification{
				Hash: ref.Verification.Hash,
			},
		})
	}
	return
}

func translateUsers3to2(users []new.PasswdUser) (ret []old.PasswdUser) {
	for _, u := range users {
		ret = append(ret, old.PasswdUser{
			Name:              u.Name,
			PasswordHash:      u.PasswordHash,
			SSHAuthorizedKeys: translateUserSSH3to2(u.SSHAuthorizedKeys),
			UID:               u.UID,
			Gecos:             strV(u.Gecos),
			HomeDir:           strV(u.HomeDir),
			NoCreateHome:      boolV(u.NoCreateHome),
			PrimaryGroup:      strV(u.PrimaryGroup),
			Groups:            translateUserGroups3to2(u.Groups),
			NoUserGroup:       boolV(u.NoUserGroup),
			NoLogInit:         boolV(u.NoLogInit),
			Shell:             strV(u.Shell),
			System:            boolV(u.System),
		})
	}
	return
}

func translateUserSSH3to2(in []new.SSHAuthorizedKey) (ret []old.SSHAuthorizedKey) {
	for _, k := range in {
		ret = append(ret, old.SSHAuthorizedKey(k))
	}
	return
}

func translateUserGroups3to2(in []new.Group) (ret []old.Group) {
	for _, g := range in {
		ret = append(ret, old.Group(g))
	}
	return
}

func translateGroups3to2(groups []new.PasswdGroup) (ret []old.PasswdGroup) {
	for _, g := range groups {
		ret = append(ret, old.PasswdGroup{
			Name:         g.Name,
			Gid:          g.Gid,
			PasswordHash: strV(g.PasswordHash),
			System:       boolV(g.System),
		})
	}
	return
}

func translateUnits3to2(units []new.Unit) (ret []old.Unit) {
	for _, u := range units {
		ret = append(ret, old.Unit{
			Name:     u.Name,
			Enabled:  u.Enabled,
			Mask:     boolV(u.Mask),
			Contents: strV(u.Contents),
			Dropins:  translateDropins3to2(u.Dropins),
		})
	}
	return
}

func translateDropins3to2(dropins []new.Dropin) (ret []old.SystemdDropin) {
	for _, d := range dropins {
		ret = append(ret, old.SystemdDropin{
			Name:     d.Name,
			Contents: strV(d.Contents),
		})
	}
	return
}

func translateDisks3to2(disks []new.Disk) (ret []old.Disk) {
	for _, d := range disks {
		ret = append(ret, old.Disk{
			Device:     d.Device,
			WipeTable:  boolV(d.WipeTable),
			Partitions: translatePartitions3to2(d.Partitions),
		})
	}
	return
}

func translatePartitions3to2(parts []new.Partition) (ret []old.Partition) {
	for _, p := range parts {
		ret = append(ret, old.Partition{
			Label:    strV(p.Label),
			Number:   p.Number,
			Size:     intV(p.SizeMiB),
			Start:    intV(p.StartMiB),
			TypeGUID: strV(p.TypeGUID),
			GUID:     strV(p.GUID),
		})
	}
	return
}

func translateRaid3to2(raids []new.Raid) (ret []old.Raid) {
	for _, r := range raids {
		ret = append(ret, old.Raid{
			Name:    r.Name,
			Level:   r.Level,
			Devices: translateDevices3to2(r.Devices),
			Spares:  intV(r.Spares),
			Options: translateRaidOptions3to2(r.Options),
		})
	}
	return
}

func translateDevices3to2(devices []new.Device) (ret []old.Device) {
	for _, d := range devices {
		ret = append(ret, old.Device(d))
	}
	return
}

func translateRaidOptions3to2(options []new.RaidOption) (ret []old.RaidOption) {
	for _, o := range options {
		ret = append(ret, old.RaidOption(o))
	}
	return
}

func translateFilesystems3to2(fss []new.Filesystem) (ret []old.Filesystem) {
	for _, f := range fss {

		ret = append(ret, old.Filesystem{
			// To construct a mapping for files/directories, we name the filesystem by path uniquely.
			// TODO: check if its ok to leave out "Path" since we are mapping it via Name later
			Name: strV(f.Path),
			Mount: &old.Mount{
				Device:         f.Device,
				Format:         strV(f.Format),
				WipeFilesystem: boolV(f.WipeFilesystem),
				Label:          f.Label,
				UUID:           f.UUID,
				Options:        translateFilesystemOptions3to2(f.Options),
			},
		})
	}
	return
}

func translateFilesystemOptions3to2(options []new.FilesystemOption) (ret []old.MountOption) {
	for _, o := range options {
		ret = append(ret, old.MountOption(o))
	}
	return
}

func translateNode3to2(n new.Node, fss []string) old.Node {
	fsname := ""
	path := n.Path
	for _, fs := range fss {
		if strings.HasPrefix(n.Path, fs) && len(fs) > len(fsname) {
			fsname = fs
			path = strings.TrimPrefix(n.Path, fsname)
		}
	}
	if len(fsname) == 0 {
		// TODO check if this properly implies root
		fsname = "root"
	}
	return old.Node{
		Filesystem: fsname,
		Path:       path,
		User: &old.NodeUser{
			ID:   n.User.ID,
			Name: strV(n.User.Name),
		},
		Group: &old.NodeGroup{
			ID:   n.Group.ID,
			Name: strV(n.Group.Name),
		},
		Overwrite: n.Overwrite,
	}
}

func translateFiles3to2(files []new.File, fss []string) (ret []old.File) {
	for _, f := range files {
		file := old.File{
			Node: translateNode3to2(f.Node, fss),
			FileEmbedded1: old.FileEmbedded1{
				Mode: f.Mode,
			},
		}
		// TODO check if spec 3 Append and Contents can simultaneously exist
		if f.FileEmbedded1.Contents.Source != nil {
			file.FileEmbedded1.Contents = old.FileContents{
				Compression: strV(f.Contents.Compression),
				Source:      strV(f.Contents.Source),
			}
			file.FileEmbedded1.Append = false
			ret = append(ret, file)
		} else if f.FileEmbedded1.Append != nil {
			for _, fc := range f.FileEmbedded1.Append {
				appendFile := old.File{
					Node:          file.Node,
					FileEmbedded1: file.FileEmbedded1,
				}
				appendFile.FileEmbedded1.Contents = old.FileContents{
					Compression: strV(fc.Compression),
					Source:      strV(fc.Source),
				}
				appendFile.FileEmbedded1.Append = true
				ret = append(ret, appendFile)
			}
		}
	}
	return
}

func translateLinks3to2(links []new.Link, fss []string) (ret []old.Link) {
	for _, l := range links {
		ret = append(ret, old.Link{
			Node: translateNode3to2(l.Node, fss),
			LinkEmbedded1: old.LinkEmbedded1{
				Hard:   boolV(l.Hard),
				Target: l.Target,
			},
		})
	}
	return
}

func translateDirectories3to2(dirs []new.Directory, fss []string) (ret []old.Directory) {
	for _, d := range dirs {
		ret = append(ret, old.Directory{
			Node: translateNode3to2(d.Node, fss),
			DirectoryEmbedded1: old.DirectoryEmbedded1{
				Mode: d.Mode,
			},
		})
	}
	return
}

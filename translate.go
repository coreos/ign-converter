package ign2to3

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	old "github.com/coreos/ignition/config/v2_4_experimental/types"
	oldValidate "github.com/coreos/ignition/config/validate"
	"github.com/coreos/ignition/v2/config/v3_0/types"
)

// Error definitions

// Error type for when a filesystem is referenced in a config but there's no mapping to where
// it should be mounted (i.e. `path` in v3+ configs)
type NoFilesystemError string

func (e NoFilesystemError) Error() string {
	return fmt.Sprintf("Config defined filesystem %q but no mapping was defined", string(e))
}

// DuplicateInodeError is for when files, directories, or links both specify the same path
type DuplicateInodeError struct {
	Old string // first occurance of the path
	New string // second occurance of the path
}

func (e DuplicateInodeError) Error() string {
	return fmt.Sprintf("Config has conflicting inodes: %q and %q", e.Old, e.New)
}

// UsesOwnLinkError is for when files, directories, or links use symlinks defined in the config
// in their own path. This is disallowed in v3+ configs.  type UsesOwnLinkError struct { LinkPath string
type UsesOwnLinkError struct {
	LinkPath string
	Name     string
}

func (e UsesOwnLinkError) Error() string {
	return fmt.Sprintf("%s uses link in config %q", e.Name, e.LinkPath)
}

// UsesNetworkdError is the error for inlcuding networkd configs
var UsesNetworkdError = errors.New("Config includes a networkd section")

// Check returns if the config is translatable but does not do any translation.
// fsMap is a map from v2 filesystem names to the paths under which they should
// be mounted in v3.
func Check(cfg old.Config, fsMap map[string]string) error {
	// TODO: validate cfg
	rpt := oldValidate.ValidateWithoutSource(reflect.ValueOf(cfg))
	if rpt.IsFatal() || rpt.IsDeprecated() {
		// disallow any deprecated fields
		return fmt.Errorf("Invalid input config:\n%s", rpt.String())
	}

	if len(cfg.Networkd.Units) != 0 {
		return UsesNetworkdError
	}

	// check that all filesystems have a path
	fsMap["root"] = "/"
	for _, fs := range cfg.Storage.Filesystems {
		if _, ok := fsMap[fs.Name]; !ok {
			return NoFilesystemError(fs.Name)
		}
	}

	// check that there are no duplicates with files, links, or directories
	// from path to a pretty-printing description of the entry
	entryMap := map[string]string{}
	links := make([]string, 0, len(cfg.Storage.Links))
	// build up a list of all the links we write. We're not allow to use links
	// that we write
	for _, link := range cfg.Storage.Links {
		path := filepath.Join("/", fsMap[link.Filesystem], link.Path)
		links = append(links, path)
	}

	for _, file := range cfg.Storage.Files {
		path := filepath.Join("/", fsMap[file.Filesystem], file.Path)
		name := fmt.Sprintf("File: %s", path)
		if duplicate, isDup := entryMap[path]; isDup {
			return DuplicateInodeError{duplicate, name}
		}
		if l := checkPathUsesLink(links, path); l != "" {
			return &UsesOwnLinkError{
				LinkPath: l,
				Name:     name,
			}
		}
		entryMap[path] = name
	}
	for _, dir := range cfg.Storage.Directories {
		path := filepath.Join("/", fsMap[dir.Filesystem], dir.Path)
		name := fmt.Sprintf("Directory: %s", path)
		if duplicate, isDup := entryMap[path]; isDup {
			return DuplicateInodeError{duplicate, name}
		}
		if l := checkPathUsesLink(links, path); l != "" {
			return &UsesOwnLinkError{
				LinkPath: l,
				Name:     name,
			}
		}
		entryMap[path] = name
	}
	for _, link := range cfg.Storage.Links {
		path := filepath.Join("/", fsMap[link.Filesystem], link.Path)
		name := fmt.Sprintf("Directory: %s", path)
		if duplicate, isDup := entryMap[path]; isDup {
			return &DuplicateInodeError{duplicate, name}
		}
		entryMap[path] = name
		if l := checkPathUsesLink(links, path); l != "" {
			return &UsesOwnLinkError{
				LinkPath: l,
				Name:     name,
			}
		}
	}
	return nil
}

func checkPathUsesLink(links []string, path string) string {
	for _, l := range links {
		if strings.HasPrefix(path, l) {
			return l
		}
	}
	return ""
}

func Translate(cfg old.Config, fsMap map[string]string) (types.Config, error) {
	if err := Check(cfg, fsMap); err != nil {
		return types.Config{}, err
	}
	res := types.Config{
		// Ignition section
		Ignition: types.Ignition{
			Version: "3.0.0",
			Config: types.IgnitionConfig{
				Replace: translateCfgRef(cfg.Ignition.Config.Replace),
				Merge: translateCfgRefs(cfg.Ignition.Config.Append),
			},
			Security: types.Security{
				TLS: types.TLS{
					CertificateAuthorities: translateCAs(cfg.Ignition.Security.TLS.CertificateAuthorities),
				},
			},
			Timeouts: types.Timeouts{
				HTTPResponseHeaders: cfg.Ignition.Timeouts.HTTPResponseHeaders,
				HTTPTotal: cfg.Ignition.Timeouts.HTTPTotal,
			},
		},
		// Passwd section
		Passwd: types.Passwd{
			Users: translateUsers(cfg.Passwd.Users),
			Groups: translateGroups(cfg.Passwd.Groups),
		},
		Systemd: types.Systemd{
			Units: translateUnits(cfg.Systemd.Units),
		},
		Storage: types.Storage{
			Disks: translateDisks(cfg.Storage.Disks),
			Raid: translateRaid(cfg.Storage.Raid),
			Filesystems: translateFilesystems(cfg.Storage.Filesystems),
			Files: translateFiles(cfg.Storage.Files),
			Directories: translateDirectories(cfg.Storage.Directories),
			Links: translateLinks(cfg.Storage.Links),
		},
	}
	return res, nil
}

func translateCfgRef(ref *old.ConfigReference) (ret types.ConfigReference) {
	if ref == nil {
		return
	}
	ret.Source = &ref.Source
	ret.Verification.Hash = ref.Verification.Hash
	return
}

func translateCfgRefs(refs []old.ConfigReference) (ret []types.ConfigReference) {
	for _, ref := range refs {
		ret = append(ret, translateCfgRef(&ref))
	}
	return
}

func translateCAs(refs []old.CaReference) (ret []types.CaReference) {
	for _, ref := range refs {
		ret = append(ret, types.CaReference{
			Source: ref.Source,
			Verification: types.Verification{
				Hash: ref.Verification.Hash,
			},
		})
	}
	return
}

func translateUsers(users []old.PasswdUser) (ret []types.PasswdUser) {
	for _, u := range users {
		ret = append(ret, types.PasswdUser{
			Name: u.Name,
			PasswordHash: u.PasswordHash,
			SSHAuthorizedKeys: translateUserSSH(u.SSHAuthorizedKeys),
			UID: u.UID,
			Gecos: strP(u.Gecos),
			HomeDir: strP(u.HomeDir),
			NoCreateHome: boolP(u.NoCreateHome),
			PrimaryGroup: strP(u.PrimaryGroup),
			Groups: translateUserGroups(u.Groups),
			NoUserGroup: boolP(u.NoUserGroup),
			NoLogInit: boolP(u.NoLogInit),
			Shell: strP(u.Shell),
			System: boolP(u.System),
		})
	}
	return
}

func translateUserSSH(in []old.SSHAuthorizedKey) (ret []types.SSHAuthorizedKey) {
	for _, k := range in {
		ret = append(ret, types.SSHAuthorizedKey(k))
	}
	return
}

func translateUserGroups(in []old.Group) (ret []types.Group) {
	for _, g := range in {
		ret = append(ret, types.Group(g))
	}
	return
}

func translateGroups(groups []old.PasswdGroup) (ret []types.PasswdGroup) {
	return
}

func translateUnits(units []old.Unit) (ret []types.Unit) {
	return
}

func translateDisks(disks []old.Disk) (ret []types.Disk) {
	return
}

func translateRaid(raids []old.Raid) (ret []types.Raid) {
	return
}

func translateFilesystems(fss []old.Filesystem) (ret []types.Filesystem) {
	return
}

func translateFiles(files []old.File) (ret []types.File) {
	return
}

func translateLinks(links []old.Link) (ret []types.Link) {
	return
}

func translateDirectories(dirs []old.Directory) (ret []types.Directory) {
	return
}

package ign2to3

import (
	"fmt"
	"path/filepath"
	"strings"

	old "github.com/coreos/ignition/config/v2_4_experimental/types"
	"github.com/coreos/ignition/v2/config/v3_0/types"
)

type NoFilesystemError string

func (e NoFilesystemError) Error() string {
	return fmt.Sprintf("Config defined filesystem %q but no mapping was defined", string(e))
}

type DuplicateInodeError struct {
	Old string
	New string
}

func (e DuplicateInodeError) Error() string {
	return fmt.Sprintf("Config has conflicting inodes: %q and %q", e.Old, e.New)
}

type UsesOwnLinkError struct {
	LinkPath string
	Name     string
}

func (e UsesOwnLinkError) Error() string {
	return fmt.Sprintf("%s uses link in config %q", e.Name, e.LinkPath)
}

// Check returns if the config is translatable but does not do any translation.
// fsMap is a map from v2 filesystem names to the paths under which they should
// be mounted in v3.
func Check(cfg old.Config, fsMap map[string]string) error {
	// TODO: validate cfg

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
	return types.Config{}, nil
}

package structure

import (
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PackageType string

const (
	// External packages are external dependencies
	External PackageType = "external"
	// Common packages contain shared utilities and can import only External packages
	Common PackageType = "common"
	// Service packages are deployable services and can import Common and External packages
	Service PackageType = "service"
	// Utility packages can import whatever they want
	Utility PackageType = "utility"
)

type PackageInfo struct {
	// path is the path that can be used to run tests on a specific package.
	path string
	// name is the friendly name of a package.
	name string
	typ  PackageType
}

var RootPkg = PackageInfo{
	path: ".",
	name: "project root",
	typ:  Common,
}

func (p PackageInfo) String() string {
	return fmt.Sprintf("%s (%s)", p.name, p.typ)
}

func (p PackageInfo) Path() string {
	return fmt.Sprintf("%s/...", p.path)
}

func (p PackageInfo) Name() string {
	return p.name
}

func (p PackageInfo) Type() PackageType {
	return p.typ
}

type baseStruct struct {
	base     string
	rootDir  string
	external []string
	utility  []string
	common   []string
	service  []string
}

type BaseStruct interface {
	GetPackageInfo(pkgPath string) PackageInfo
	BuildPackageTree(path string) ImportTree
	GetChangedPackages(r io.Reader) PackageSet
}

func NewBaseStruct() BaseStruct {
	res := baseStruct{
		base:     viper.GetString("module_name"),
		rootDir:  viper.GetString("root_dir"),
		external: viper.GetStringSlice("folders.external"),
		utility:  viper.GetStringSlice("folders.utility"),
		common:   viper.GetStringSlice("folders.common"),
		service:  viper.GetStringSlice("folders.service"),
	}
	return &res
}

func (s *baseStruct) match(pkgPath string, relPath string, typ PackageType) (match bool, info PackageInfo) {
	prefix := strings.Join([]string{s.base, relPath}, "/")
	hasWildcard := strings.HasSuffix(prefix, "*")
	if hasWildcard {
		prefix = strings.TrimSuffix(prefix, "*")
	}
	if strings.HasPrefix(pkgPath, prefix) {
		var res PackageInfo
		if hasWildcard {
			name := strings.Split(strings.Replace(pkgPath, prefix, "", 1), "/")[0]
			res = PackageInfo{
				path: strings.Join([]string{prefix, name}, "/"),
				name: name,
				typ:  typ,
			}
		} else {
			res = PackageInfo{
				path: prefix,
				name: relPath,
				typ:  typ,
			}
		}
		log.Debugf("found %s matched %s against %s with base %s", res, pkgPath, relPath, s.base)
		return true, res
	}
	return false, PackageInfo{}
}

func (s *baseStruct) GetPackageInfo(pkgPath string) PackageInfo {
	if !strings.HasPrefix(pkgPath, s.base) {
		return PackageInfo{
			path: pkgPath,
			name: pkgPath,
			typ:  External,
		}
	}

	for _, path := range s.utility {
		if ok, match := s.match(pkgPath, path, Utility); ok {
			return match
		}
	}

	for _, path := range s.external {
		if ok, match := s.match(pkgPath, path, External); ok {
			return match
		}
	}

	for _, path := range s.common {
		if ok, match := s.match(pkgPath, path, Common); ok {
			return match
		}
	}

	for _, path := range s.service {
		if ok, match := s.match(pkgPath, path, Service); ok {
			return match
		}
	}

	return PackageInfo{
		path: pkgPath,
		name: pkgPath,
		typ:  External,
	}
}

func (p PackageInfo) CanImport(o PackageInfo) bool {
	if p.typ == External {
		log.Tracef("skipping external package %s", p)
		return true
	}

	if o.typ == External {
		log.Tracef("package %s can import package %s: imported package is external", p, o)
		return true
	}

	if p.typ == Utility {
		log.Debugf("package %s can import package %s: package is utility", p, o)
		return true
	}

	if o.typ == Utility {
		log.Errorf("package %s cannot import package %s: no package can import utilities", p, o)
		return false
	}

	if o.typ == Common {
		log.Debugf("package %s can import package %s: imported package is common", p, o)
		return true
	}

	if p.typ == Common {
		log.Errorf("package %s cannot import package %s: common package cannot import service", p, o)
		return false
	}

	if p.path != o.path {
		log.Errorf("package %s cannot import package %s: packages belong to different services", p, o)
		return false
	}

	return true
}

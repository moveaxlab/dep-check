package structure

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
)

type packageTree struct {
	m map[PackageInfo]map[PackageInfo]bool
}

type importTree struct {
	packageTree
}

type dependencyTree struct {
	packageTree
}

type PackageTree interface {
	Enumerate() map[PackageInfo][]PackageInfo
}

type ImportTree interface {
	PackageTree
	ToDependencyTree() DependencyTree
	PrintEnumerate() string
}

func (t importTree) PrintEnumerate() string {
	res := ""
	for pkg, imps := range t.Enumerate() {
		res += fmt.Sprintf("%s:\n", pkg)
		for _, imp := range imps {
			res += fmt.Sprintf("\t%s\n", imp)
		}
	}
	return res
}

type DependencyTree interface {
	PackageTree
	ToImportTree() ImportTree
	ExpandDependencies(set PackageSet)
	PrintEnumerate() string
}

func (t dependencyTree) PrintEnumerate() string {
	res := ""
	for pkg, deps := range t.Enumerate() {
		res += fmt.Sprintf("%s:\n", pkg)
		for _, dep := range deps {
			res += fmt.Sprintf("\t%s\n", dep)
		}
	}
	return res
}

func (s *baseStruct) BuildPackageTree(path string) ImportTree {
	log.Debugf("building package tree from %s", path)
	cfg := &packages.Config{}
	cfg.Mode = packages.NeedImports | packages.NeedName
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		panic(fmt.Errorf("failed to load packages: %w", err))
	}
	res := importTree{
		packageTree{
			m: make(map[PackageInfo]map[PackageInfo]bool),
		},
	}
	for _, pkg := range pkgs {
		pkgInfo := s.GetPackageInfo(pkg.PkgPath)
		if pkgInfo.typ == External {
			continue
		}
		if _, ok := res.m[pkgInfo]; !ok {
			log.Debugf("found package %s", pkgInfo)
			res.m[pkgInfo] = make(map[PackageInfo]bool)
		}
		for _, imp := range pkg.Imports {
			impInfo := s.GetPackageInfo(imp.PkgPath)
			if impInfo.typ == External {
				continue
			}
			if !res.m[pkgInfo][impInfo] {
				log.Debugf("package %s depends on %s", pkgInfo, impInfo)
			}
			res.m[pkgInfo][impInfo] = true
		}
	}
	return res
}

func (t packageTree) flip() packageTree {
	res := packageTree{
		m: make(map[PackageInfo]map[PackageInfo]bool),
	}
	for pkg, imps := range t.m {
		for imp := range imps {
			if _, ok := res.m[imp]; !ok {
				res.m[imp] = make(map[PackageInfo]bool)
			}
			res.m[imp][pkg] = true
		}
	}
	return res
}

func (t importTree) ToDependencyTree() DependencyTree {
	return dependencyTree{t.packageTree.flip()}
}

func (t dependencyTree) ToImportTree() ImportTree {
	return importTree{t.packageTree.flip()}
}

func (t packageTree) Enumerate() map[PackageInfo][]PackageInfo {
	res := make(map[PackageInfo][]PackageInfo)
	for pkg, imps := range t.m {
		res[pkg] = make([]PackageInfo, 0, len(imps))
		for imp := range imps {
			res[pkg] = append(res[pkg], imp)
		}
	}
	return res
}

func (t dependencyTree) ExpandDependencies(set PackageSet) {
	log.Debugf("expanding dependency tree")

	if set.Contains(RootPkg) {
		log.Debugf("detected root package change, adding all packages to dependency tree")
		for pkg := range t.ToImportTree().Enumerate() {
			set.Add(pkg)
		}
		return
	}
	time.Sleep(10 * time.Second)

	changed := true
	visited := NewPackageSet()
	for changed {
		changed = false

		for _, pkg := range set.Enumerate() {
			deps, ok := t.m[pkg]
			if !ok {
				log.Debugf("package %s has no dependencies, skipping it", pkg)
				continue
			}

			for dep := range deps {
				if !visited.Contains(dep) && !set.Contains(dep) {
					log.Debugf("package %s has package %s as a dependency, expanding tree", pkg, dep)
					changed = true
					set.Add(dep)
					visited.Add(dep)
				}
			}
		}
	}
}

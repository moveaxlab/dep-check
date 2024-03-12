package structure

import (
	"fmt"

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
}

type DependencyTree interface {
	PackageTree
	ToImportTree() ImportTree
	ExpandDependencies(set PackageSet)
}

func (s *baseStruct) BuildPackageTree(path string) ImportTree {
	cfg := &packages.Config{}
	cfg.Mode = packages.NeedImports | packages.NeedName
	pkgs, err := packages.Load(cfg, path)
	log.Debugf("%v", pkgs)
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
	if set.Contains(RootPkg) {
		log.Debugf("detected root package change, adding all packages")
		for pkg := range t.ToImportTree().Enumerate() {
			set.Add(pkg)
		}
		return
	}

	changed := true
	for changed {
		changed = false
		for _, pkg := range set.Enumerate() {
			if deps, ok := t.m[pkg]; ok {
				for dep := range deps {
					if !set.Contains(dep) {
						log.Debugf("%s depends on %s", dep, pkg)
						changed = true
						set.Add(dep)
					}
				}
			} else {
				log.Debugf("skipping invalid package %s", pkg)
				set.Remove(pkg)
				changed = true
			}
		}
	}
}

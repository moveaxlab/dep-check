package structure

type packageSet struct {
	m map[PackageInfo]bool
}

type PackageSet interface {
	Add(p PackageInfo)
	Contains(p PackageInfo) bool
	Remove(p PackageInfo)
	Enumerate() []PackageInfo
}

func NewPackageSet() PackageSet {
	return &packageSet{
		m: make(map[PackageInfo]bool),
	}
}

func (s *packageSet) Add(p PackageInfo) {
	s.m[p] = true
}

func (s *packageSet) Contains(p PackageInfo) bool {
	_, ok := s.m[p]
	return ok
}

func (s *packageSet) Remove(p PackageInfo) {
	delete(s.m, p)
}

func (s *packageSet) Enumerate() []PackageInfo {
	res := make([]PackageInfo, 0, len(s.m))
	for pkg := range s.m {
		res = append(res, pkg)
	}
	return res
}

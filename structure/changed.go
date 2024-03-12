package structure

import (
	"bufio"
	"io"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (s *baseStruct) GetChangedPackages(r io.Reader) PackageSet {
	scanner := bufio.NewScanner(r)

	res := NewPackageSet()

	for scanner.Scan() {
		changedFile := scanner.Text()

		if !strings.HasPrefix(changedFile, s.rootDir) {
			continue
		}

		changedFile = strings.TrimPrefix(changedFile, s.rootDir)

		if changedFile == "go.mod" || changedFile == "go.sum" {
			res.Add(RootPkg)
			return res
		}

		changedPath := path.Join(s.base, changedFile)

		if path.Dir(changedPath) == s.base {
			log.Debugf("skipping root directory file %s", changedFile)
			continue
		}

		info := s.GetPackageInfo(changedPath)

		if !res.Contains(info) {
			log.Debugf("found changed package %s", info)
		}

		res.Add(info)
	}

	return res
}

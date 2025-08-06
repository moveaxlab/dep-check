package cmd

import (
	"github.com/moveaxlab/dep-check/config"
	"github.com/moveaxlab/dep-check/graph"
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// const (
// 	goModFileName = "go.mod"
// )

var (
	// moduleName string
	// // packagesConfig = &packages.Config{Mode: packages.NeedImports | packages.NeedName}
	jsCreatePrefixWith = func(n string) (pref string) {
		return n
	}
	jsCmd = &cobra.Command{
		Use:   "js",
		Short: "build dependencies graph based on the JavaScript language",
		PreRun: func(cmd *cobra.Command, args []string) {
			jsPackages := make([]*jsPackage, 0)
			err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				slog.Debug(fmt.Sprint("current walking dir ", path))
				slog.Debug(fmt.Sprint(" current file ", d.Name()))
				if strings.Contains(path, "node_modules") || (len(path) > 1 && strings.HasPrefix(path, ".") && d.IsDir()) {
					slog.Debug(" SkipDir")
					return filepath.SkipDir
				}

				if d.Name() == "package.json" {
					jsPackage, err := readPackageJSON(path)
					if err != nil {
						slog.Error(fmt.Sprint("cannot load package ", path))
						return filepath.SkipDir
					}
					jsPackages = append(jsPackages, jsPackage)
				}

				return nil
			})

			if err != nil {
				panic(fmt.Sprint("JavaScript packages can not be loaded", err))
			}

			slog.Debug(fmt.Sprintf("%d Javaxcript packages found", len(jsPackages)))
			if DebugMode {
				for _, n := range jsPackages {
					slog.Debug(fmt.Sprintf(" %s : %s", n.Name, n.Path))
				}
			}

			// Create NODES
			for _, pkg := range jsPackages {
				createOrGetNodeFromJs(pkg)
			}
			slog.Debug(fmt.Sprintf("Nodes num : %d", len(graph.Nodes)))
			if DebugMode {
				for _, n := range graph.Nodes {
					slog.Debug(fmt.Sprintf(" %s : %s", n.Name, n.Path))
				}
			}

			// CREATE EDGES
			for _, pkg := range jsPackages {
				if pkgNode, ok := graph.Nodes[pkg.Name]; ok {
					for imp := range pkg.Imports {
						if impNode, ok := graph.Nodes[imp]; ok {
							e := &graph.Edge{From: pkgNode.Name, To: impNode.Name}
							if _, ok := graph.Edges[e.String()]; !ok {
								graph.Edges[e.String()] = e
							}
						}
					}
				}
			}

			slog.Debug(fmt.Sprintf("Edges num : %d", len(graph.Edges)))
		},
		Run: func(cmd *cobra.Command, args []string) {
			changedPackages := make(map[string]string)

			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				changedFile := scanner.Text()
				slog.Info(fmt.Sprintf("Input file : %s", changedFile))

				if filepath.Dir(changedFile) == "." {
					continue
				}

				filePackage, err := graph.FindNodeFromFilePath(changedFile)
				if err != nil {
					slog.Warn("No node found")
					continue
				}
				slog.Debug("Corresponding node :")
				slog.Debug(fmt.Sprintf("%s : %s", filePackage.Name, filePackage.Path))

				if _, ok := changedPackages[filePackage.Name]; !ok {
					for _, parent := range filePackage.GetParents() {
						if _, ok := changedPackages[parent.Name]; !ok {
							changedPackages[parent.Name] = strings.Join([]string{parent.Path, "/..."}, "")
							fmt.Println(changedPackages[parent.Name])
						}
					}
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(jsCmd)
}

type jsPackage struct {
	Name    string            `json:"name"`
	Imports map[string]string `json:"dependencies"`
	Path    string
}

func readPackageJSON(filePath string) (*jsPackage, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error(fmt.Sprint("Cannot read file ", filePath, err))
		return nil, err
	}

	jsPackage := jsPackage{Path: strings.Replace(filePath, "/package.json", "", 1)}
	err = json.Unmarshal(fileContent, &jsPackage)
	if err != nil {
		slog.Error(fmt.Sprint("Cannot unmarshal file ", filePath, err))
		return nil, err
	}

	return &jsPackage, nil
}

func createOrGetNodeFromJs(pkg *jsPackage) (n *graph.Node) {
	for _, artif := range config.DepCheckConfig.Artifacts {
		hasWildcard := strings.HasSuffix(artif, "*")
		if hasWildcard {
			artif = strings.TrimSuffix(artif, "*")
		}
		if strings.HasPrefix(pkg.Path, artif) {
			name := artif
			if hasWildcard {
				name = pkg.Name
			}
			if _, ok := graph.Nodes[name]; !ok {
				n = &graph.Node{Name: name, Path: pkg.Path}
				graph.Nodes[name] = n
			} else {
				n = graph.Nodes[name]
			}
			return n
		}
	}
	return nil
}

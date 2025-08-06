package cmd

import (
	"github.com/moveaxlab/dep-check/config"
	"github.com/moveaxlab/dep-check/graph"
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

const (
	goModFileName = "go.mod"
)

var (
	moduleName string
	goCmd      = &cobra.Command{
		Use:   "go",
		Short: "build dependencies graph based on the go language",
		PreRun: func(cmd *cobra.Command, args []string) {
			// Load Module Name
			data, err := os.ReadFile(goModFileName)
			if err != nil {
				panic(fmt.Sprint("No go.mod file", err))
			}
			modFile, err := modfile.Parse(goModFileName, data, nil)
			if err != nil {
				panic("No go.mod file")
			}
			moduleName = modFile.Module.Mod.Path
			slog.Debug(fmt.Sprintf("Name of the module is : %s", moduleName))

			// Load packages
			packagesConfig := &packages.Config{Mode: packages.NeedImports | packages.NeedName}

			goPackages, err := packages.Load(packagesConfig, "./...")
			if err != nil {
				panic(fmt.Sprint("Go packages can not be loaded", err))
			}

			slog.Debug(fmt.Sprintf("%d go packages found", len(goPackages)))
			if DebugMode {
				for _, n := range goPackages {
					slog.Debug(fmt.Sprintf(" %s : %s", n.Name, n.PkgPath))
				}
			}

			// Build Graph
			for _, pkg := range goPackages {
				pkgNode := createOrGetNodeFromGo(pkg.PkgPath)
				if pkgNode != nil {
					for _, imp := range pkg.Imports {
						impNode := createOrGetNodeFromGo(imp.PkgPath)
						if impNode != nil && pkgNode.Name != impNode.Name {
							e := &graph.Edge{From: pkgNode.Name, To: impNode.Name}
							if _, ok := graph.Edges[e.String()]; !ok {
								graph.Edges[e.String()] = e
							}
						}
					}
				}
			}

			slog.Debug(fmt.Sprintf("Nodes num : %d", len(graph.Nodes)))
			if DebugMode {
				for _, n := range graph.Nodes {
					slog.Debug(fmt.Sprintf(" %s : %s", n.Name, n.Path))
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

				filePackage := createOrGetNodeFromGo(path.Join(moduleName, changedFile))
				if filePackage == nil {
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
	rootCmd.AddCommand(goCmd)
}

func createOrGetNodeFromGo(pkgPath string) (n *graph.Node) {
	for _, artif := range config.DepCheckConfig.Artifacts {
		prefix := path.Join(moduleName, artif)
		hasWildcard := strings.HasSuffix(prefix, "*")
		if hasWildcard {
			prefix = strings.TrimSuffix(prefix, "*")
		}
		if strings.HasPrefix(pkgPath, prefix) {
			name := artif
			if hasWildcard {
				name = strings.Replace(artif, "*", strings.Split(strings.Replace(pkgPath, prefix, "", 1), "/")[0], 1)
				prefix = path.Join(moduleName, name)
			}
			if _, ok := graph.Nodes[name]; !ok {
				n = &graph.Node{Name: name, Path: prefix}
				graph.Nodes[name] = n
			} else {
				n = graph.Nodes[name]
			}
			return n
		}
	}
	return nil
}

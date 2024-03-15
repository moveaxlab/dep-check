package cmd

import (
	"fmt"
	"os"

	"github.com/moveaxlab/dep-check/structure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var changedPackagesCommand = &cobra.Command{
	Use:   "changed-packages",
	Short: "detect changed packages from a git diff",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseStruct := structure.NewBaseStruct()

		changedPackages := baseStruct.GetChangedPackages(os.Stdin)

		if changedPackages.Contains(structure.RootPkg) {
			log.Infof("change detected in root package")
			fmt.Fprintln(os.Stdout, structure.RootPkg.Path())
			return nil
		}

		dependencies := baseStruct.BuildPackageTree("./...").
			ToDependencyTree()

		dependencies.ExpandDependencies(changedPackages)

		for _, pkg := range changedPackages.Enumerate() {
			log.Infof("change detected in package %s", pkg)
			fmt.Fprintln(os.Stdout, pkg.Path())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(changedPackagesCommand)
}

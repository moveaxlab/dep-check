package cmd

import (
	"fmt"
	"os"

	"github.com/moveaxlab/dep-check/structure"
	"github.com/spf13/cobra"
)

var changedPackagesCommand = &cobra.Command{
	Use:   "changed-packages",
	Short: "detect changed packages from a git diff",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseStruct := structure.NewBaseStruct()
		changedPackages := baseStruct.GetChangedPackages(os.Stdin)

		if changedPackages.Contains(structure.RootPkg) {
			fmt.Println(structure.RootPkg.Path())
			return nil
		}

		dependencies := baseStruct.BuildPackageTree("./...").
			ToDependencyTree()

		dependencies.ExpandDependencies(changedPackages)

		for _, pkg := range changedPackages.Enumerate() {
			fmt.Println(pkg.Path())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(changedPackagesCommand)
}

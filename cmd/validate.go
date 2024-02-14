package cmd

import (
	"fmt"

	"github.com/moveaxlab/deps-check/structure"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validates the project structure",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseStruct := structure.NewBaseStruct()
		tree := baseStruct.BuildPackageTree("./...")

		valid := true
		for pkg, imports := range tree.Enumerate() {
			for _, imp := range imports {
				if !pkg.CanImport(imp) {
					valid = false
				}
			}
		}

		if !valid {
			return fmt.Errorf("invalid imports detected, check error log for more details")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

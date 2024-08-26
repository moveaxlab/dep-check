package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/moveaxlab/dep-check/structure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var changedServicesCommand = &cobra.Command{
	Use:   "changed-services",
	Short: "detect changed services from a git diff",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseStruct := structure.NewBaseStruct()

		changedServices := baseStruct.GetChangedPackages(os.Stdin)

		if changedServices.Contains(structure.RootPkg) {
			log.Infof("change detected in root package")
		}

		dependencies := baseStruct.BuildPackageTree("./...").
			ToDependencyTree()

		dependencies.ExpandDependencies(changedServices)

		var services []string
		for _, pkg := range changedServices.Enumerate() {
			if pkg.Type() == structure.Service {
				services = append(services, pkg.Name())
			}
		}
	
		log.Println("changes detected for services", services)
		fmt.Println(strings.Join(services, " "))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(changedServicesCommand)
}

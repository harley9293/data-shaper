package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var create bool
var export bool

var excelCmd = &cobra.Command{
	Use:   "excel",
	Short: "excel file conversion and generation",
	RunE: func(cmd *cobra.Command, args []string) error {
		if create == export {
			return errors.New("either 'create' or 'export' must be specified, but not both")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(excelCmd)

	excelCmd.Flags().BoolVarP(&create, "create", "c", false, "create")
	excelCmd.Flags().BoolVarP(&export, "export", "e", false, "export")

	excelCmd.DisableFlagsInUseLine = true
}

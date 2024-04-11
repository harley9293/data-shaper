package cmd

import (
	"errors"
	"fmt"
	"github.com/harley9293/data-shaper/internal/converter"
	"github.com/spf13/cobra"
)

func ExcelCmd() *cobra.Command {
	var create bool
	var export bool
	cmd := &cobra.Command{
		Use:   "excel",
		Short: "excel file conversion and generation",
		RunE: func(cmd *cobra.Command, args []string) error {
			if create == export {
				return errors.New("either 'create' or 'export' must be specified, but not both")
			}

			return nil
		},
	}
	cmd.DisableFlagsInUseLine = true

	return cmd
}

func CreateCmd() *cobra.Command {
	var protoPath string
	var excelPath string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an Excel file from a proto file",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := converter.ProtoToExcel(protoPath, excelPath)
			if err != nil {
				return err
			}
			fmt.Println("Excel file created successfully at", excelPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&protoPath, "proto", "", "Path to the proto file")
	cmd.MarkFlagRequired("proto")
	cmd.Flags().StringVar(&excelPath, "excel", "", "Path for the resulting Excel file")
	cmd.MarkFlagRequired("excel")

	return cmd
}

func init() {
	excelCmd := ExcelCmd()
	excelCmd.AddCommand(CreateCmd())
	rootCmd.AddCommand(excelCmd)
}

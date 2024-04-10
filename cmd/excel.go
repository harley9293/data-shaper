package cmd

import (
	"errors"
	"fmt"
	"github.com/harley9293/data-shaper/internal/converter"
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

var protoPath string
var excelPath string
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Excel file from a proto file",
	Long:  `Create an Excel file from a proto file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 调用转换逻辑
		err := converter.ProtoToExcel(protoPath, excelPath)
		if err != nil {
			return err
		}
		fmt.Println("Excel file created successfully at", excelPath)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVar(&protoPath, "proto", "", "Path to the proto file")
	createCmd.MarkFlagRequired("proto")
	createCmd.Flags().StringVar(&excelPath, "excel", "", "Path for the resulting Excel file")
	createCmd.MarkFlagRequired("excel")
	excelCmd.AddCommand(createCmd)

	excelCmd.DisableFlagsInUseLine = true
	rootCmd.AddCommand(excelCmd)
}

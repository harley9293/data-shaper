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
	var outputPath string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an Excel file from a proto file",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := converter.ProtoToExcel(protoPath, outputPath)
			if err != nil {
				return err
			}
			fmt.Println("Excel file created successfully at", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&protoPath, "proto", "p", "", "Path to the proto file")
	cmd.MarkFlagRequired("proto")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path for the resulting Excel directory")
	cmd.MarkFlagRequired("output")

	return cmd
}

func ExportCmd() *cobra.Command {
	var outputPath string
	var excelPath string
	var exportType string
	var protoPath string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export an Excel file to support config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if exportType != "json" {
				return errors.New("unsupported export type, only 'json' is supported")
			}

			err := converter.ExcelToJson(protoPath, excelPath, outputPath)
			if err != nil {
				return err
			}
			fmt.Println("Excel file convert to json successfully at", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to the config file")
	cmd.MarkFlagRequired("output")
	cmd.Flags().StringVarP(&excelPath, "excel", "e", "", "Path to the Excel file")
	cmd.MarkFlagRequired("excel")
	cmd.Flags().StringVarP(&exportType, "type", "t", "", "Type of the config file to export to")
	cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(&protoPath, "proto", "p", "", "Path to the proto file")
	cmd.MarkFlagRequired("proto")

	return cmd
}

func init() {
	excelCmd := ExcelCmd()
	excelCmd.AddCommand(CreateCmd(), ExportCmd())
	rootCmd.AddCommand(excelCmd)
}

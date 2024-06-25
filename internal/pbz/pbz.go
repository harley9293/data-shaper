package pbz

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pbz",
		Short: "config tools with proto file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("must specify valid subcommand")
		},
	}
	cmd.DisableFlagsInUseLine = true
	cmd.AddCommand(excel(), export())
	return cmd
}

func excel() *cobra.Command {
	var protoPath string
	var outputPath string
	cmd := &cobra.Command{
		Use:   "excel",
		Short: "create excel file from proto file",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := protoToExcel(protoPath, outputPath)
			if err != nil {
				return err
			}
			fmt.Println("excel file created successfully at", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&protoPath, "proto", "p", "", "path to the proto file")
	cmd.MarkFlagRequired("proto")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path")
	cmd.MarkFlagRequired("output")

	return cmd
}

func export() *cobra.Command {
	var protoPath string
	var outputPath string
	var excelPath string
	var exportType string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "export excel content to support config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if exportType != "cfg" {
				return errors.New("unsupported export type, only 'cfg' is supported")
			}

			fmt.Println("excel file convert to cfg successfully at", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&protoPath, "proto", "p", "", "path to the proto file")
	cmd.MarkFlagRequired("proto")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path")
	cmd.MarkFlagRequired("output")
	cmd.Flags().StringVarP(&excelPath, "excel", "e", "", "Path to the Excel file")
	cmd.MarkFlagRequired("excel")
	cmd.Flags().StringVarP(&exportType, "type", "t", "", "type of the config file to export to")
	cmd.MarkFlagRequired("type")

	return cmd
}

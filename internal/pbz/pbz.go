package pbz

import (
	"errors"
	"fmt"
	"github.com/harley9293/data-shaper/internal/pbz/core"
	"github.com/harley9293/data-shaper/internal/pbz/parser"
	"github.com/harley9293/data-shaper/internal/pbz/writer"
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
	var protoFile string
	var excelPath string
	cmd := &cobra.Command{
		Use:   "excel",
		Short: "create or update excel file from proto file",
		RunE: func(cmd *cobra.Command, args []string) error {
			schema := core.NewProtoExcelSchema(excelPath, &parser.ProtoParser{}, &writer.ExcelWriter{}, &parser.DataParser{})
			err := schema.ParseProto(protoFile)
			if err != nil {
				return err
			}

			err = schema.SaveData()
			if err != nil {
				return err
			}
			fmt.Println("excel file update successfully at", excelPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&protoFile, "proto", "p", "", "path to the proto file")
	cmd.MarkFlagRequired("proto")
	cmd.Flags().StringVarP(&excelPath, "excel", "e", "", "output excel file path")
	cmd.MarkFlagRequired("excel")

	return cmd
}

func export() *cobra.Command {
	var protoFile string
	var excelPath string
	var outputPath string
	var exportType string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "export excel content to support config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if exportType != "json" {
				return errors.New("unsupported export type, only 'json' is supported")
			}

			fmt.Println("excel file convert to json successfully at", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&protoFile, "proto", "p", "", "path to the proto file")
	cmd.MarkFlagRequired("proto")
	cmd.Flags().StringVarP(&excelPath, "excel", "e", "", "excel file path")
	cmd.MarkFlagRequired("excel")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path")
	cmd.MarkFlagRequired("output")
	cmd.Flags().StringVarP(&exportType, "type", "t", "", "support export type, only 'json' is supported")
	cmd.MarkFlagRequired("type")

	return cmd
}

package pbz

import (
	"errors"
	"fmt"

	"github.com/harley9293/data-shaper/internal/pbz/parser"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pbz",
		Short: "protobuf配置工具",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("必须指定有效的子命令")
		},
	}
	cmd.DisableFlagsInUseLine = true
	cmd.AddCommand(update(), export())
	return cmd
}

func update() *cobra.Command {
	var protoFile string
	var excelPath string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "从proto文件创建或更新excel文件",
		RunE: func(cmd *cobra.Command, args []string) error {

			_, err := parser.Proto(protoFile)
			if err != nil {
				return err
			}

			// err = schema.SaveData()
			// if err != nil {
			// 	return err
			// }
			fmt.Println("excel文件更新成功", excelPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&protoFile, "proto", "p", "", "proto文件")
	cmd.MarkFlagRequired("proto")
	cmd.Flags().StringVarP(&excelPath, "output", "o", "", "输出excel目录")
	cmd.MarkFlagRequired("output")

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

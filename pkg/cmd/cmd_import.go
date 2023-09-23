// Copyright 2023 The KCL Authors. All rights reserved.

package cmd

import (
	"github.com/urfave/cli/v2"
	"kcl-lang.io/kcl-go/pkg/tools/gen"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/reporter"
	"os"
)

// NewImportCmd new a Command for `kpm import`.
func NewImportCmd(_ *client.KpmClient) *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "import",
		Usage:  "convert other formats to KCL file",
		Description: `import converts other formats to KCL file.

	Supported conversion modes:
	json:            convert JSON data to KCL data
	yaml:            convert YAML data to KCL data
	gostruct:        convert Go struct to KCL schema
	jsonschema:      convert JSON schema to KCL schema
	terraformschema: convert Terraform schema to KCL schema
	auto:            automatically detect the input format`,
		ArgsUsage: "<file>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "mode",
				Aliases:     []string{"m"},
				Usage:       "mode of import",
				DefaultText: "auto",
				Value:       "auto",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "output filename",
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "force overwrite output file",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				reporter.Report("kpm: invalid arguments")
				reporter.ExitWithReport("kpm: run 'kpm import help' for more information.")
			}
			inputFile := c.Args().First()

			opt := &gen.GenKclOptions{}
			switch c.String("mode") {
			case "json":
				opt.Mode = gen.ModeJson
			case "yaml":
				opt.Mode = gen.ModeYaml
			case "gostruct":
				opt.Mode = gen.ModeGoStruct
			case "jsonschema":
				opt.Mode = gen.ModeJsonSchema
			case "terraformschema":
				opt.Mode = gen.ModeTerraformSchema
			case "auto":
				opt.Mode = gen.ModeAuto
			default:
				reporter.Report("kpm: invalid mode: ", c.String("mode"))
				reporter.ExitWithReport("kpm: run 'kpm import help' for more information.")
			}

			outputFile := c.String("output")
			if outputFile == "" {
				outputFile = "generated.k"
				reporter.Report("kpm: output file not specified, use default: ", outputFile)
			}

			if _, err := os.Stat(outputFile); err == nil && !c.Bool("force") {
				reporter.ExitWithReport("kpm: output file already exist, use --force to overwrite: ", outputFile)
			}

			outputWriter, err := os.Create(outputFile)
			if err != nil {
				reporter.ExitWithReport("kpm: failed to create output file: ", outputFile)
			}

			return gen.GenKcl(outputWriter, inputFile, nil, opt)
		},
	}
}

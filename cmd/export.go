/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hgkcho/matpw/pkg/password"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export password data",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: run,
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.PersistentFlags().BoolP("csv", "c", false, "export as csv")
	exportCmd.PersistentFlags().StringP("output", "o", "", "export file path")
	viper.BindPFlag("csv", exportCmd.PersistentFlags().Lookup("csv"))
	viper.BindPFlag("output", exportCmd.PersistentFlags().Lookup("output"))
}

func run(cmd *cobra.Command, args []string) error {
	m := password.Manager{}
	m.Init()

	// determine export path, stdout or csv file
	if viper.GetString("output") != "" {
		if err := exportAsCSV(m); err != nil {
			return nil
		}
	} else {
		if err := m.ToCSV(os.Stdout); err != nil {
			return nil
		}
	}
	return nil
}

func exportAsCSV(m password.Manager) error {
	var f *os.File

	output := viper.GetString("output")
	if output == "" {
		curDir, err := os.Getwd()
		if err != nil {
			return err
		}
		output = curDir + "/matpw.csv"
	}

	_, err := os.Stat(output)
	// when file is exist, ask you can overwrite
	if err == nil {
		prompt := promptui.Select{
			Label: fmt.Sprintf("Overwrite this file : %s [y/N]", output),
			Items: []string{"Yes", "No"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		if result == "No" {
			fmt.Println("cancelled! ")
			os.Exit(1)
		}
	}

	f, err = os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := m.ToCSV(f); err != nil {
		return nil
	}

	return nil
}

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
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hgkcho/matpw/pkg/password"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type search struct {
	meta
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search password matrix with title",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(viper.AllKeys())
		s := &search{}
		if err := s.meta.init(); err != nil {
			return err
		}
		return s.run(args)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (s *search) run(args []string) error {
	funcMap := promptui.FuncMap
	funcMap["time"] = humanize.Time
	templates := &promptui.SelectTemplates{
		Label:    "Select a service",
		Active:   promptui.IconSelect + " {{ .Service | cyan }}",
		Inactive: "  {{ .Service | faint }}",
		Selected: promptui.IconGood + " {{ .Service }}",
		// 		Details: `
		// {{ "ID:" | faint }}	{{ .ID }}
		// {{ "Description:" | faint }}	{{ .Title }}
		// {{ "Account:" | faint }}	{{ not .Account }}
		// {{ "Last modified:" | faint }}	{{ .UpdatedAt | time }}
		// {{ "Content:" | faint }}	{{ .Content | head }}
		// 		`,
		Details: `
{{ "Service:" | faint }}	{{ .Service }}
{{ "Account:" | faint }}	{{ .Account }}
{{ "Last modified:" | faint }}	{{ .UpdatedAt | time }}
		`,
		FuncMap: funcMap,
	}

	searcher := func(input string, index int) bool {
		p := s.passwords[index]
		name := strings.Replace(strings.ToLower(p.Service), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:             "select a service",
		Items:             s.passwords,
		Searcher:          searcher,
		Templates:         templates,
		StartInSearchMode: true,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return err
	}
	pw := prompt.Items.([]password.Password)
	selected := pw[i]

	fmt.Fprintf(os.Stdout, "[service]: %v \n", selected.Service)
	fmt.Fprintf(os.Stdout, "[account]: %v \n", selected.Account)
	fmt.Fprintf(os.Stdout, "[descripiton]: %v \n", selected.Descripiton)
	fmt.Fprintf(os.Stdout, "[createdAt]: %v \n", selected.CreatedAt)
	fmt.Fprintf(os.Stdout, "[updatedAt]: %v \n", selected.UpdatedAt)
	fmt.Println(selected.Render())
	return nil
}

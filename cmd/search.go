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
	"golang.org/x/crypto/ssh/terminal"
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

func render(pwSet string) (ret string) {
	for k, v := range pwSet {
		if k == 0 {
			ret += fmt.Sprintln("------------------------------------")
		}
		if k%5 == 4 {
			ret += fmt.Sprintf("|  %s  |\n", string(v))
			ret += fmt.Sprintln("                                            ")
		} else {
			ret += fmt.Sprintf("|  %s  ", string(v))
		}
		if k == 24 {
			ret += fmt.Sprintln("------------------------------------")
		}
	}
	return "-------------------"
}

func (s *search) run(args []string) error {
	funcMap := promptui.FuncMap
	funcMap["render"] = render
	funcMap["time"] = humanize.Time
	templates := &promptui.SelectTemplates{
		Label:    "Select a title",
		Active:   promptui.IconSelect + " {{ .Title | cyan }}",
		Inactive: "  {{ .Title | faint }}",
		// Selected: promptui.IconGood + " {{ .Title }}",
		// 		Details: `
		// {{ "ID:" | faint }}	{{ .ID }}
		// {{ "Description:" | faint }}	{{ .Title }}
		// {{ "Account:" | faint }}	{{ not .Account }}
		// {{ "Last modified:" | faint }}	{{ .UpdatedAt | time }}
		// {{ "Content:" | faint }}	{{ .Content | head }}
		// 		`,
		Details: `
{{ "Title:" | faint }}	{{ .Title }}
{{ "Account:" | faint }}	{{ .Account }}
		`,
		FuncMap: funcMap,
	}

	searcher := func(input string, index int) bool {
		p := s.passwords[index]
		name := strings.Replace(strings.ToLower(p.Title), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:             "select a title",
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

	fmt.Fprintf(os.Stdout,"[service]: %v \n",selected.Title)
	fmt.Fprintf(os.Stdout,"[account]: %v \n",selected.Account)
	fmt.Fprintf(os.Stdout,"[descripiton]: %v \n",selected.Descripiton)
	fmt.Println(selected.Render())
	return nil
}

func head(content string) string {
	wrap := func(line string) string {
		line = strings.ReplaceAll(line, "\t", "  ")
		id := int(os.Stdout.Fd())
		width, _, _ := terminal.GetSize(id)
		if width < 10 {
			return line
		}
		if len(line) < width-10 {
			return line
		}
		return line[:width-10] + "..."
	}
	lines := strings.Split(content, "\n")
	content = "\n"
	for i := 0; i < len(lines); i++ {
		if i > 4 {
			content += "  ...\n"
			break
		}
		content += "  " + wrap(lines[i]) + "\n"
	}
	return content
}

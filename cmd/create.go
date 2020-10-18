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
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hgkcho/matpw/pkg/password"
	"github.com/spf13/cobra"
)

type create struct {
	meta
	newpw password.Password
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create generate password matrix",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := &create{}
		err := c.init()
		if err != nil {
			return err
		}
		fmt.Print("input service: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		c.newpw.Title = input.Text()

		fmt.Print("input account: ")
		input.Scan()
		c.newpw.Account = input.Text()

		fmt.Print("input descripiton: ")
		input.Scan()
		c.newpw.Descripiton = input.Text()

		return c.run(args)
	},
}

func (c *create) run(args []string) error {
		err := c.newpw.Create()
		if err != nil {
			return err
		}
		c.meta.passwords = append(c.meta.passwords, c.newpw)
		f, err := os.Create(c.meta.pwFilePath)
		// f, err := os.OpenFile(c.meta.pwFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
		defer f.Close()
		if err != nil {
			return err
		}
		err = json.NewEncoder(f).Encode(c.meta.passwords)
		if err != nil {
			return err
		}
		fmt.Println(c.newpw.Render())
		return nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

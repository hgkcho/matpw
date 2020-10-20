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
	"errors"
	"fmt"
	"os"

	"github.com/hgkcho/matpw/pkg/password"
	"github.com/spf13/cobra"
)

type create struct {
	meta
	newpw   password.Password
	manager password.Manager
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create generate matrix password",
	Long: `Create is the command that generate matrix password.
By following instructions in prompt, answer some simple questions.
Then you can register password.
You can confirm your password incrementally by 'matpw search' `,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := &create{}
		c.manager.Init()
		err := c.init()
		if err != nil {
			return err
		}

		input := bufio.NewScanner(os.Stdin)

		fmt.Print("input service(required): ")
		input.Scan()
		c.newpw.Service = input.Text()

		if c.newpw.Service == "" {
			return errors.New("service name should be filled")
		}

		for _, pw := range c.passwords {
			if c.newpw.Service == pw.Service {
				return errors.New("this service has been already registered")
			}
		}

		fmt.Print("input account(required): ")
		input.Scan()
		c.newpw.Account = input.Text()

		if c.newpw.Account == "" {
			return errors.New("service name should be filled")
		}

		fmt.Print("input descripiton: ")
		input.Scan()
		c.newpw.Descripiton = input.Text()

		c.newpw.UseUppercase = yesNo("use upper case? [")
		c.newpw.UseDigit = yesNo("use digit?")
		c.newpw.UseSymbol = yesNo("use use symbol?")

		return c.run(args)
	},
}

func (c *create) run(args []string) error {
	c.newpw.Create()
	c.meta.passwords = append(c.meta.passwords, c.newpw)
	c.meta.Save()
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

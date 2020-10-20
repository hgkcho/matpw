package cmd

import (
	"log"

	"github.com/manifoldco/promptui"
)

func yesNo(lavel string) bool {
	prompt := promptui.Select{
		Label: lavel + "[Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

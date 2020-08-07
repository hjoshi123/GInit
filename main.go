// Copyright 2020 Hemant Joshi. All rights reserved.
// Use of this source code is governed by MIT License.

// Package main deals with the calling of functions from package prompt and displaying
// error messages using chalk module. Command Line Arguments are not included yet
// and will be out soon.
package main

// TODO Create a new package for CLI which will take optional commands.

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/common-nighthawk/go-figure"
	"github.com/hjoshi123/ginit/cli"
	"github.com/hjoshi123/ginit/prompt"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	"github.com/ttacon/chalk"
)

var (
	result string
	err    error
)

func init() {
	handleControlC()

	if err := godotenv.Load(); err != nil {
		fmt.Print("No .env file found")
	}

	gitPat, exists := os.LookupEnv("GITHUB_PAT")

	validatePat := func(input string) error {
		if len(input) < 40 {
			return errors.New("Personal Access Token should be 40 characters")
		}
		return nil
	}
	if exists {
		result = gitPat
	} else {
		pat := promptui.Prompt{
			Label:     "Paste your Personal Access Token here",
			Templates: prompt.PromptTemplate,
			Validate:  validatePat,
			Mask:      ' ',
		}

		result, err = pat.Run()
		if err != nil {
			fmt.Println(chalk.Red.NewStyle().WithBackground(chalk.White).WithTextStyle(chalk.Bold).Style(err.Error()))
			os.Exit(0)
		}

		write, _ := godotenv.Unmarshal(fmt.Sprint("GITHUB_PAT=" + result))
		err = godotenv.Write(write, "./.env")
	}
}

func main() {
	handleControlC()

	myFigure := figure.NewColorFigure("GInit", "", "yellow", true)
	myFigure.Print()

	fmt.Println()

	app := cli.GetCli(result)

	err = app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func handleControlC() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Run Cleanup
		os.Exit(1)
	}()
}

// Copyright 2020 Hemant Joshi. All rights reserved.
// Use of this source code is governed by MIT License.

// Package main deals with the calling of functions from package prompt and displaying
// error messages using chalk module. Command Line Arguments are not included yet
// and will be out soon.
package main

import (
	"errors"
	"fmt"
	"html"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/common-nighthawk/go-figure"
	"github.com/hjoshi123/ginit/prompt"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	"github.com/ttacon/chalk"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Print("No .env file found")
	}
}

func main() {
	handleControlC()

	myFigure := figure.NewColorFigure("GInit", "", "yellow", true)
	myFigure.Print()

	fmt.Println()
	fmt.Println(chalk.Blue.NewStyle().WithTextStyle(chalk.Bold).Style("Please Generate Your GitHub personal Tokens with Repo scope. "))
	fmt.Println(chalk.Yellow.NewStyle().WithTextStyle(chalk.Italic).Style("https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token"))
	fmt.Println()

	gitPat, exists := os.LookupEnv("GITHUB_PAT")

	validatePat := func(input string) error {
		if len(input) < 40 {
			return errors.New("Personal Access Token should be 40 characters")
		}
		return nil
	}

	var result string
	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond) // Build our new spinner
	s.Suffix = fmt.Sprint(chalk.Yellow.NewStyle().Style(" Authenticating..."))
	s.Color("yellow", "bold") // Set the spinner color to a bold re
	s.FinalMSG = fmt.Sprint(chalk.Green.NewStyle().WithTextStyle(chalk.Bold).Style("Authenticated!!!"))

	if exists {
		result = gitPat
		s.Start()
	} else {
		pat := promptui.Prompt{
			Label:    "Paste your Personal Access Token here",
			Validate: validatePat,
			Mask:     ' ',
		}

		result, err := pat.Run()

		if err != nil {
			fmt.Println(chalk.Red.NewStyle().WithBackground(chalk.White).WithTextStyle(chalk.Bold).Style(err.Error()))
			return
		}
		s.Start()

		write, _ := godotenv.Unmarshal(fmt.Sprintf("GITHUB_PAT=%s", result))
		err = godotenv.Write(write, "./.env")
	}

	s.Stop()

	repoCreationStatus := prompt.RepoPrompt(result)

	if repoCreationStatus {
		// InitCommit process started
		commitStatus, err := prompt.InitCommit(result)

		if err != nil {
			fmt.Println(chalk.Red.NewStyle().WithTextStyle(chalk.Bold).Style(commitStatus + err.Error()))
			return
		}

		fmt.Println(chalk.Green.NewStyle().WithTextStyle(chalk.Bold).Style(commitStatus))

		pushStatus, err := prompt.PushCommit(result)

		if err != nil {
			fmt.Println(chalk.Red.NewStyle().WithTextStyle(chalk.Bold).Style(pushStatus + err.Error()))
			return
		}

		fmt.Println(chalk.Green.NewStyle().WithTextStyle(chalk.Bold).Style(pushStatus))

		fmt.Println()

		str := html.UnescapeString("&#" + strconv.Itoa(129395) + ";")
		fmt.Println(chalk.Yellow.NewStyle().WithTextStyle(chalk.Bold).Style("You're good to go.. Happy Hacking " + str))
		fmt.Println()
	} else {
		return
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

// Copyright 2020 Hemant Joshi. All rights reserved.
// Use of this source code is governed by MIT License.

// Package cli deals with the command line arguments and their behaviour.
// This package takes in arguments from OS and uses optional and mandatory flags to
// create a directory/repo and `git init`, create .gitignore and README and push it to the repo.
// This package requires the presence of <b>.env file or GITHUB_PAT</b> environmental variable.
package cli

import (
	"fmt"
	"html"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/hjoshi123/ginit/prompt"
	"github.com/ttacon/chalk"
	"github.com/urfave/cli/v2"
)

// GetCli returns an instance of cli application which is then used in main() to run the app
func GetCli(pat string) *cli.App {
	app := &cli.App{
		Name:  "ginit",
		Usage: "Use manual input to create and push to github",
		Action: func(c *cli.Context) error {
			promptMode(pat)
			return nil
		},
	}
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "usage",
			Aliases: []string{"u"},
			Usage:   "Usage of the app",
			Action: func(c *cli.Context) error {
				fmt.Println()
				fmt.Println(chalk.Blue.NewStyle().WithTextStyle(chalk.Bold).Style("Please Generate Your GitHub personal Tokens with Repo scope. "))
				fmt.Println(chalk.Yellow.NewStyle().WithTextStyle(chalk.Italic).Style("https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token"))
				fmt.Println()
				return nil
			},
		},
		{
			Name:    "repo",
			Aliases: []string{"r"},
			Usage:   "GitHub repo related operations",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "Mention the repository name",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "description",
					Aliases: []string{"d"},
					Usage:   "Description of the repo",
				},
				&cli.BoolFlag{
					Name:     "private",
					Aliases:  []string{"p"},
					Usage:    "Private Repo (T or F)",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				repoName := c.String("name")
				fmt.Println(repoName)

				repoDesc := c.String("description")
				repoStatusBool := c.Bool("private")

				s := spinner.New(spinner.CharSets[33], 100*time.Millisecond) // Build our new spinner
				s.Suffix = fmt.Sprint(chalk.Yellow.NewStyle().Style(" Creating the repo..."))
				s.Color("yellow", "bold") // Set the spinner color to a bold re
				s.FinalMSG = fmt.Sprint(chalk.Green.NewStyle().WithTextStyle(chalk.Bold).Style("Successfully created new repo"))
				s.Start()
				time.Sleep(100 * time.Millisecond)

				repo, profile, err := prompt.CreateRepo(pat, repoName, repoDesc, repoStatusBool)

				if err != nil {
					s.Stop()
					return cli.Exit(fmt.Sprintf(chalk.Red.NewStyle().WithTextStyle(chalk.Bold).Style("Repository not created")), 2)
				}

				s.Stop()

				commitStatus, err := prompt.InitCommit(pat, profile.GetLogin(), *repo.Name, profile.GetEmail())

				if err != nil {
					return cli.Exit(fmt.Sprintf(chalk.Red.NewStyle().WithTextStyle(chalk.Bold).Style("Commit failed")), 3)
				}
				fmt.Println(chalk.Green.NewStyle().WithTextStyle(chalk.Bold).Style(commitStatus))

				pushStatus, err := prompt.PushCommit(pat, profile.GetLogin(), *repo.CloneURL)

				if err != nil {
					return cli.Exit(fmt.Sprintf(chalk.Red.NewStyle().WithTextStyle(chalk.Bold).Style("Push failed")), 4)
				}

				fmt.Println(chalk.Green.NewStyle().WithTextStyle(chalk.Bold).Style(pushStatus))
				fmt.Println()
				str := html.UnescapeString("&#" + strconv.Itoa(129395) + ";")
				fmt.Println(chalk.Yellow.NewStyle().WithTextStyle(chalk.Bold).Style("You're good to go.. Happy Hacking " + str))
				fmt.Println()
				return nil
			},
		},
	}

	return app
}

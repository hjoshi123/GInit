// Copyright 2020 Hemant Joshi. All rights reserved.
// Use of this source code is governed by MIT License.

// Package prompt handles the prompt for creation of repository in GitHub.
// It takes in the GitHub personal access token to create your repo and then do
// an init commit with a README.md and .gitignore file.
// Currently only 6 tempaltes i.e. Node, Android, Golang, Java, Python, Rails.
// More templates will be added through Command Line arguments.
package prompt

// TODO Remove error msgs from Println and put erro msgs in a logger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"github.com/ttacon/chalk"
	"golang.org/x/oauth2"
)

var (
	profile  *github.User
	repoName string
	err      error
	repoURL  string
	r        *git.Repository
)

// RepoPrompt gives the prompt related to Repository
func RepoPrompt(pat string) bool {
	repoNamePromptValidation := func(input string) error {
		if len(input) < 1 {
			return errors.New("Repository Name should be greater than 1 character")
		}
		return nil
	}

	repoNamePrompt := promptui.Prompt{
		Label:     "Enter your repo name",
		Templates: promptTemplate,
		Validate:  repoNamePromptValidation,
	}

	repoName, err = repoNamePrompt.Run()

	if err != nil {
		fmt.Println(chalk.Red.NewStyle().WithBackground(chalk.White).WithTextStyle(chalk.Bold).Style(err.Error()))
		os.Exit(0)
	}

	repoDescPrompt := promptui.Prompt{
		Label:     "Optionally enter your repository description",
		Templates: promptTemplate,
	}

	repoDesc, err := repoDescPrompt.Run()
	if err != nil {
		fmt.Println(chalk.Red.NewStyle().WithBackground(chalk.White).WithTextStyle(chalk.Bold).Style(err.Error()))
		os.Exit(0)
	}

	repoStatusPrompt := promptui.Select{
		Label:     "Public or Private",
		Templates: selectTemplate,
		Items:     []string{"Public", "Private"},
	}

	_, repoStatus, err := repoStatusPrompt.Run()
	if err != nil {
		fmt.Println(chalk.Red.NewStyle().WithBackground(chalk.White).WithTextStyle(chalk.Bold).Style(err.Error()))
		os.Exit(0)
	}

	var repoStatusBool bool = false
	if repoStatus == "Private" {
		repoStatusBool = true
	}

	s := spinner.New(spinner.CharSets[33], 100*time.Millisecond) // Build our new spinner
	s.Suffix = fmt.Sprint(chalk.Yellow.NewStyle().Style(" Creating the repo..."))
	s.Color("yellow", "bold") // Set the spinner color to a bold re
	s.FinalMSG = fmt.Sprint(chalk.Green.NewStyle().WithTextStyle(chalk.Bold).Style("Successfully created new repo"))
	s.Start()
	time.Sleep(100 * time.Millisecond)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: pat})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	r := &github.Repository{Name: &repoName, Private: &repoStatusBool, Description: &repoDesc}
	repo, _, err := client.Repositories.Create(ctx, "", r)

	if err != nil {
		log.Output(1, err.Error())
		s.FinalMSG = fmt.Sprint(chalk.Red.NewStyle().WithTextStyle(chalk.Bold).Style("Problem creating the repo"))
		s.Stop()
		return false
	}

	user := repo.GetOwner()
	profile, _, err = client.Users.Get(ctx, *user.Login)

	repoURL = repo.GetCloneURL()

	s.Stop()

	return true
}

// InitCommit does the following things
// 1. Clone the remote repo and store it to the location
// 2. Get the gitignore file from GitHub if avaliable and write it to your the folder.
// 3. Init the commit with a commit msg
func InitCommit(pat string) (string, error) {
	username := profile.GetLogin()
	directory := repoName
	err = os.Mkdir(directory, 0755)

	gitIgnoreContent := GitIgnorePrompt()
	if gitIgnoreContent != "" {
		filePath := filepath.Join(directory, ".gitignore")
		file, err := os.Create(filePath)
		defer file.Close()

		gitIgnoreContent = strings.Replace(gitIgnoreContent, `\n`, "\n", -1)

		_, err = file.Write([]byte(gitIgnoreContent))

		if err != nil {
			log.Println(err)
		}
	}

	readmePath := filepath.Join(directory, "README.md")
	file, err := os.Create(readmePath)
	defer file.Close()

	_, err = file.Write([]byte("# " + repoName))

	r, err = git.PlainInit(repoName, false)
	w, err := r.Worktree()

	_, err = w.Add(".gitignore")
	_, err = w.Add("README.md")

	if err != nil {
		return "Not able to git add", err
	}

	commit, err := w.Commit("[Init]: Initalised Repostiory", &git.CommitOptions{
		Author: &object.Signature{
			Name:  username,
			Email: profile.GetEmail(),
			When:  time.Now(),
		},
	})

	if err != nil {
		return "Commit error", err
	}

	return ("Commit Successful " + commit.String()), nil
}

// PushCommit pushses the InitCommit
func PushCommit(pat string) (string, error) {
	username := profile.GetLogin()

	s := spinner.New(spinner.CharSets[8], 100*time.Millisecond)
	s.Suffix = fmt.Sprint(chalk.Yellow.NewStyle().Style(" Pushing the commit..."))
	s.Color("yellow", "bold") // Set the spinner color to a bold re
	s.Start()

	_, err := r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})

	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: username,
			Password: pat,
		},
	})

	if err != nil {
		s.Stop()
		return "Push Failed", err
	}

	s.Stop()
	return "Pushed successfully", nil
}

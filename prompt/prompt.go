package prompt

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
		Label:    "Enter your repo name",
		Validate: repoNamePromptValidation,
	}

	repoName, err = repoNamePrompt.Run()

	if err != nil {
		fmt.Println(chalk.Red.NewStyle().WithBackground(chalk.White).WithTextStyle(chalk.Bold).Style(err.Error()))
		os.Exit(0)
	}

	repoDescPrompt := promptui.Prompt{
		Label: "Optionally enter your repository description",
	}

	repoDesc, _ := repoDescPrompt.Run()

	repoStatusPrompt := promptui.Select{
		Label: "Public or Private",
		Items: []string{"Public", "Private"},
	}

	_, repoStatus, _ := repoStatusPrompt.Run()

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
	err := os.Mkdir(directory, 0755)

	gitIgnoreContent := GitIgnoreContent()
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

	r, err = git.PlainInit(repoName, false)
	w, err := r.Worktree()

	_, err = w.Add(".gitignore")
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

package prompt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"github.com/ttacon/chalk"
	"golang.org/x/oauth2"
)

type gitIgnoreResponse struct {
	Name   string
	Source string
}

var (
	// PromptTemplate which displays green tick when the input is valid and red text when it
	// is invalid. `.` indicates the text to be displayed.
	PromptTemplate = &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: fmt.Sprintf("%s {{ . }} ", promptui.Styler(promptui.FGGreen)("✔")),
	}

	// selectTemplate displays white tick and blurs the other options when selected it shows
	// green tick and displays it in bold.
	selectTemplate = &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   fmt.Sprintf("%s {{ . }} ", promptui.Styler(promptui.FGWhite)("✔")),
		Inactive: "{{ . | cyan }}",
		Selected: fmt.Sprintf("%s {{ . | bold }} ", promptui.Styler(promptui.FGGreen)("✔")),
	}
)

// GitIgnorePrompt prompts the user to select the gitignore template and then calls
// gitIgnoreContent to return the content of the file
func GitIgnorePrompt() string {
	gitIgnorePrompt := promptui.Select{
		Label:     "Select your type of project",
		Templates: selectTemplate,
		Items:     []string{"Node", "Android", "Java", "Python", "Go", "Rails", "None"},
	}

	_, gitIgnore, err := gitIgnorePrompt.Run()
	if err != nil {
		fmt.Println(chalk.Red.NewStyle().WithBackground(chalk.White).WithTextStyle(chalk.Bold).Style(err.Error()))
		os.Exit(0)
	}

	return gitIgnoreContent(gitIgnore)
}

// CreateRepo creates a repo by taking in parameters and returning the repository
// and the user's profile of GitHub
func CreateRepo(pat, repoName, repoDesc string, privateRepo bool) (*github.Repository, *github.User, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: pat})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	r := &github.Repository{Name: &repoName, Private: &privateRepo, Description: &repoDesc}
	repo, _, err := client.Repositories.Create(ctx, "", r)

	user := repo.GetOwner()
	profile, _, err := client.Users.Get(ctx, *user.Login)

	if err != nil {
		log.Output(1, err.Error())
		return nil, nil, err
	}

	return repo, profile, nil
}

// GitIgnoreContent represents the .gitignore file template.
// Returns the template from GitHub's available list.
func gitIgnoreContent(gitIgnore string) string {
	if gitIgnore != "None" {
		resp, err := http.Get(fmt.Sprintf("https://api.github.com/gitignore/templates/%s", gitIgnore))

		if err != nil {
			fmt.Println("Error caused " + err.Error())
			return ""
		}

		defer resp.Body.Close()

		gitIgnoreResp := new(gitIgnoreResponse)
		err = json.NewDecoder(resp.Body).Decode(&gitIgnoreResp)

		if err != nil {
			fmt.Println(err)
			return ""
		}

		return gitIgnoreResp.Source
	}

	return ""
}

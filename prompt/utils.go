package prompt

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/manifoldco/promptui"
)

type gitIgnoreResponse struct {
	Name   string
	Source string
}

// GitIgnorePrompt prompts the user to select the gitignore template and then calls
// gitIgnoreContent to return the content of the file
func GitIgnorePrompt() string {
	gitIgnorePrompt := promptui.Select{
		Label: "Select your type of project",
		Items: []string{"Node", "Android", "Java", "Python", "Go", "Rails", "None"},
	}

	_, gitIgnore, _ := gitIgnorePrompt.Run()

	return gitIgnoreContent(gitIgnore)
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

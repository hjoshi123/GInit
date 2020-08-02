package prompt

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/manifoldco/promptui"
)

// GitIgnoreContent represents the .gitignore file template.
// Returns the template from GitHub's available list.
func GitIgnoreContent() string {
	gitIgnorePrompt := promptui.Select{
		Label: "Select your type of project",
		Items: []string{"Node", "Android", "Java", "Python", "Go", "Rails", "None"},
	}

	_, gitIgnore, _ := gitIgnorePrompt.Run()

	if gitIgnore != "None" {
		resp, err := http.Get(fmt.Sprintf("https://api.github.com/gitignore/templates/%s", gitIgnore))

		if err != nil {
			fmt.Println("Error caused " + err.Error())
			return ""
		}

		defer resp.Body.Close()

		data, _ := ioutil.ReadAll(resp.Body)

		return string(data)
	}

	return ""
}

package cli

import (
	"fmt"
	"html"
	"strconv"

	"github.com/hjoshi123/ginit/prompt"
	"github.com/ttacon/chalk"
)

func promptMode(result string) {

	fmt.Println()
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

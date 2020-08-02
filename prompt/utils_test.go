package prompt

import "testing"

func TestGitIgnoreContent(t *testing.T) {
	output := gitIgnoreContent("Go")

	if output == "" {
		t.Errorf("Error in fetching results")
	}
}

package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIssueTracker_InjectLinks(t *testing.T) {
	cases := []struct {
		Title    string
		URL      string
		Body     string
		Expected string
	}{
		{
			Title:    "one link",
			URL:      "https://jira.com/board",
			Body:     "task ABC-12",
			Expected: "task [ABC-12](https://jira.com/board/ABC-12)",
		},
	}

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			tracker := NewIssueTracker(c.URL)

			assert.Equal(t, c.Expected, tracker.InjectLinks(c.Body))
		})
	}
}

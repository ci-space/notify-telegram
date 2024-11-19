package internal

import (
	"fmt"
	"regexp"
	"strings"
)

type IssueTracker struct {
	url string
}

func NewIssueTracker(url string) *IssueTracker {
	return &IssueTracker{url: strings.TrimRight(url, "/")}
}

var issueTrackerRegex = regexp.MustCompile(`(?m)([A-Z]+-[0-9]+)`)

func (t *IssueTracker) InjectLinks(body string) string {
	return issueTrackerRegex.ReplaceAllString(body, fmt.Sprintf("[$1](%s/$1)", t.url))
}

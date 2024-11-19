package internal

import (
	"fmt"
	"testing"
)

func TestMessenger_createTgMessage(t *testing.T) {
	m := &Messenger{
		markdownConverter: NewMarkdownRenderer(),
	}

	got := m.createTgMessage(Message{
		Body:            "## Heading 1\n- 1\n- 2",
		ConvertMarkdown: true,
	})

	fmt.Println(got.Text)
}

# tgapi

```
go get github.com/ci-space/notify-telegram/pkg/tgapi@v0.1.0
```

**Usage**

```go
package main

import (
	"context"
	"github.com/ci-space/notify-telegram/pkg/tgapi"
)

func main() {
	client := tgapi.NewClient("<bot-token>", "api.telegram.org")

	client.SendMessage(context.Background(), tgapi.SendingMessage{
		Body: "message",
		ChatID: "chat id",
    })
}
```

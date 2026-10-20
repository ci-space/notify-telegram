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
	client, err := tgapi.NewClient(tgapi.Config{
		BotToken: "<bot-token>",
		Host:     "api.telegram.org",
	})
	if err != nil {
		panic(err)
	}

	client.SendMessage(context.Background(), tgapi.SendingMessage{
		Body: "message",
		ChatID: "chat id",
    })
}
```

To send requests through a SOCKS5 proxy, set its address in the config:

```go
client, err := tgapi.NewClient(tgapi.Config{
	BotToken:   "<bot-token>",
	Host:       "api.telegram.org",
	Socks5Addr: "127.0.0.1:1080",
})
if err != nil {
	panic(err)
}
```

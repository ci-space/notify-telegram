# notify-telegram

notify-telegram - GitHub Action for sending messages to Telegram

📦 [md2html](./pkg/md2html) 📦 [tgapi](./pkg/tgapi)


**Features**
- Convert Markdown to Telegram HTML
- Inject links to issue trackers

**Inputs**

| Parameter             | Description                                                                                             | Required | Default value      |
|-----------------------|---------------------------------------------------------------------------------------------------------|----------|--------------------|
| **token**             | Token for Telegram bot                                                                                  | ✔️       | -                  |
| **chat_id**           | ID of Telegram chat                                                                                     | ✔️       | -                  |
| **chat_thread_id**    | ID of Telegram chat thread                                                                              |          | -                  |
| **message**           | Text for message                                                                                        | ✔️       | -                  |
| **host**              | Telegram host                                                                                           |          | `api.telegram.org` |
| **convert_markdown**  | Flag that indicates whether markdown should be converted from **message**. Possible values: true, false |          | `false`            |
| **issue_tracker_url** | URL to issue tracker. Example: https://my-project.atlassian.net/browse                                  |          | -                  |
| **mode**              | Mode of send. Possible values: `create`, `update`                                                       |          | `create`           |
| **message_id**        | ID of updating message. Required if mode is `update`                                                    |          | `update`           |

**Outputs**

| Parameter      | Description        |
|----------------|--------------------|
| **message_id** | ID of sent message |

## Usage

### Send message on Release

.github/workflows/release.yaml:

```yaml
name: Release

on:
  release:
    types:
      - published

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - name: Send notification
        uses: ci-space/notify-telegram@v0.3.2
        with:
          token: ${{ secrets.TELEGRAM_TOKEN }}
          chat_id: ${{ secrets.TELEGRAM_CHAT }}
          chat_thread_id: ${{ secrets.TELEGRAM_CHAT_THREAD }}
          convert_markdown: true
          message: |
            ${{ github.repository }} deployed on tag ${{ github.event.release.tag_name }}
```

### Get ID of sent message

.github/workflows/release.yaml:

```yaml
name: Release

on:
  release:
    types:
      - published

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - name: Send notification
        id: send
        uses: ci-space/notify-telegram@v0.3.2
        with:
          # omit config

      - name: Print message id
        run: echo ${{ steps.send.outputs.message_id }}
```


### Create and update message

.github/workflows/release.yaml:

```yaml
name: Release

on:
  release:
    types:
      - published

jobs:
  counter:
    runs-on: ubuntu-latest

    steps:
      - name: Send count 1
        id: first
        uses: ci-space/notify-telegram@v0.3.2
        with:
          token: ${{ secrets.TELEGRAM_TOKEN }}
          chat_id: ${{ secrets.TELEGRAM_CHAT }}
          message: |
            Count: 1
            
      - name: Send count 2
        uses: ci-space/notify-telegram@v0.3.2
        with:
          token: ${{ secrets.TELEGRAM_TOKEN }}
          chat_id: ${{ secrets.TELEGRAM_CHAT }}
          message_id: ${{ steps.first.outputs.message_id }}
          mode: update
          message: |
            Count: 2
```


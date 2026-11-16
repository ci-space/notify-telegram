# notify-telegram

notify-telegram - GitHub Action for sending messages to Telegram

📦 [md2html](./pkg/md2html) 📦 [tgapi](./pkg/tgapi)


**Features**
- Convert Markdown to Telegram HTML
- Inject links to issue trackers
- Attach an inline keyboard with URL buttons

**Inputs**

| Parameter             | Description                                                                                             | Required | Default value      |
|-----------------------|---------------------------------------------------------------------------------------------------------|----------|--------------------|
| **token**             | Token for Telegram bot                                                                                  | ✔️       | -                  |
| **chat_id**           | ID of Telegram chat                                                                                     | ✔️       | -                  |
| **chat_thread_id**    | ID of Telegram chat thread                                                                              |          | -                  |
| **message**           | Text for message                                                                                        | ✔️       | -                  |
| **buttons**           | Inline keyboard as a YAML list of rows containing URL buttons                                           |          | -                  |
| **host**              | Telegram host                                                                                           |          | `api.telegram.org` |
| **socks5_addr**       | SOCKS5 proxy address. Example: `127.0.0.1:1080`                                                         |          | -                  |
| **convert_markdown**  | Flag that indicates whether markdown should be converted from **message**. Possible values: true, false |          | `false`            |
| **issue_tracker_url** | URL to issue tracker. Example: https://my-project.atlassian.net/browse                                  |          | -                  |
| **mode**              | Mode of send. Possible values: `create`, `update`                                                       |          | `create`           |
| **message_id**        | ID of updating message. Required if mode is `update`                                                    |          | -                  |

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

### Send a message with URL buttons

Use an action release that includes inline keyboard support in its Docker image.
Replace `RELEASE_TAG` with that release tag.

```yaml
name: Notify Telegram

on:
  release:
    types:
      - published

jobs:
  notify:
    runs-on: ubuntu-latest
    steps:
      - name: Send release notification
        uses: ci-space/notify-telegram@RELEASE_TAG
        with:
          token: ${{ secrets.TELEGRAM_TOKEN }}
          chat_id: ${{ secrets.TELEGRAM_CHAT }}
          chat_thread_id: ${{ secrets.TELEGRAM_CHAT_THREAD }}
          message: Release ${{ github.event.release.tag_name }} is available.
          buttons: |
            - - text: Open release
                url: ${{ github.event.release.html_url }}
              - text: Workflow logs
                url: ${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}
            - - text: All releases
                url: ${{ github.server_url }}/${{ github.repository }}/releases
```

The outer YAML list defines the rows. Each inner list defines the buttons in
one row. In the example, the first row has two buttons and the second has one.
Each button requires non-empty `text` and `url` fields. Use an absolute
`http://`, `https://`, or `tg://` URL.

You can use `buttons` in both `create` and `update` mode. Omit the input or
leave it empty to send the request without keyboard markup.

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

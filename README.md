# telegram-action

telegram-action - GitHub Action for sending messages to Telegram

Features:
- Convert Markdown to Telegram HTML
- Inject links to issue trackers

## Usage

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
        name: Send notification
        uses: ci-space/notify-telegram@v0.3.0
        with:
          token: ${{ secrets.TELEGRAM_TOKEN }}
          chat_id: ${{ secrets.TELEGRAM_CHAT }}
          chat_thread_id: ${{ secrets.TELEGRAM_CHAT_THREAD }}
          convert_markdown: true
          message: |
            ${{ github.repository }} deployed on tag ${{ github.event.release.tag_name }}
```
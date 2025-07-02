module github.com/ci-space/notify-telegram/pkg/tgapi

go 1.23.3

require github.com/ci-space/notify-telegram/pkg/md2html v0.1.0

require github.com/gomarkdown/markdown v0.0.0-20250311123330-531bef5e742b // indirect

replace github.com/ci-space/notify-telegram/pkg/md2html => ./../md2html

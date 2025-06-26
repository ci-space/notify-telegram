module github.com/ci-space/notify-telegram

go 1.23.3

replace github.com/ci-space/notify-telegram/pkg/md2html => ./pkg/md2html

require (
	github.com/artarts36/singlecli v0.0.0-20241017172045-f9a31a534745
	github.com/caarlos0/env/v11 v11.2.2
	github.com/ci-space/notify-telegram/pkg/md2html v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gomarkdown/markdown v0.0.0-20250311123330-531bef5e742b // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

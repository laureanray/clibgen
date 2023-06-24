module github.com/laureanray/clibgen

go 1.18

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/fatih/color v1.13.0
	github.com/kennygrant/sanitize v1.2.4
	github.com/manifoldco/promptui v0.9.0
	github.com/schollz/progressbar/v3 v3.8.6
	github.com/spf13/cobra v1.4.0
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.0.0-20220131195533-30dcbda58838 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
)

replace github.com/laureanray/clibgen/internal/book => ./internal/book
replace github.com/laureanray/clibgen/internal/libgen => ./internal/libgen
replace github.com/laureanray/clibgen/internal/console => ./internal/console
replace github.com/laureanray/clibgen/internal/mirror => ./internal/mirror
replace github.com/laureanray/clibgen/internal/page => ./internal/page
replace github.com/laureanray/clibgen/internal/downloader => ./internal/downloader


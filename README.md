# clibgen
[![Go Report Card](https://goreportcard.com/badge/github.com/laureanray/clibgen)](https://goreportcard.com/report/github.com/laureanray/clibgen)  
clibgen is a TUI for library genesis. (https://libgen.li, https://libgen.is)

![ezgif com-gif-maker (7)](https://user-images.githubusercontent.com/22195710/180980454-4e0c95b5-1df3-4891-84f0-9b92d0ac12d4.gif)

### Installation

Shell Installation (Mac, Linux)

```shell
curl -fsSL https://raw.githubusercontent.com/laureanray/clibgen/main/install.sh | sh
```

PowerShell Installation (Windows)

```powershell
iwr https://raw.githubusercontent.com/laureanray/clibgen/main/install.ps1 -useb | iex
```

Manual Installation (Binaries)

- Download the file for your platform and version you need in the [releases page](https://github.com/laureanray/clibgen/releases)
- Extract to where you store your binaries
- Add the binaries directory to your PATH

### Usage

Search for a book using the old website (default)

```shell
clibgen search "Eloquent JavaScript"
```

#### Search

Usage:
`  clibgen search [flags]`

Flags:
-  -f, --filter string           search by [title, author, isbn] (default "title")
-  -h, --help                    help for search
-  -n, --number of results int   number of result(s) to be displayed maximum: 25 (default 10)
-  -o, --output string           Output directory (default "./")
-  -s, --site string             which website to use [legacy, new] (default "legacy")



### Found an issue?

Please open up an issue or a PR if you have a fix for it. 

# clibgen

clibgen is a TUI for library genesis. (https://libgen.li, https://libgen.is)

![ezgif-3-2ec26c641b](https://user-images.githubusercontent.com/22195710/176466306-0dd493dd-5a3b-494a-96c5-e2380b830275.gif)

### Install

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

Search for a book using old the old website (default)
```shell
clibgen search "Eloquent JavaScript"
```

Search for a book using the newer website (useful if for some reason the old website is down or the mirrors are not working)
`-s or -site` flag
```shell
clibgen search -m "new" "Eloquent JavaScript"
```

Limit search results (default: 10)
```shell
clibgen search -n 5 "Eloquent JavaScript"
```


### Found an issue?
Please open up an issue or a PR if you have a fix for it ;)

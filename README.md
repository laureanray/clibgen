# clibgen

clibgen is a TUI for library genesis. (https://libgen.li, https://libgen.is)

![ezgif-3-2ec26c641b](https://user-images.githubusercontent.com/22195710/176466306-0dd493dd-5a3b-494a-96c5-e2380b830275.gif)

### Install

Shell Installation (Mac, Linux)

```shell
curl -fsSL https://raw.githubusercontent.com/laureanray/clibgen/main/install.sh | sh
```

PowerShell Installation (Work in progress)

### Usage 

Search for a book using old the old website (default)
```shell
clibgen search "Eloquent JavaScript"
```

Search for a book using the newer website (useful if for some reason the old website is down or the mirrors are not working)
`-m or -mirror` flag
```shell
clibgen search -m "new" "Eloquent JavaScript"
```

Limit search results
```shell
clibgen search -n 5 "Eloquent JavaScript"
```


### Found an issue?
Please open up an issue or a PR if you have a fix for it ;)

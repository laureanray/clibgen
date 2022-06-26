#!/bin/sh
# Copyright 2019 the Deno authors. All rights reserved. MIT license.
# TODO(everyone): Keep this script simple and easily auditable.

set -e

if ! command -v unzip >/dev/null; then
	echo "Error: unzip is required to install clibgen" 1>&2
	exit 1
fi

if [ "$OS" = "Windows_NT" ]; then
	target="Windows_x86_64"
else
	case $(uname -sm) in
	"Darwin x86_64") target="Darwin_x86_64" ;;
	"Darwin arm64") target="Darwin_arm64" ;;
	*) target="Linux_x86_64" ;;
	esac
fi

if [ $# -eq 0 ]; then
	clibgen_uri="https://github.com/laureanray/clibgen/releases/latest/download/clibgen_${target}.zip"
else
	clibgen_uri="https://github.com/laureanray/clibgen/releases/download/${1}/clibgen_${target}.zip"
fi

clibgen_install="${CLIBGEN_INSTALL:-$HOME/.clibgen}"
bin_dir="$clibgen_install/bin"
exe="$bin_dir/deno"

if [ ! -d "$bin_dir" ]; then
	mkdir -p "$bin_dir"
fi

echo "$clibgen_uri"
curl --fail --location --progress-bar --output "$exe.zip" "$clibgen_uri"
unzip -d "$bin_dir" -o "$exe.zip"
chmod +x "$exe"
rm "$exe.zip"

echo "Clibgen was installed successfully to $exe"
if command -v deno >/dev/null; then
	echo "Run 'clibgen --help' to get started"
else
	case $SHELL in
	/bin/zsh) shell_profile=".zshrc" ;;
	*) shell_profile=".bashrc" ;;
	esac
	echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
	echo "  export CLIBGEN_INSTALL=\"$clibgen_install:was\""
	echo "  export PATH=\"\$CLIBGEN_INSTALL/bin:\$PATH\""
	echo "Run '$exe --help' to get started"
fi
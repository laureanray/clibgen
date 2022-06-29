#!/bin/sh

set -e

if ! command -v tar >/dev/null; then
	echo "Error: tar is required to install clibgen" 1>&2
	exit 1
fi

if [ "$OS" = "Windows_NT" ]; then
	target="Windows_x86_64"
else
	case $(uname -sm) in
	"Darwin x86_64") target="Darwin-x86_64" ;;
	"Darwin arm64") target="Darwin-arm64" ;;
	*) target="Linux_x86_64" ;;
	esac
fi

if [ $# -eq 0 ]; then
	clibgen_uri="https://github.com/laureanray/clibgen/releases/latest/download/clibgen_${target}.tar.gz"
else
	clibgen_uri="https://github.com/laureanray/clibgen/releases/download/${1}/clibgen_${target}.tar.gz"
fi

clibgen_install="${CLIBGEN_INSTALL:-$HOME/.clibgen}"
bin_dir="$clibgen_install/bin"
exe="$bin_dir/clibgen"

if [ ! -d "$bin_dir" ]; then
	mkdir -p "$bin_dir"
fi

echo "$clibgen_uri"
curl --fail --location --progress-bar --output "$exe.tar.gz" "$clibgen_uri"
tar -xf "$exe.tar.gz" -C "$bin_dir"
chmod +x "$exe"
# Clean up
rm "$exe.tar.gz"

echo "Clibgen was installed successfully to $exe"
if command -v clibgen >/dev/null; then
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
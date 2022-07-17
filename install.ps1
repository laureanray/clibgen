#!/usr/bin/env pwsh
# Copyright 2022 Laurean Ray Bahala. All rights reserved. MIT license.
# TODO(everyone): Keep this script simple and easily auditable.

$ErrorActionPreference = 'Stop'

if ($v) {
  $Version = "v${v}"
}
if ($args.Length -eq 1) {
  $Version = $args.Get(0)
}

$ClibgenInstall = $env:CLIBGEN_INSTALL
$BinDir = if ($ClibgenInstall) {
  "$ClibgenInstall\bin"
} else {
  "$Home\.clibgen\bin"
}

$ClibgenZip = "$BinDir\clibgen.zip"
$ClibgenExe = "$BinDir\clibgen.exe"
$Target = 'Windows-x86_64'

# GitHub requires TLS 1.2
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$ClibgenUri = if (!$Version) {
  "https://github.com/laureanray/clibgen/releases/latest/download/clibgen_${Target}.zip"
} else {
  "https://github.com/laureanray/clibgen/releases/download/${Version}/clibgen_${Target}.zip"
}

Write-Output $ClibgenUri

if (!(Test-Path $BinDir)) {
  New-Item $BinDir -ItemType Directory | Out-Null
}

curl.exe -Lo $ClibgenZip $ClibgenUri

tar.exe xf $ClibgenZip -C $BinDir

Remove-Item $ClibgenZip

$User = [EnvironmentVariableTarget]::User
$Path = [Environment]::GetEnvironmentVariable('Path', $User)
if (!(";$Path;".ToLower() -like "*;$BinDir;*".ToLower())) {
  [Environment]::SetEnvironmentVariable('Path', "$Path;$BinDir", $User)
  $Env:Path += ";$BinDir"
}

Write-Output "Clibgen was installed successfully to $ClibgenExe"
Write-Output "Run 'clibgen --help' to get started"
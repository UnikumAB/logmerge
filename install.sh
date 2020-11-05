#!/usr/bin/env bash
set -e

wget -O - https://api.github.com/repos/UnikumAB/logmerge/releases/latest |\
  jq -r '.assets[]|select(.name|contains("'$(uname -s)'"))|select(.name|contains("'$(uname -m)'"))|.browser_download_url'|\
  wget -i - -O logmerge.tar.gz

tar xfvz logmerge.tar.gz

install -p logmerge /usr/local/bin

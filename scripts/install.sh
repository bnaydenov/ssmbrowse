#!/bin/sh
# Usage: [sudo] [BINDIR=/usr/local/bin] ./install.sh [<BINDIR>]
#
# Example:
#     1. sudo ./install.sh /usr/local/bin
#     2. sudo ./install.sh /usr/bin
#     3. ./install.sh $HOME/usr/bin
#     4. BINDIR=$HOME/usr/bin ./install.sh
#
# Default BINDIR=/usr/local/bin

set -euf

if [ -n "${DEBUG-}" ]; then
    set -x
fi

: "${BINDIR:=/usr/local/bin}"

if [ $# -gt 0 ]; then
  BINDIR=$1
fi

_can_install() {
  if [ ! -d "${BINDIR}" ]; then
    mkdir -p "${BINDIR}" 2> /dev/null
  fi
  [ -d "${BINDIR}" ] && [ -w "${BINDIR}" ]
}

if ! _can_install && [ "$(id -u)" != 0 ]; then
  printf "Run script as sudo\n"
  exit 1
fi

if ! _can_install; then
  printf -- "Can't install to %s\n" "${BINDIR}"
  exit 1
fi

# machine=$(uname -m)

case $(uname -m) in
    x86_64)
       machine="x86_64"
        ;;
    aarch64)
        machine="arm64"
        ;;
    *)
        printf "Arch not supported\n"
        exit 1
        ;;
esac

case $(uname -s) in
    Linux)
        os="linux"
        ;;
    Darwin)
        os="macos"
        ;;
    *)
        printf "OS not supported\n"
        exit 1
        ;;
esac

# exit if no unzip command found 
if ! [ -x "$(command -v unzip)" ]; then
    echo "unzip command could not be found, instalation can not proceed. "
    exit 1
fi

printf "Fetching latest stable version\n"
latest="$(curl -sL 'https://api.github.com/repos/bnaydenov/ssmbrowse/releases/latest' | grep 'tag_name' | grep --only 'v[0-9\.]\+' | cut -c 2- | awk -F  '-' '/1/ {print $1}')"
tempFolder="/tmp/ssmbrowse-${latest}"

printf -- "Found version %s\n" "${latest}"
 
mkdir -p "${tempFolder}" 2> /dev/null
printf -- "Downloading ssmbrowse_%s_%s_%s.zip\n" "${latest}" "${os}" "${machine}"

curl -sL --output ${tempFolder}/ssmbrowse.zip "https://github.com/bnaydenov/ssmbrowse/releases/download/v${latest}/ssmbrowse_${latest}_${os}_${machine}.zip"
unzip -d ${tempFolder} ${tempFolder}/ssmbrowse.zip


printf -- "Installing ssmbrowse into ${BINDIR}\n"
install -m755 "${tempFolder}/ssmbrowse" "${BINDIR}/ssmbrowse"

printf "Cleaning up temp files ${tempFolder}\n"
rm -rf "${tempFolder}"

printf -- "Successfully installed ssmbrowse into %s/\n" "${BINDIR}"

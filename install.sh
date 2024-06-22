#!/bin/bash

if [ "$(whoami)" != "root" ]; then
	echo "please run this script as root"
	exit 1
fi

if command -v go &>/dev/null; then
	echo "Go version found"
	go version
else
	echo "golang must be installed"
	exit 1
fi

mkdir -pv /etc/some/ /var/some/

cp -v config.yml /etc/some/
go mod init "some-pkgmgr" && go mod tidy
go build && cp -v some-pkgmgr /usr/bin/some

echo
echo "----some-pkgmgr installed----"
echo

echo "Usage: some -d [install|update|remove] [package name]"

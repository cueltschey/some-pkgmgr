# Some Package manager
A simple package manager that works with multiple package formats

## Dependencies:

General:
* golang
* binutils (ar)

Go Dependencies:
* github.com/mattn/go-sqlite3 v1.14.22
* github.com/ulikunitz/xz v0.5.12
* gopkg.in/yaml.v3 v3.0.1

## Installation:

Install with script:
```
git clone https://github.com/cueltschey/some-pkgmgr
cd some-pkgmgr
sudo install.sh
```
Or:
```
git clone https://github.com/cueltschey/some-pkgmgr
cd some-pkgmgr
go mod init "some-pkgmgr"
go build
mkdir /etc/some && sudo cp -v config.yml /etc/some
sudo cp -v some-pkgmgr /usr/bin/some
```

## Configuration:

Paths used by some package manager are located in config.yml

- **deb-uri** -> base uri for debian packages
- **keyring** -> your system gpg keyring for ubuntu
- **tmpdir** -> directory for extracting temporary files
- **dbpath** -> where to put the package database



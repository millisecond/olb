#!/bin/bash -e

readonly prgdir=$(cd $(dirname $0); pwd)
readonly basedir=$(cd $prgdir/..; pwd)
v=$1

[[ -n "$v" ]] || read -p "Enter version (e.g. 1.0.4): " v
if [[ -z "$v" ]] ; then
	echo "Usage: $0 [<version>] (e.g. 1.0.4)"
	exit 1
fi

go get -u github.com/mitchellh/gox
for go in go1.8.3; do
	echo "Building olb with ${go}"
	gox -gocmd ~/${go}/bin/go -tags netgo -output "${basedir}/build/builds/olb-${v}/olb-${v}-${go}-{{.OS}}_{{.Arch}}"
done

( cd ${basedir}/build/builds/olb-${v} && shasum -a 256 olb-${v}-* > olb-${v}.sha256 )
( cd ${basedir}/build/builds/olb-${v} && gpg2 --output olb-${v}.sha256.sig --detach-sig olb-${v}.sha256 )

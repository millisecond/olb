#!/bin/bash
#
# homebrew.sh creates an updated homebrew on olb

set -o nounset
set -o errexit

readonly prgdir=$(cd $(dirname $0); pwd)
readonly brewdir=$(brew --prefix)/Homebrew/Library/Taps/homebrew/homebrew-core

v=${1:-}
[[ -n "$v" ]] || read -p "Enter version (e.g. 1.0.4): " v
if [[ -z "$v" ]] ; then
	echo "Usage: $0 <version> (e.g. 1.0.4)"
	exit 1
fi

srcurl=https://github.com/millisecond/olb/archive/v${v}.tar.gz
shasum=$(wget -O- -q "$srcurl" | shasum -a 256 | awk '{ print $1; }')
echo -e "/url
DAurl \"$srcurl\"/sha256
DAsha256 \"$shasum\":wq" > $prgdir/homebrew.vim

brew update
brew update
(
	cd $brewdir
	git checkout -b olb-$v origin/master
	vim -u NONE -s $prgdir/homebrew.vim $brewdir/Formula/olb.rb
	git add Formula/olb.rb
	git commit -m "olb $v"
	git push --set-upstream magiconair olb-$v
)

echo "Goto https://github.com/Homebrew/homebrew-core to create pull request"
open https://github.com/Homebrew/homebrew-core

exit 0

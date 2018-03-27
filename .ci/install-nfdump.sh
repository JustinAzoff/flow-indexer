#!/bin/sh

set -ex

VER=1.6.16
DEST=$HOME/.sources/nfdump-${VER}.tar.gz

mkdir -p $HOME/.sources

if ! [ -e $DEST ] ; then
	curl -Lo $DEST https://github.com/phaag/nfdump/archive/v${VER}.tar.gz
fi

tar xvzf $DEST
cd nfdump-$VER
./configure --prefix=$HOME/local && make && make install
$HOME/local/bin/nfdump -V

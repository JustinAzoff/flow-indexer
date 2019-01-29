#!/bin/sh

set -ex

VER=3.0.8.2
DEST=$HOME/.sources/argus-clients_${VER}.tar.gz

mkdir -p $HOME/.sources


if ! [ -e $DEST ] ; then
    curl -Lo $DEST http://deb.debian.org/debian/pool/main/a/argus-clients/argus-clients_${VER}.orig.tar.gz
fi

tar xvzf $DEST
cd argus*/
./configure --prefix=$HOME/local && make && make install

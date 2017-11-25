#!/bin/sh

# From http://hyperledger-fabric.readthedocs.io/en/latest/samples.html

mkdir fabric-test && cd fabric-test
git clone -b master https://github.com/hyperledger/fabric-samples.git
mkdir platform-specific-binaries
cd platform-specific-binaries
curl -sSL https://goo.gl/fMh2s3 | bash
export PATH=$(pwd)/bin:$PATH
cd ../fabric-samples/basic-network
./generate.sh
./start.sh
./stop.sh

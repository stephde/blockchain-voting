#!/bin/sh

# From http://hyperledger-fabric.readthedocs.io/en/latest/samples.html
# Installs the Hyperledger Sample Projects and Platform-specific Docker Images

mkdir platform-specific-binaries
cd platform-specific-binaries
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/release/scripts/bootstrap-1.0.4.sh | bash
export PATH=$(pwd)/bin:$PATH
echo $PATH
cd ../

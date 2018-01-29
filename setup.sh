#!/bin/sh

# From http://hyperledger-fabric.readthedocs.io/en/latest/samples.html
# Installs the Hyperledger Sample Projects and Platform-specific Docker Images

pgrep -f docker > /dev/null || {
  echo "Docker Daemon is not running"
  exit 1
}

mkdir platform-specific-binaries
cd platform-specific-binaries
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/132daddf1156a0c70429af2e8c3ca86bbdbf8c31/scripts/bootstrap-1.1.0-preview.sh | bash
export PATH=$(pwd)/bin:$PATH
echo $PATH
cd ../

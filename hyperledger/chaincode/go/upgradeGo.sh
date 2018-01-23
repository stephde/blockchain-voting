#!/usr/bin/env bash

go version
ls /opt/go
sudo rm -rf /opt/go

wget https://dl.google.com/go/go1.9.2.linux-amd64.tar.gz
tar -C /opt -xzf go1.9.2.linux-amd64.tar.gz
export PATH=/opt/go/bin:$PATH

go version
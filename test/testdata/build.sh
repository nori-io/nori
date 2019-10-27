#!/bin/bash

touch empty.so
go build -buildmode=plugin -o plugin.so plugin/plugin.go
go build -buildmode=plugin -o interface.so interface/plugin.go
go build -buildmode=plugin -o novariable.so novariable/plugin.go
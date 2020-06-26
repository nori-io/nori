#!/bin/bash

# todo: move command into Makefile (?)
touch empty.so
go build -buildmode=plugin -o plugin.so plugin/plugin.go
go build -buildmode=plugin -o no_interface.so interface/plugin.go
go build -buildmode=plugin -o no_variable.so novariable/plugin.go
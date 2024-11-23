#!/bin/bash

export CGO_ENABLED=1 
export CGO_CFLAGS=""
go build -tags "sqlite_math_functions"
#!/bin/bash
go get github.com/99designs/gqlgen/internal/imports
go get github.com/99designs/gqlgen/internal/code
go get github.com/99designs/gqlgen
#go run github.com/99designs/gqlgen generate
go run gen/genModel.go
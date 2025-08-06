#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o ./bin/ ./...


mv bin/nancalacc bin/nancalacc-linux-amd64

cgc build -f

cgc pack

cgc bundle -o ~/Downloads
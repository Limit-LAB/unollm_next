#!/bin/bash

set -e

rm -fr ./tmp
mkdir -p ./tmp

openapi-generator-cli generate -i ./UnoLLM.openapi.json -g go -o ./tmp --additional-properties=packageName=apimodel

rm -fr ./model/apimodel
mkdir -p ./model/apimodel
mv -f ./tmp/model_*.go ./model/apimodel
mv -f ./tmp/utils.go ./model/apimodel
rm -rf ./tmp
go fmt ./model/...
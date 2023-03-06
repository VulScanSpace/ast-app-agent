#!/bin/bash

build(){
  echo "building ast-app with $(go -v) in $(pwd)..."
  make build
  echo "build ast-app successful."
}

build

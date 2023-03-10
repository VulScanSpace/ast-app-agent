#!/bin/bash

buildAstApp() {
  echo "building ast-app with $(go -v) in $(pwd)..."
  make build
  echo "build ast-app successful."
}

buildAstFramework() {
  echo "start build ast-framework"
  # TODO write custom public key to os
  git clone git@github.com:VulScanSpace/AST-Framework.git
  mvn -DskipTests clean package
  echo "finished build ast-framework"
}

downloadJAttach() {
  echo "start download latest jattach"
  version=$(curl -sL "https://api.github.com/repos/jattach/jattach/releases/latest" | grep -E 'tag_name\": \"' | head -n 1 | tr -d 'tag_name\": ' | tr -d ',')
  curl -sL -o "jattach-linux" "https://github.com/jattach/jattach/releases/download/$(version)/jattach"
  echo "finished download jattach"
}

generateRpmSource() {
  echo "generate rpm source file..."
  cp build/ci/bin/* /
  cp AST-Framework/libs/* /
  echo "finished generate rpm source file..."
}

if [ "$1" == "astApp"]; then
  buildAstApp
elif [ "$1" == "astFramework" ]; then
  buildAstFramework
elif [ "$1" == "rpm" ]; then
  generateRpmSource
fi

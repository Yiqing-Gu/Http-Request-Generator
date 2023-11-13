#!/bin/bash

echo "please input os code (1:linux, 2:darwin, 3:windows):"
read oscode

date=`date +%Y%m%d`
if [ ${oscode} = 1 ]; then
  osname="linux"
  out="mockServer-linux"
fi
if [ ${oscode} = 2 ]; then
  osname="darwin"
  out="mockServer-darwin"
fi
if [ ${oscode} = 3 ]; then
  osname="windows"
  out="mockServer-windows.exe"
fi


echo "select build os: ${osname}"
echo "please wait for a while ..."
CGO_ENABLED=0 GOOS=${osname} GOARCH=amd64 go build -o dist/${date}/${out}
cp ./serverConfig.json dist/${date}/serverConfig.json
cp ./readme.md dist/${date}/readme.md
echo "build finish!"
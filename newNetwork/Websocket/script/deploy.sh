#!/usr/bin/env bash

func=$(certutil.exe -hashfile main.go md5)

hash="${func}"

echo $hash > hash.txt

echo " " &>> hash.txt


func=$(certutil.exe -hashfile server.wasm md5)

hash="${func}"

echo $hash &>> hash.txt

echo " " &>> hash.txt
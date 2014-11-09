#!/bin/bash

# TODO: Remove following line when script works
echo $(pwd)

wget https://storage.googleapis.com/appengine-sdks/featured/go_appengine_sdk_linux_amd64-1.9.14.zip -O go_appengine.zip -nv
unzip -q go_appengine.zip

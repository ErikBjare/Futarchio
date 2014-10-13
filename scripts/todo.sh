#!/bin/bash

HTMLFILES=`find ../src | grep "site/src/.*\.html$"`
GOFILES=`find ../src | grep ".*\.go$"`
JSFILES=`find ../src | grep "site/src/.*.js"`

grep TODO $HTMLFILES $GOFILES $JSFILES

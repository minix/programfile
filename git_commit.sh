#!/usr/bin/env bash

[ $# != 2 ] && echo "please specify commit directory and input commit content" && exit 1
cd $1
git add . && git commit -m "$2" && git push origin master

#!/bin/bash

CUR_PATH=`pwd`
echo "current work path: $CUR_PATH"

rm -rf $CUR_PATH/bin/nld_*

dlv=false
race=""


# 解析命令行选项
while getopts ":dr" opt; do
  case $opt in
    d)
      dlv=true
      echo "build in dlv mode"
      ;;
    r)
      race="-race"
      echo "build in race mode"
      ;;
    \?)
      echo "无效的选项: -$OPTARG" >&2
      exit 1
      ;;
    :)
      echo "选项 -$OPTARG 需要参数." >&2
      exit 1
      ;;
  esac
done

BuildServer() {
    echo "start build server $1"
#    cd $CUR_PATH/$1/server
    #go build -o "../../bin/nld_$1" -ldflags "-s -w" -race -v

    local buildStr="./bin/nld_$1"
    if [ "$2" == false ];then
      buildStr="$buildStr -v -ldflags '-s -w' -gcflags '-m -l'"
    else
      buildStr="$buildStr -v -gcflags 'all=-N -l'"
    fi

    echo "go build -o $buildStr $3"
    eval  "go build -o $buildStr $3"
}

BuildServer server $dlv $race

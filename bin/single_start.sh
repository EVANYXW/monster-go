#!/bin/bash

NAME=$1
./nld_server run --server_name $NAME &

echo "start ok, name:${NAME}"
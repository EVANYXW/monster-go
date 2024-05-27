#!/bin/bash

ps -ef | grep nld_ | grep -v grep | awk '{print $2}' | xargs kill -9

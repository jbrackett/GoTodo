#!/bin/bash

START_DEPS_CMD="docker compose up -d --build postgres_go"
START_CMD=""
DEBUG_CMD="-Xms128m -Xmx128m -Xdebug -Xrunjdwp:transport=dt_socket,address=*:5004,server=y,suspend=y"
EXEC="./main"

# load the local environment variable settings
#export $(cat local.env | xargs)
#export $(grep -v '^#' local.env | xargs)
while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ $line =~ ^[a-zA-Z_][a-zA-Z0-9_]*=.* ]]; then
        eval "export $line"
    fi
done < local.env

env

ARGS="-m"
if [[ "$1" != "" ]]; then
  ARGS=$1
fi

case $ARGS in
  -m|--minimal)
    eval $START_DEPS_CMD
    LOG_APPENDER=console $START_CMD $EXEC
    ;;
  -D|--debug)
    eval $START_DEPS_CMD
    $START_CMD $DEBUG_CMD $EXEC
    ;;
  -d|--docker)
    docker build -t gotodo .
    docker compose up -d
    docker compose logs -f -t gotodo
    ;;
  -h|--help)
    cat <<EOM
Usage: ./run.sh [-d|--docker] [-m|--minimal] [-j|--json] [-r|--raw] [-h|--help]
   -m --minimal   tiny messages (default)
   -d --docker    run in docker container
   -D --debug     run in debug mode
EOM
    ;;
esac

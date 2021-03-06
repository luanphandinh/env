#!/usr/bin/env bash

usage() {
  cli_name=${0##*/}
  echo "
${0##*/}
usage: ${cli_name} [options]
options:
  -i | --info        Log info to var/log/
  -d | --debug       Log debug to var/log/
  -e | --error       Log error to var/log/
  -w | --warning     Log warning to var/log/
  *                  Help
"
  exit 1
}

color() {
  # Black        0;30     Dark Gray     1;30
  # Red          0;31     Light Red     1;31
  # Green        0;32     Light Green   1;32
  # Brown/Orange 0;33     Yellow        1;33
  # Blue         0;34     Light Blue    1;34
  # Purple       0;35     Light Purple  1;35
  # Cyan         0;36     Light Cyan    1;36
  # Light Gray   0;37     White         1;37
  COLOR='\033[0;32m'
  case $1 in
  info)
    COLOR='\033[0;36m'
    ;;

  debug)
    COLOR='\033[0;32m'
    ;;

  error)
    COLOR='\033[0;31m'
    ;;

  warning)
    COLOR='\033[1;33m'
    ;;
  esac

  echo $COLOR
}

log() {
  type=$1
  BEGIN_COLOR="$(color ${type})"
  END_COLOR='\033[0m'

  shift
  echo -e "[ ${BEGIN_COLOR}${type}${END_COLOR} ][$(date)]: $@"

  log_file="${__VAR_LOG_DIR__}/$(date +'%Y%m%d').log"
  echo "[ ${type} ][$(date)]: $@" >> "${log_file}"
}

while [ "$1" != "" ]; do
  case $1 in
  -i | --info)
    shift
    log info $@
    exit
    ;;

  -d | --debug)
    shift
    log debug $@
    exit
    ;;

  -e | --error)
    shift
    log error $@
    exit
    ;;

  -w | --warning)
    shift
    log warning $@
    exit
    ;;

  -h | --help)
    usage
    exit
    ;;

  *)
    usage
    exit 1
    ;;
  esac
  shift
done

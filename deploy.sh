#!/bin/bash
function error_exit {
	echo [$(date +%Y/%m/%d:%H:%M:%S)] 1>&2
	exit 1
}

function log {
  echo [$(date +%Y/%m/%d:%H:%M:%S)] $1
}

echo "---------------------------------------------------------------------------------------------------"
BUILDDATE=`date -u "+%Y:%m:%d %H:%M:%S"`
log "start building version: ${BUILDDATE}"


if rm -rf ./target; then
  log "target directory cleaned"
else
  error_exit "error while deleting target directory"
fi


echo "---------------------------------------------------------------------------------------------------"

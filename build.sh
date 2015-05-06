#!/bin/bash
function error_exit {
	echo [$(date +%Y/%m/%d:%H:%M:%S)] 1>&2
	exit 1
}

function log {
  echo [$(date +%Y/%m/%d:%H:%M:%S)] $1
}

echo "---------------------------------------------------------------------------------------------------"
log "start building version: ${BUILD_NUMBER}"

if rm -rf ./target; then
  log "target directory cleaned"
else
  error_exit "error while deleting target directory"
fi

if mkdir -p ./target; then
  log "new target direcory created"
else
  error_exit "error while creating target directory"
fi

if BUILDDATE=`date -u "+%Y:%m:%d %H:%M:%S"` && CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-X main.version '${VERSION}' -X main.builddate '${BUILDDATE}'" -o ./target/go-read main.go; then
  log "go build completed"
else
  error_exit "error while building static binary"
fi

if cp -r ./static ./target/static; then
  log "static files copied to target directory"
else
  error_exit "error while copying static files to target directory"
fi

echo "---------------------------------------------------------------------------------------------------"

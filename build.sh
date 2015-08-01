#!/bin/bash
function error_exit {
	echo [$(date +%Y/%m/%d:%H:%M:%S)] 1>&2
	exit 1
}

function log {
  echo [$(date +%Y/%m/%d:%H:%M:%S)] $1
}

echo "---------------------------------------------------------------------------------------------------"
BUILDDATE=`date "+%Y:%m:%d %H:%M:%S"`
log "start building version: ${BUILDDATE}"

if rm -rf ./target; then
  log "target directory cleaned"
else
  error_exit "error while deleting target directory"
fi

if mkdir -p ./target/logs; then
  log "logs directory created"
else
  error_exit "error while deleting logs directory"
fi

if mkdir -p ./target; then
  log "new target direcory created"
else
  error_exit "error while creating target directory"
fi

if go build -a -ldflags "-X main.BuildDate '${BUILDDATE}'" -tags netgo -installsuffix netgo -o ./target/go-read main.go; then
  log "go build completed"
else
  error_exit "error while building static binary"
fi

# go-bindata -o statics.go static/...

if cp -r ./static ./target/static; then
  log "static files copied to target directory"
else
  error_exit "error while copying static files to target directory"
fi

if cp ./run.sh ./target/; then
  log "copy run.sh to target"
else
  error_exit "error while copying run.sh to target directory"
fi

echo "---------------------------------------------------------------------------------------------------"

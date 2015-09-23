#!/bin/bash
function error_exit {
	echo [$(date +%Y/%m/%d:%H:%M:%S)] 1>&2
	exit 1
}

function log {
  echo [$(date +%Y/%m/%d:%H:%M:%S)] $1
}

echo "---------------------------------------------------------------------------------------------------"
install_dir="/smb/www/go-read"
log "start deploying to ${install_dir}"

if pkill go-read; then
  log "service stopped"
else
  log "NOTICE: go-read did not run at all"
fi

if rm -rf ${install_dir}/static ${install_dir}/go-read; then
  log "deleted old version in ${install_dir}"
else
  error_exit "something wrong deleting old version"
fi

if cp -R ./target/* ${install_dir}; then
  log "copied new version to ${install_dir}"
else
  error_exit "error copying to installtion directory"
fi

if tmux send -t server:3 ./go-read ENTER; then
  log "Restarted go-read"
else
  error_exit "error restarting go-read after deployment"
fi

echo "---------------------------------------------------------------------------------------------------"

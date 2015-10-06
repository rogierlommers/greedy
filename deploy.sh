#!/bin/bash
function error_exit {
	echo [$(date +%Y/%m/%d:%H:%M:%S)] 1>&2
	exit 1
}

function log {
  echo [$(date +%Y/%m/%d:%H:%M:%S)] $1
}

echo "---------------------------------------------------------------------------------------------------"
install_dir="/smb/www/greedy"
log "start deploying to ${install_dir}"

if pkill greedy; then
  log "service stopped"
else
  log "NOTICE: greedy did not run at all"
fi

if rm -rf ${install_dir}/static ${install_dir}/greedy; then
  log "deleted old version in ${install_dir}"
else
  error_exit "something wrong deleting old version"
fi

if cp -R ./target/* ${install_dir}; then
  log "copied new version to ${install_dir}"
else
  error_exit "error copying to installtion directory"
fi

if tmux send -t server:3 ./greedy ENTER; then
  log "Restarted greedy"
else
  error_exit "error restarting greedy after deployment"
fi

echo "---------------------------------------------------------------------------------------------------"

package main

import log "gopkg.in/inconshreveable/log15.v2"

func main() {
	log.Info("systemtest", "status", "starting")

	log.Info("systemtest", "status", "done")
}

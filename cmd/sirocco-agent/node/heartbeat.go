package node

import (
	"time"
	"sirocco-agent/config"
)

func StartHeartbeat(cfg config.Config) {
	for {
		// send heartbeat to controller
		time.Sleep(5 * time.Second)
	}
}
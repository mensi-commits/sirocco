package shard

import (
	"os"
	"os/exec"
	"sirocco-agent/config"
)

func Delete(cfg config.Config, name string) error {
	exec.Command("docker", "rm", "-f", name).Run()
	return os.RemoveAll(cfg.DataDir + "/" + name)
}
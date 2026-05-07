package shard

import (
	"fmt"
	"os/exec"
	"sirocco-agent/config"
)

func CreateShard(cfg config.Config, name string, port int, password string) error {

	dataPath := fmt.Sprintf("%s/%s", cfg.DataDir, name)

	cmd := exec.Command("docker", "run", "-d",
		"--name", name,
		"-e", "MYSQL_ROOT_PASSWORD="+password,
		"-p", fmt.Sprintf("%d:3306", port),
		"-v", fmt.Sprintf("%s:/var/lib/mysql", dataPath),
		"mysql:8.0",
	)

	return cmd.Run()
}
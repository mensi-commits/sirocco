package shard

import "os/exec"

func Start(name string) error {
	return exec.Command("docker", "start", name).Run()
}

func Stop(name string) error {
	return exec.Command("docker", "stop", name).Run()
}
package docker

import "os/exec"

func Run(args ...string) error {
	return exec.Command("docker", args...).Run()
}
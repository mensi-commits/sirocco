package sirocco

import (
	"fmt"
	"os/exec"
)

// Frankenstrike creates a MySQL shard container with persistent storage
func Frankenstrike(volumeName string, rootPassword string) (string, error) {

	containerName := "Frankenstrike"

	// Create persistent volume
	cmdVolume := exec.Command(
		"docker", "volume", "create", volumeName,
	)
	if err := cmdVolume.Run(); err != nil {
		return "", fmt.Errorf("failed to create volume: %w", err)
	}

	// Run MySQL shard container
	cmd := exec.Command(
		"docker", "run", "-d",
		"--name", containerName,

		// persistent storage
		"-v", fmt.Sprintf("%s:/var/lib/mysql", volumeName),

		// MySQL config
		"-e", fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", rootPassword),

		// optional port exposure (can be removed in production)
		"-p", "3306",

		// image
		"mysql:8.0",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create Frankenstrike container: %w - %s", err, string(output))
	}

	return string(output), nil
}
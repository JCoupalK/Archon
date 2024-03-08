package utils

import (
	"log"
	"path/filepath"
	"strings"
)

// BackupFS orchestrates the backup process, handling file transfers based on specified flags.
func BackupFS(username, password, sshKeyFile, hostname string, port int, sourceDir, outputDir, exclude string, recursive, compress bool, direction string) {
	// Establish the SSH connection
	client, err := establishSSHConnection(username, password, sshKeyFile, hostname, port)
	if err != nil {
		log.Fatalf("Failed to establish SSH connection: %v", err)
	}
	defer client.Close()

	// Determine the operation direction and perform the file transfer
	switch direction {
	case "to-remote":
		// Validate or prepare the source and destination paths
		if !filepath.IsAbs(sourceDir) {
			sourceDir, err = filepath.Abs(sourceDir)
			if err != nil {
				log.Fatalf("Failed to get absolute path of source directory: %v", err)
			}
		}

		// This example assumes CopyLocalToRemote is a function you would implement
		// similar to CopyRemoteToLocal, but for uploading files to the remote server.
		err = CopyLocalToRemote(client, sourceDir, outputDir, []string{exclude}, recursive)
		if err != nil {
			log.Fatalf("Failed to copy files to remote: %v", err)
		}
	case "from-remote":
		err = CopyRemoteToLocal(client, sourceDir, outputDir, strings.Split(exclude, ","), recursive)
		if err != nil {
			log.Fatalf("Failed to copy files from remote: %v", err)
		}
	default:
		log.Fatalf("Invalid direction specified: %s. Use 'to-remote' or 'from-remote'.", direction)
	}
}

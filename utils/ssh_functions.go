package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// Establish an SSH connection using either a password or an SSH key.
func establishSSHConnection(username, password, sshKeyFile, hostname string, port int) (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	if sshKeyFile != "" {
		key, err := os.ReadFile(sshKeyFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %w", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	} else {
		return nil, fmt.Errorf("either sshKeyFile or password must be provided")
	}

	// Path to the known hosts file
	knownHostsFile := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	hostKeyCallback, err := knownhosts.New(knownHostsFile)
	if err != nil {
		return nil, fmt.Errorf("could not create hostkeycallback function: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return client, nil
}

// RemoteFileList retrieves a list of files and directories from a remote path.
func RemoteFileList(client *ssh.Client, remotePath string) ([]string, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// list files
	cmd := "ls -a " + remotePath
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(output), "\n")
	return files, nil
}

func copyFileFromRemote(sftpClient *sftp.Client, remotePath, localPath string) error {
	// Open the remote file for reading.
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file (%s): %w", remotePath, err)
	}
	defer remoteFile.Close()

	// Create the local file.
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file (%s): %w", localPath, err)
	}
	defer localFile.Close()

	// Copy the contents of the remote file to the local file.
	if _, err := io.Copy(localFile, remoteFile); err != nil {
		return fmt.Errorf("failed to copy file from remote (%s) to local (%s): %w", remotePath, localPath, err)
	}

	return nil
}

func CopyRemoteToLocal(client *ssh.Client, remotePath, localPath string, exclusions []string, recursive bool) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	// Adjust this function to correctly handle the creation of local files within directories.
	walker := sftpClient.Walk(remotePath)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			fmt.Println(err)
			continue
		}

		relativePath, err := filepath.Rel(remotePath, walker.Path())
		if err != nil {
			return fmt.Errorf("failed to calculate relative file path: %v", err)
		}

		localFilePath := filepath.Join(localPath, relativePath)

		// Skip if it matches any of the exclusion patterns
		if matchesExclusion(relativePath, exclusions) {
			continue
		}

		if walker.Stat().IsDir() {
			// Create directory if it does not exist
			if err := os.MkdirAll(localFilePath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create local directory: %v", err)
			}
		} else {
			// Prepare to copy file, ensure the containing directory exists
			if err := os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm); err != nil {
				return fmt.Errorf("failed to create local directory for file: %v", err)
			}

			// Now, copy the file from remote to local
			err := copyFileFromRemote(sftpClient, walker.Path(), localFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyLocalToRemote uploads files or directories from the local file system to a remote system.
// It supports recursive copying and exclusions.
func CopyLocalToRemote(client *ssh.Client, localPath, remotePath string, exclusions []string, recursive bool) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	return uploadDirectory(sftpClient, localPath, remotePath, exclusions, recursive)
}

// uploadDirectory recursively uploads a directory and its contents to the remote server.
func uploadDirectory(sftpClient *sftp.Client, localPath, remotePath string, exclusions []string, recursive bool) error {
	localFiles, err := os.ReadDir(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local directory: %v", err)
	}

	// Ensure the remote directory exists
	err = sftpClient.MkdirAll(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote directory: %v", err)
	}

	for _, entry := range localFiles {
		localEntryPath := filepath.Join(localPath, entry.Name())
		remoteEntryPath := filepath.Join(remotePath, entry.Name())

		// Check for exclusions
		if matchesExclusion(entry.Name(), exclusions) {
			continue
		}

		if entry.IsDir() {
			if recursive {
				err = uploadDirectory(sftpClient, localEntryPath, remoteEntryPath, exclusions, recursive)
				if err != nil {
					return err
				}
			}
		} else {
			err = uploadFile(sftpClient, localEntryPath, remoteEntryPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// uploadFile uploads a single file to the remote server.
func uploadFile(sftpClient *sftp.Client, localFilePath, remoteFilePath string) error {
	localFile, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %v", err)
	}
	defer remoteFile.Close()

	_, err = remoteFile.ReadFrom(localFile)
	if err != nil {
		return fmt.Errorf("failed to copy local file to remote: %v", err)
	}

	return nil
}

// matchesExclusion checks if the file matches any of the provided exclusion patterns.
func matchesExclusion(file string, exclusions []string) bool {
	for _, pattern := range exclusions {
		matched, err := filepath.Match(pattern, file)
		if err != nil {
			fmt.Printf("Error matching pattern: %v\n", err)
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

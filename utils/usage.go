package utils

import (
	"fmt"
	"os"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "General options:")
	fmt.Fprintln(os.Stderr, "  -t, --type             	Backup type (mysql, mariadb, fs)")
	fmt.Fprintln(os.Stderr, "  -od, --outputdir       	Output directory")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Database backup options (mysql, mariadb):")
	fmt.Fprintln(os.Stderr, "  -u, --user             	Database username")
	fmt.Fprintln(os.Stderr, "  -pw, --password        	Database password")
	fmt.Fprintln(os.Stderr, "  -h, --host             	Database host")
	fmt.Fprintln(os.Stderr, "  -p, --port             	Database port")
	fmt.Fprintln(os.Stderr, "  -d, --database         	Database name")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Filesystem backup options (fs):")
	fmt.Fprintln(os.Stderr, "  -u, --user             	SSH username (for fs type)")
	fmt.Fprintln(os.Stderr, "  -pw, --password        	SSH password (for fs type)")
	fmt.Fprintln(os.Stderr, "  -h, --host             	SSH host (for fs type)")
	fmt.Fprintln(os.Stderr, "  -p, --port             	SSH port (for fs type)")
	fmt.Fprintln(os.Stderr, "  -sd, --sourcedir       	Source directory for filesystem backup")
	fmt.Fprintln(os.Stderr, "  -r, --recursive        	Recursive backup (includes subdirectories)")
	fmt.Fprintln(os.Stderr, "  -e, --exclude          	Exclude pattern (files/directories to skip)")
	fmt.Fprintln(os.Stderr, "  -sk, --sshkey          	SSH key file for authentication (alternative to password)")
	fmt.Fprintln(os.Stderr, "  -c, --compress         	Compress the backup")
	fmt.Fprintln(os.Stderr, "  -dir, --direction      	Backup direction (to-remote, from-remote)")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Ensure all required flags are provided for the selected backup type.")
	os.Exit(1)
}

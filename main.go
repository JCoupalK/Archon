package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/KeepSec-Technologies/Archon/utils"
	// "github.com/KeepSec-Technologies/Archon/utils/restore"
)

var (
	// Long-form flags
	Type       string
	username   string
	password   string
	hostname   string
	db         string
	port       int
	outputDir  string
	sourceDir  string
	recursive  bool
	exclude    string
	sshKeyFile string
	compress   bool
	direction  string

	// Short-form flags
	TypeShort       string
	usernameShort   string
	passwordShort   string
	hostnameShort   string
	dbShort         string
	portShort       int
	outputDirShort  string
	sourceDirShort  string
	recursiveShort  bool
	excludeShort    string
	sshKeyFileShort string
	compressShort   bool
	directionShort  string
)

func init() {
	// Long-form flags
	flag.StringVar(&Type, "type", "", "Backup type (mysql, mariadb, fs)")
	flag.StringVar(&username, "user", "", "Database or SSH username (depends on type)")
	flag.StringVar(&password, "password", "", "Database or SSH password (depends on type)")
	flag.StringVar(&hostname, "host", "", "Database or SSH host (depends on type)")
	flag.IntVar(&port, "port", 0, "Database or SSH port (depends on type)")
	flag.StringVar(&db, "database", "", "Database name")
	flag.StringVar(&outputDir, "outputdir", "", "Output directory")
	flag.StringVar(&sourceDir, "sourcedir", "", "Source directory for filesystem backup")
	flag.BoolVar(&recursive, "recursive", false, "Recursive backup (includes subdirectories)")
	flag.StringVar(&exclude, "exclude", "", "Exclude pattern (files/directories to skip)")
	flag.StringVar(&sshKeyFile, "sshkey", "", "SSH key file for authentication")
	flag.BoolVar(&compress, "compress", false, "Compress the backup")
	flag.StringVar(&direction, "direction", "", "Backup direction (to-remote, from-remote)")

	// Short-form flags
	flag.StringVar(&TypeShort, "t", "", "Backup type (mysql, mariadb, fs)")
	flag.StringVar(&usernameShort, "u", "root", "Database or SSH username (depends on type)")
	flag.StringVar(&passwordShort, "pw", "", "Database or SSH password (depends on type)")
	flag.StringVar(&hostnameShort, "h", "", "Database or SSH host (depends on type)")
	flag.IntVar(&portShort, "p", 0, "Database or SSH port (depends on type)")
	flag.StringVar(&dbShort, "d", "", "Database name")
	flag.StringVar(&outputDirShort, "od", "", "Output directory")
	flag.StringVar(&sourceDirShort, "sd", "", "Source directory for filesystem backup")
	flag.BoolVar(&recursiveShort, "r", false, "Recursive backup (includes subdirectories)")
	flag.StringVar(&excludeShort, "e", "", "Exclude pattern (files/directories to skip)")
	flag.StringVar(&sshKeyFileShort, "sk", "", "SSH key file for authentication")
	flag.BoolVar(&compressShort, "c", false, "Compress the backup")
	flag.StringVar(&directionShort, "di", "", "Backup direction (to-remote, from-remote)")

	// Override the default flag.Usage
	flag.Usage = utils.Usage

	flag.Parse()

	// Apply short-form flags if set
	if TypeShort != "" {
		Type = TypeShort
	}
	if usernameShort != "" {
		username = usernameShort
	}
	if passwordShort != "" {
		password = passwordShort
	}
	if hostnameShort != "" {
		hostname = hostnameShort
	}
	if portShort != 0 {
		port = portShort
	}
	if dbShort != "" {
		db = dbShort
	}
	if outputDirShort != "" {
		outputDir = outputDirShort
	}
	if sourceDirShort != "" {
		sourceDir = sourceDirShort
	}
	if recursiveShort {
		recursive = recursiveShort
	}
	if excludeShort != "" {
		exclude = excludeShort
	}
	if sshKeyFileShort != "" {
		sshKeyFile = sshKeyFileShort
	}
	if compressShort {
		compress = compressShort
	}
	if directionShort != "" {
		direction = directionShort
	}
}

func main() {

	switch Type {
	case "mysql", "mariadb":
		if username == "" || password == "" || hostname == "" || port == 0 || db == "" {
			fmt.Fprintln(os.Stderr, "Error: Required flags are missing.")
			utils.Usage()
		}
		if outputDir == "" {
			outputDir, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			utils.BackupMySQL(username, password, hostname, db, outputDir, port)
		} else {
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputDir := fmt.Sprintf("%s/%s", pwd, outputDir)
			utils.BackupMySQL(username, password, hostname, db, outputDir, port)
		}
	case "fs":
		// Check required flags for filesystem backup
		if username == "" || password == "" && sshKeyFile == "" || hostname == "" || port == 0 || sourceDir == "" || outputDir == "" || direction == "" {
			fmt.Fprintln(os.Stderr, "Error: Required flags for filesystem backup are missing.")
			utils.Usage()
		}

		// Call to filesystem backup function (to be implemented in backup_fs.go)
		utils.BackupFS(username, password, sshKeyFile, hostname, port, sourceDir, outputDir, exclude, recursive, compress, direction)
	default:
		fmt.Fprintln(os.Stderr, "Error: Required flags are missing.")
		utils.Usage()
	}
}

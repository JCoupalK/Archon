package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/JCoupalK/Archon/utils"
)

var (
	// Long-form flags
	Type      string
	username  string
	password  string
	hostname  string
	db        string
	port      int
	outputDir string

	// Short-form flags
	TypeShort      string
	usernameShort  string
	passwordShort  string
	hostnameShort  string
	dbShort        string
	portShort      int
	outputDirShort string
)

func init() {
	// Define flags
	flag.StringVar(&Type, "type", "", "Backup type (mysql, mariadb)")
	flag.StringVar(&username, "user", "", "Database username")
	flag.StringVar(&password, "password", "", "Database password")
	flag.StringVar(&hostname, "host", "", "Database host")
	flag.IntVar(&port, "port", 0, "Database port")
	flag.StringVar(&db, "database", "", "Database name")
	flag.StringVar(&outputDir, "outputdir", "", "Output directory for the backup files")

	// Short-form flags
	flag.StringVar(&TypeShort, "t", "", "Backup type (mysql, mariadb)")
	flag.StringVar(&usernameShort, "u", "", "Database username")
	flag.StringVar(&passwordShort, "p", "", "Database password")
	flag.StringVar(&hostnameShort, "h", "", "Database host")
	flag.IntVar(&portShort, "P", 0, "Database port")
	flag.StringVar(&dbShort, "d", "", "Database name")
	flag.StringVar(&outputDirShort, "od", "", "Output directory for the backup files")

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
	case "postgresql", "postgres":
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
			utils.BackupPostgreSQL(username, password, hostname, db, outputDir, port)
		} else {
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputDir := fmt.Sprintf("%s/%s", pwd, outputDir)
			utils.BackupPostgreSQL(username, password, hostname, db, outputDir, port)
		}
	default:
		fmt.Fprintln(os.Stderr, "Error: Required flags are missing.")
		utils.Usage()
	}
}

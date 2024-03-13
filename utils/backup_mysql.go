package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jamf/go-mysqldump"
)

func BackupMySQL(username, password, hostname, dbname, outputDir string, port int) {
	// Open connection to database
	config := mysql.NewConfig()
	config.User = username
	config.Passwd = password
	config.DBName = dbname
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", hostname, port)

	dumpDir := outputDir
	currentTime := time.Now()
	dumpFilenameFormat := fmt.Sprintf("%s-%s", config.DBName, currentTime.Format("20060102T150405"))

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	if err != nil {
		fmt.Println("Error registering database:", err)
		return
	}

	// Capture the actual file name from dumper.Out before checking for an error
	var actualDumpFilePath string
	if file, ok := dumper.Out.(*os.File); ok {
		actualDumpFilePath = file.Name()
	}

	// Attempt to dump database to file
	if err := dumper.Dump(); err != nil {
		fmt.Println("Error dumping:", err)
		os.Remove(actualDumpFilePath)
		// Error handling
	} else {
		// Rename the file if the actual dump file path doesn't match the expected format
		expectedFilePath := fmt.Sprintf("%s/%s.sql", dumpDir, dumpFilenameFormat)
		if actualDumpFilePath != expectedFilePath {
			if err := os.Rename(actualDumpFilePath, expectedFilePath); err != nil {
				fmt.Printf("Failed to rename the dump file from '%s' to '%s': %v\n", actualDumpFilePath, expectedFilePath, err)
			} else {
				fmt.Println("Backup successfully saved to", expectedFilePath)
			}
		}
	}
	// Close dumper, connected database, and file stream.
	dumper.Close()
}

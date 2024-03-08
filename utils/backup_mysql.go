package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", config.DBName)

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return
	}

	// Dump database to file
	if err := dumper.Dump(); err != nil {
		fmt.Println("Error dumping:", err)
		return
	}
	if file, ok := dumper.Out.(*os.File); ok {
		fmt.Println("Backup successfully saved to", file.Name())
	} else {
		fmt.Println("It's not part of *os.File, but dump is done")
	}

	// Close dumper, connected database and file stream.
	dumper.Close()
}

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/JCoupalK/go-pgdump" // Ensure this path matches your library's import path
)

func BackupPostgreSQL(username, password, hostname, dbname, outputDir string, port int) {
	// PostgreSQL connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		hostname, port, username, password, dbname)

	// Prepare filename for the unified dump
	currentTime := time.Now()
	dumpFilename := filepath.Join(outputDir, fmt.Sprintf("%s-%s.sql", dbname, currentTime.Format("20060102T150405")))

	// Create a new dumper instance
	dumper := pgdump.NewDumper(psqlInfo)

	if err := dumper.DumpDatabase(dumpFilename); err != nil {
		fmt.Printf("Error dumping database: %v", err)
		os.Remove(dumpFilename) // Cleanup on failure
		return
	}

	fmt.Println("Backup successfully saved to", dumpFilename)
}

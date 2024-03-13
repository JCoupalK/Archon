package utils

import (
	"fmt"
	"os"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "General options:")
	fmt.Fprintln(os.Stderr, "  -t,    --type             	Database type (mysql, mariadb)")
	fmt.Fprintln(os.Stderr, "  -u,    --user             	Database username")
	fmt.Fprintln(os.Stderr, "  -p,   --password        	Database password")
	fmt.Fprintln(os.Stderr, "  -h,    --host             	Database host")
	fmt.Fprintln(os.Stderr, "  -P,    --port             	Database port")
	fmt.Fprintln(os.Stderr, "  -d,    --database         	Database name")
	fmt.Fprintln(os.Stderr, "  -od,   --outputdir       	Output directory")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Ensure all required flags are provided for the selected backup type.")
	os.Exit(1)
}

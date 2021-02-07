package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
)

const usage = `Usage: %s [OPTIONS] [file1 file2 ...]

OPTIONS:
`

const appVersion = "ImageFix v0.1.0"

func main() {
	dryRunPtr := flag.Bool("dry-run", false, "Do not modify files (dry-run)")
	helpPtr := flag.Bool("help", false, "Display this help info")
	verbosePrt := flag.Bool("verbose", false, "Display verbose output")
	versionPrt := flag.Bool("version", false, "Display app version")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, filepath.Base(os.Args[0]))

		flag.VisitAll(func(f *flag.Flag) {
			optionName := fmt.Sprintf("-%s", f.Name)
			fmt.Fprintf(flag.CommandLine.Output(), "  %-6s  %s (default: %v)\n", optionName, f.Usage, f.DefValue) //
		})
	}

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		fmt.Println()
		os.Exit(0)
	}

	if !*dryRunPtr && !*helpPtr && !*verbosePrt && !*versionPrt && len(flag.Args()) == 0 {
		flag.Usage()
		fmt.Println()
		os.Exit(1)
	}

	if *versionPrt {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	for _, filePath := range flag.Args() {
		renameFile := false

		f, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error Opening File: %s\n", err.Error())
			continue
		}
		fileDir, fileName := filepath.Split(filePath)
		if len(fileDir) > 0 && len(fileName) == 0 {
			fileInfo, err := os.Lstat(fileDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "LStat Error: %s\n", err.Error())
			}
			if fileInfo.IsDir() == false {
				fileName = fileDir // NOTE: Should check if a long path and pop the filename? May only happen with no dir present in path.
				fileDir = ""
			}
		}
		fileBase := fileName
		fileExt := filepath.Ext(strings.ToLower(fileName))
		if len(fileExt) > 0 {
			baselength := len(fileName) - len(fileExt)
			fileBase = fileName[0:baselength]
		}
		// fmt.Printf("filePath: %q | fileDir: %q | fileName: %q | fileExt: %q\n", filePath, fileDir, fileName, fileExt)

		// fmt.Printf("filePath: %q |\n\tfileDir: %q | fileName: %q\n\tfileBase: %q | fileExt: %q\n", filePath, fileDir, fileName, fileBase, fileExt)

		// We only have to pass the file header = first 261 bytes
		head := make([]byte, 261)
		f.Read(head)

		if len(head) == 0 || filetype.IsImage(head) == false {
			if *verbosePrt {
				fmt.Fprintf(os.Stderr, "Skipping %q (it is not an image)\n", filePath)
			}
			continue
		}

		if filetype.IsImage(head) {
			kind, _ := filetype.Match(head)
			switch fileExt {
			case ".jpeg":
				renameFile = !filetype.Is(head, "jpg")
			case ".tiff":
				renameFile = !filetype.Is(head, "tif")
			default:
				if len(fileExt) > 0 {
					if fileExt[1:] != kind.Extension {
						renameFile = true
					}
				} else {
					renameFile = true
				}
			}

			if renameFile {
				fileNewName := fmt.Sprintf("%s.%s", fileBase, kind.Extension)
				fileNewPath := fmt.Sprintf("%s%s.%s", fileDir, fileBase, kind.Extension)

				newFileExists := true
				_, err := os.Stat(fileNewPath)
				if err != nil && os.IsNotExist(err) {
					newFileExists = false
				}

				if newFileExists {
					layout := "Not renaming %q to %q as the later already exists"
					if *dryRunPtr {
						layout += " (dry-run)"
					}
					fmt.Fprintf(os.Stderr, layout+"\n", filePath, fileNewName)
				} else {
					layout := "Renaming %q to %q"
					if *dryRunPtr {
						layout += " (dry-run)"
					}
					fmt.Printf(layout+"\n", filePath, fileNewPath)

					if *dryRunPtr == false {
						err := os.Rename(filePath, fileNewPath)
						if err != nil {
							fmt.Fprintf(os.Stderr, "File Rename Error: %s\n", err.Error())
						}
					}
				}
			} else {
				if *verbosePrt {
					layout := "Skipping %q (no fixing needed)"
					if *dryRunPtr {
						layout += " (dry-run)"
					}
					fmt.Printf(layout+"\n", filePath)
				}
			}
		}
	}
}

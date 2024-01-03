package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync source dir to destination",
	Long: `Sync source directory to destination.

Only regular files and directories are supported.`,
	RunE: syncRun,
	Args: cobra.ExactArgs(2),
}

func syncRun(cmd *cobra.Command, args []string) error {
	ignore := []string{".skan-filesync"}
	srcFiles, err := walkdir(args[0], ignore)
	if err != nil {
		return err
	}
	log.Printf("src files:")
	printList(srcFiles)

	dstFiles, err := walkdir(args[1], ignore)
	if err != nil {
		return err
	}
	log.Printf("dst files:")
	printList(dstFiles)

	return nil
}

func printList(l []fileInfo) {
	for _, f := range l {
		d := "r"
		if f.IsDir {
			d = "d"
		}
		log.Printf("%s %s %10d %s", d, f.Perm.String(), f.Size, f.Name)
	}
}

type fileInfo struct {
	Name  string
	Size  int64
	IsDir bool
	Perm  os.FileMode
}

func walkdir(root string, ignore []string) ([]fileInfo, error) {
	ret := []fileInfo{}
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// make sure the root is a direcory, and otherwise skip it.
			if path == root {
				if !info.IsDir() {
					return fmt.Errorf("%s is not a directory", root)
				}
				return nil
			}
			name, _ := strings.CutPrefix(path, root+"/")
			if info.IsDir() || info.Mode().IsRegular() {
				ret = append(ret, fileInfo{
					Name:  name,
					Size:  info.Size(),
					IsDir: info.IsDir(),
					Perm:  info.Mode().Perm(),
				})
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

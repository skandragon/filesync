package cmd

import (
	"log"
	"os"
	"path/filepath"

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
	files, err := walkdir(args[0])
	if err != nil {
		return err
	}
	log.Printf("files: %#v", files)
	return nil
}

type fileInfo struct {
	Name  string
	Size  int64
	IsDir bool
	Perm  os.FileMode
}

func walkdir(root string) ([]fileInfo, error) {
	ret := []fileInfo{}
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path != root && (info.IsDir() || info.Mode().IsRegular()) {
				ret = append(ret, fileInfo{
					Name:  info.Name(),
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

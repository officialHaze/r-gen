package helper

import (
	"os"
	"reportgenengine/interfaces"
	"reportgenengine/log"
	"strings"
)

// RmFileExt separates a filename into name and extension
// and returns only the name without any trailing extension
func RmFileExt(name string) string {
	parts := strings.Split(name, ".")
	// The last part should always be an extension
	partswithoutext := parts[:len(parts)-1]
	// re-join the extracted parts
	joined := strings.Join(partswithoutext, ".")

	return joined
}

// ListFileInfo reads a dir and creates a separate
// list of only file entries with their info as per a
// custom FileInfo schema. And returns the created list
func ListFileInfo(dir string) []interfaces.FileInfo {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Errorf("%v", err)
	}

	infos := make([]interfaces.FileInfo, 0, len(entries))

	for _, e := range entries {
		// Skip in case a dir
		if e.IsDir() {
			continue
		}
		i, _ := e.Info()

		id := RmFileExt(i.Name())
		info := interfaces.FileInfo{
			ID:      id,
			ModTime: i.ModTime(),
			Size:    i.Size(),
		}
		infos = append(infos, info)
	}

	return infos
}

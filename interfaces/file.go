package interfaces

import "time"

type FileInfo struct {
	ID      string    `json:"fileId"`
	ModTime time.Time `json:"modifiedAt"`
	Size    int64     `json:"size"`
}

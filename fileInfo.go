package bigbrother

import "time"

type FileInfo struct {
	ID   string //fixme:uuid maybe ??
	Name string
	Path string

	CreateAt time.Time
	Version  int
}

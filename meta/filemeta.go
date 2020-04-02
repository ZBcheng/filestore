package meta

// FileMeta : file struct
type FileMeta struct {
	FileID   string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta : add or update a file
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileID] = fmeta
}

// GetFileMeta : get a file
func GetFileMeta(fileID string) FileMeta {
	return fileMetas[fileID]
}

// RemoveFileMeta : delete a file
func RemoveFileMeta(fileID string) {
	delete(fileMetas, fileID)
}

func GetFileSize(filename string, filepath string) {}
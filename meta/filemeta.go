package meta

// FileMeta : file struct
type FileMeta struct {
	FileHash string
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
	fileMetas[fmeta.FileHash] = fmeta
}

// GetFileMeta : get a file
func GetFileMeta(fileHash string) FileMeta {
	return fileMetas[fileHash]
}

// RemoveFileMeta : delete a file
func RemoveFileMeta(fileHash string) {
	delete(fileMetas, fileHash)
}

func GetFileSize(filename string, filepath string) {}

package store

type FileSystem struct {
	Spaces     map[int]string
	NOpen      int //	ibd文件打开数量
	NameHashes map[string]int
}

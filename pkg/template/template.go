package template

type File struct {
	name     string
	filePath string
}

func (f *File) Name() string {
	return f.name
}

func (f *File) FilePath() string {
	return f.filePath
}

func NewFile(name string, filePath string) *File {
	return &File{
		name:     name,
		filePath: filePath,
	}
}

type FileStorage interface {
	List() ([]*File, error)
}

package template

type FilePath string

type File struct {
	name string
	path FilePath
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Path() FilePath {
	return f.path
}

func NewFile(name string, filePath string) *File {
	return &File{
		name: name,
		path: FilePath(filePath),
	}
}

type FileStorage interface {
	List() ([]*File, error)
	LoadTemplate(FilePath) (*Template, error)
	Save(*Template) (*Template, error)
}

package template

type CreateTemplate struct {
	Name string
	Path string
}

type TemplateService interface {
	Create(CreateTemplate) (*Template, error)
	Delete(FilePath) error
	Edit(FilePath) (*Template, error)
	Find(FilePath) (*Template, error)
	List() ([]*Template, error)
	ListFiles() ([]*File, error)
	Update(*Template) (*Template, error)
}

type TemplateServiceImpl struct {
	fileStorage      FileStorage
	DefaultDirectory string
}

func NewTemplateService(fileStorage FileStorage, defaultDirectory string) TemplateService {
	return &TemplateServiceImpl{
		fileStorage:      fileStorage,
		DefaultDirectory: defaultDirectory,
	}
}

func (t *TemplateServiceImpl) Create(act CreateTemplate) (*Template, error) {
	path := act.Path
	if path == "" {
		path = t.DefaultDirectory
	}

	name := act.Name
	if name == "" {
		name = path
	}

	templ := DefaultTemplate(name, path)
	return t.fileStorage.Save(templ)
}

func (t *TemplateServiceImpl) Delete(path FilePath) error {
	return nil
}

func (t *TemplateServiceImpl) Edit(path FilePath) (*Template, error) {
	return nil, nil
}

func (t *TemplateServiceImpl) Find(path FilePath) (*Template, error) {
	return nil, nil
}

func (t *TemplateServiceImpl) List() ([]*Template, error) {
	files, err := t.fileStorage.List()
	if err != nil {
		return nil, err
	}

	var templates []*Template
	var errors []error
	for _, file := range files {
		templ, err := t.fileStorage.LoadTemplate(file.Path())
		if err != nil {
			errors = append(errors, err)
		} else {
			templates = append(templates, templ)
		}
	}

	// TODO: Log templates failing to load

	return templates, nil
}

func (t *TemplateServiceImpl) ListFiles() ([]*File, error) {
	return t.fileStorage.List()
}

func (t *TemplateServiceImpl) Update(templ *Template) (*Template, error) {
	return nil, nil
}

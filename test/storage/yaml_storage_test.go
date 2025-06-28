package storage_test

import (
	"os"
	"testing"
	"thop/internal/config"
	"thop/internal/storage"
	"thop/internal/types/project"
	"thop/internal/types/template"
	"thop/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_List(t *testing.T) {
	cfg := &config.Config{
		ConfigDir: "/foo/bar",
	}

	t.Run("returns empty list when no templates exist", func(t *testing.T) {
		// given
		fs := new(test.MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()
		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{}, nil).Once()

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		// when
		projects, err := st.List()

		// then
		assert.Nil(t, err)
		assert.Empty(t, projects)
		fs.AssertExpectations(t)
	})

	t.Run("returns list of projects", func(t *testing.T) {
		// given
		fs := new(test.MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()

		dir1 := new(test.MockDirEntry)
		dir1.On("IsDir").Return(true).Once()
		dir1.On("Name").Return("foo").Once()

		dir2 := new(test.MockDirEntry)
		dir2.On("IsDir").Return(true).Once()
		dir2.On("Name").Return("bar").Once()

		dir3 := new(test.MockDirEntry)
		dir3.On("IsDir").Return(true).Once()
		dir3.On("Name").Return("baz").Once()

		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{dir1, dir2, dir3}, nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/foo/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/bar/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/baz/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()

		expectedProjects := []project.Project{
			{UUID: "foo", Name: "foobar", Template: template.Template{Root: "/home/test"}},
			{UUID: "bar", Name: "foobar", Template: template.Template{Root: "/home/test"}},
			{UUID: "baz", Name: "foobar", Template: template.Template{Root: "/home/test"}},
		}

		for i := range expectedProjects {
			expectedProjects[i] = expectedProjects[i].WithDefaults()
		}

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		// when
		projects, err := st.List()

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedProjects, projects)
		fs.AssertExpectations(t)
	})

	t.Run("gracefully handles non directory files", func(t *testing.T) {
		// given
		dir := new(test.MockDirEntry)
		dir.On("IsDir").Return(false).Once()

		fs := new(test.MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()
		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{dir}, nil).Once()

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		// when
		projects, err := st.List()

		// then
		assert.Nil(t, err)
		assert.Empty(t, projects)
		fs.AssertExpectations(t)
	})

	t.Run("gracefully handles error while reading template file", func(t *testing.T) {
		// given
		dir := new(test.MockDirEntry)
		dir.On("IsDir").Return(true).Once()
		dir.On("Name").Return("foo").Once()

		fs := new(test.MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()
		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{dir}, nil).Once()

		bytes := []byte("  name: invalid:format: \ntemplate:\n  root: /home/test\n")
		fs.On("ReadFile", "/foo/bar/templates/foo/template.yaml").Return(bytes, nil).Once()

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		// when
		projects, err := st.List()

		// then
		assert.Nil(t, err)
		assert.Empty(t, projects)
		fs.AssertExpectations(t)
	})
}

func Test_ListWithInvalid(t *testing.T) {
	cfg := &config.Config{
		ConfigDir: "/foo/bar",
	}

	t.Run("returns empty list when no templates exist", func(t *testing.T) {
		// given
		fs := new(test.MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()
		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{}, nil).Once()

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		// when
		projects, err := st.ListWithInvalid()

		// then
		assert.Nil(t, err)
		assert.Empty(t, projects)
		fs.AssertExpectations(t)
	})

	t.Run("returns list of projects with invalid ones", func(t *testing.T) {
		// given
		fs := new(test.MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()

		dir1 := new(test.MockDirEntry)
		dir1.On("IsDir").Return(true).Once()
		dir1.On("Name").Return("foo").Once()

		dir2 := new(test.MockDirEntry)
		dir2.On("IsDir").Return(true).Once()
		dir2.On("Name").Return("bar").Once()

		dir3 := new(test.MockDirEntry)
		dir3.On("IsDir").Return(true).Once()
		dir3.On("Name").Return("baz").Once()

		dir4 := new(test.MockDirEntry)
		dir4.On("IsDir").Return(true).Once()
		dir4.On("Name").Return("invalid").Once()

		dir5 := new(test.MockDirEntry)
		dir5.On("IsDir").Return(true).Once()
		dir5.On("Name").Return("invalid2").Once()

		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{dir1, dir2, dir3, dir4, dir5}, nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/foo/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/bar/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/baz/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/invalid/template.yaml").Return([]byte("template:\n  root: /home/test\n  root: /home/test\n"), nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/invalid2/template.yaml").Return([]byte("name: invalid2\nname: invalid2\n"), nil).Once()

		expectedProjects := []project.Project{
			{UUID: "foo", Name: "foobar", Template: template.Template{Root: "/home/test"}},
			{UUID: "bar", Name: "foobar", Template: template.Template{Root: "/home/test"}},
			{UUID: "baz", Name: "foobar", Template: template.Template{Root: "/home/test"}},
			{UUID: "invalid", Name: "invalid", Type: project.TypeInvalid},
			{UUID: "invalid2", Name: "invalid2", Type: project.TypeInvalid},
		}

		for i := range expectedProjects {
			expectedProjects[i] = expectedProjects[i].WithDefaults()
		}

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		// when
		projects, err := st.ListWithInvalid()

		// then
		assert.Nil(t, err)
		assert.Equal(t, projects, expectedProjects)
		fs.AssertExpectations(t)
	})

	t.Run("gracefully handles non directory files", func(t *testing.T) {
		// given
		dir := new(test.MockDirEntry)
		dir.On("IsDir").Return(false).Once()

		fs := new(test.MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()
		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{dir}, nil).Once()

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		// when
		projects, err := st.ListWithInvalid()

		// then
		assert.Nil(t, err)
		assert.Empty(t, projects)
		fs.AssertExpectations(t)
	})
}

func Test_Find(t *testing.T) {
	cfg := &config.Config{
		ConfigDir: "/foo/bar",
	}

	t.Run("returns error when project does not exist", func(t *testing.T) {
		// given
		fs := new(test.MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()
		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{}, nil).Once()

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		// when
		_, err := st.Find("foobar")

		// then
		assert.True(t, storage.ErrProjectNotFound.Equal(err))
		fs.AssertExpectations(t)
	})

	t.Run("returns found project", func(t *testing.T) {
		// given
		fs := new(test.MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates").Return(nil).Once()

		dir := new(test.MockDirEntry)
		dir.On("IsDir").Return(true).Once()
		dir.On("Name").Return("foo").Once()

		fs.On("ReadDir", "/foo/bar/templates").Return([]os.DirEntry{dir}, nil).Once()
		fs.On("ReadFile", "/foo/bar/templates/foo/template.yaml").Return([]byte("name: foobar\ntemplate:\n  root: /home/test\n"), nil).Once()
		expectedProject := project.Project{UUID: "foo", Name: "foobar", Template: template.Template{Root: "/home/test"}}

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		// when
		project, err := st.Find(expectedProject.Name)
		expectedProject = expectedProject.WithDefaults()

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedProject, project)
		fs.AssertExpectations(t)
	})
}

func Test_Save(t *testing.T) {
	t.Run("saves project with a template file", func(t *testing.T) {
		cfg := &config.Config{
			ConfigDir: "/foo/bar",
		}

		// given
		fs := new(test.MockFileSystem)
		var path string
		fs.On("MkdirAll", mock.Anything).Run(func(args mock.Arguments) {
			path = args.Get(0).(string)
		}).Return(nil).Once()

		var templatePath string
		fs.On("WriteFile", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			templatePath = args.Get(0).(string)
		}).Return(nil).Once()

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}
		project := project.Project{Name: "foo"}

		// when
		err := st.Save(&project)

		// then
		assert.Nil(t, err)

		assert.NotEmpty(t, project.UUID)
		assert.Equal(t, "/foo/bar/templates/"+string(project.UUID), path)
		assert.Equal(t, "/foo/bar/templates/"+string(project.UUID)+"/template.yaml", templatePath)
	})

	t.Run("keeps UUID if it's already set", func(t *testing.T) {
		// given
		cfg := &config.Config{
			ConfigDir: "/foo/bar",
		}

		fs := new(test.MockFileSystem)
		fs.On("MkdirAll", "/foo/bar/templates/foobar").Return(nil).Once()

		fs.On("WriteFile", "/foo/bar/templates/foobar/template.yaml", mock.Anything).Return(nil).Once()

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}
		project := project.Project{UUID: "foobar", Name: "foo"}

		// when
		err := st.Save(&project)

		// then
		assert.Nil(t, err)
		fs.AssertExpectations(t)
	})
}

func Test_Delete(t *testing.T) {
	cfg := &config.Config{
		ConfigDir: "/foo/bar",
	}

	t.Run("deletes template directory", func(t *testing.T) {
		// given
		fs := new(test.MockFileSystem)
		fs.On("RemoveAll", "/foo/bar/templates/foo").Return(nil).Once()

		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		// when
		err := st.Delete("foo")

		// then
		assert.Nil(t, err)
		fs.AssertExpectations(t)
	})
}

func Test_PrepareTemplateFile(t *testing.T) {
	cfg := &config.Config{
		ConfigDir: "/foo/bar",
	}

	t.Run("returns path to template file", func(t *testing.T) {
		// given
		fs := new(test.MockFileSystem)
		st := &storage.YamlStorage{Config: cfg, FileSystem: fs}

		project := project.Project{UUID: "foo"}

		// when
		path, err := st.PrepareTemplateFile(project.UUID)

		// then
		assert.Nil(t, err)
		assert.Equal(t, "/foo/bar/templates/foo/template.yaml", path)
		fs.AssertExpectations(t)
	})
}

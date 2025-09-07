package template

import (
	"errors"
	"testing"
	"thop/internal/domain/template"
	"thop/test/gen/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_TemplateListing(t *testing.T) {
	t.Run("lists templates", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)

		expectedTemplateFiles := []*template.File{
			template.NewFile("1", "foo/bar/1"),
		}

		expectedTemplates := []*template.Template{
			template.NewTemplate(
				"foo/bar/1",
				template.V1,
				"1",
				"",
				"",
				[]string{},
				[]*template.Window{},
			),
		}

		mockFileStorage.EXPECT().List().Return(expectedTemplateFiles, nil)
		mockFileStorage.EXPECT().LoadTemplate(expectedTemplateFiles[0].Path()).Return(expectedTemplates[0], nil)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		templates, err := service.List()

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedTemplates, templates)
	})

	t.Run("ignores template loading errors", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)

		expectedTemplateFiles := []*template.File{
			template.NewFile("1", "foo/bar/1"),
			template.NewFile("2", "foo/bar/2"),
		}

		expectedTemplates := []*template.Template{
			template.NewTemplate(
				"foo/bar/1",
				template.V1,
				"1",
				"",
				"",
				[]string{},
				[]*template.Window{},
			),
		}

		mockFileStorage.EXPECT().
			List().
			Return(expectedTemplateFiles, nil)

		mockFileStorage.EXPECT().
			LoadTemplate(expectedTemplateFiles[0].Path()).
			Return(expectedTemplates[0], nil)

		mockFileStorage.EXPECT().
			LoadTemplate(expectedTemplateFiles[1].Path()).
			Return(nil, errors.New("parse error"))

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		templates, err := service.List()

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedTemplates, templates)
	})

	t.Run("propagates file listing error", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)

		expectedErr := errors.New("read error")

		mockFileStorage.EXPECT().
			List().
			Return(nil, expectedErr)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		templates, err := service.List()

		// then
		assert.Nil(t, templates)
		assert.Equal(t, err, expectedErr)
	})
}

func Test_FileListing(t *testing.T) {
	t.Run("lists files", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)

		expectedTemplateFiles := []*template.File{
			template.NewFile("1", "foo/bar/1"),
			template.NewFile("2", "foo/bar/2"),
		}

		mockFileStorage.EXPECT().List().Return(expectedTemplateFiles, nil)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		files, err := service.ListFiles()

		// then
		assert.Nil(t, err)
		assert.Equal(t, expectedTemplateFiles, files)
	})

	t.Run("propagates file listing error", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileStorage := mock.NewMockFileStorage(ctrl)

		expectedErr := errors.New("read error")

		mockFileStorage.EXPECT().
			List().
			Return(nil, expectedErr)

		service := template.NewTemplateService(
			mockFileStorage,
			"",
		)

		// when
		files, err := service.ListFiles()

		// then
		assert.Nil(t, files)
		assert.Equal(t, err, expectedErr)
	})
}

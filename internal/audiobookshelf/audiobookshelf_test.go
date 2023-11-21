package audiobookshelf_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vpoluyaktov/abb_ia/internal/audiobookshelf"
	"github.com/vpoluyaktov/abb_ia/internal/config"
)

func TestLogin(t *testing.T) {
	config.Load()
	url := config.Instance().GetAudiobookshelfUrl()
	username := config.Instance().GetAudiobookshelfUser()
	password := config.Instance().GetAudiobookshelfPassword()

	if url != "" && username != "" && password != "" {
		absClient := audiobookshelf.NewClient(url)
		err := absClient.Login(username, password)
		assert.NoError(t, err)
	}
}

func TestLibraries(t *testing.T) {
	config.Load()
	url := config.Instance().GetAudiobookshelfUrl()
	username := config.Instance().GetAudiobookshelfUser()
	password := config.Instance().GetAudiobookshelfPassword()

	if url != "" && username != "" && password != "" {
		absClient := audiobookshelf.NewClient(url)
		err := absClient.Login(username, password)
		assert.NoError(t, err)
		if err == nil {
			libraries, err := absClient.GetLibraries()
			assert.NoError(t, err)
			assert.NotNil(t, libraries)
			assert.NotEmpty(t, libraries)
		}
	}
}

func TestScan(t *testing.T) {
	config.Load()
	url := config.Instance().GetAudiobookshelfUrl()
	username := config.Instance().GetAudiobookshelfUser()
	password := config.Instance().GetAudiobookshelfPassword()
	libraryName := config.Instance().GetAudiobookshelfLibrary()

	if url != "" && username != "" && password != "" && libraryName != "" {
		absClient := audiobookshelf.NewClient(url)
		err := absClient.Login(username, password)
		assert.NoError(t, err)
		if err == nil {
			libraries, err := absClient.GetLibraries()
			assert.NoError(t, err)
			libraryID, err := absClient.GetLibraryId(libraries, libraryName)
			if err == nil {
				err = absClient.ScanLibrary(libraryID)
				assert.NoError(t, err)
			}
		}
	}
}

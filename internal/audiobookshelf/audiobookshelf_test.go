package audiobookshelf_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vpoluyaktov/abb_ia/internal/audiobookshelf"
	"github.com/vpoluyaktov/abb_ia/internal/config"
)

func TestLogin(t *testing.T) {
	config.Load()
	url := config.AudiobookshelfUrl()
	username := config.AudiobookshelfUser()
	password := config.AudiobookshelfPassword()

	if url != "" && username != "" && password != "" {
		loginResp, err := audiobookshelf.Login(url+"/login", username, password)

		assert.NoError(t, err)
		assert.NotNil(t, loginResp.User.ID)
		assert.NotNil(t, loginResp.User.Token)
	}
}

func TestLibraries(t *testing.T) {
	config.Load()
	url := config.AudiobookshelfUrl()
	username := config.AudiobookshelfUser()
	password := config.AudiobookshelfPassword()

	if url != "" && username != "" && password != "" {
		loginResp, err := audiobookshelf.Login(url+"/login", username, password)
		assert.NoError(t, err)
		if err == nil {
			libraryResponse , err := audiobookshelf.Libraries(url, loginResp.User.Token)
			assert.NoError(t, err)
			assert.NotNil(t, libraryResponse)
			assert.NotEmpty(t, libraryResponse.Libraries)
		}
	}
}

func TestScan(t *testing.T) {
	config.Load()
	url := config.AudiobookshelfUrl()
	username := config.AudiobookshelfUser()
	password := config.AudiobookshelfPassword()
	libraryName := config.AudiobookshelfLibrary()

	if url != "" && username != "" && password != ""  && libraryName != "" {
		loginResp, err := audiobookshelf.Login(url+"/login", username, password)
		if err == nil {
			assert.NoError(t, err)
			libraryResponse , err := audiobookshelf.Libraries(url, loginResp.User.Token)
			if err == nil {
				assert.NoError(t, err)
				libraryID, err := audiobookshelf.GetLibraryByName(libraryResponse.Libraries, libraryName)
				if err == nil {
					err = audiobookshelf.ScanLibrary(url, loginResp.User.Token, libraryID)
					assert.NoError(t, err)
					assert.NotNil(t, libraryResponse)
					assert.NotEmpty(t, libraryResponse.Libraries)
				}
			}
		}
	}
}
package ia_client_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"abb_ia/internal/config"
	"abb_ia/internal/ia"
	"abb_ia/internal/logger"
	"github.com/stretchr/testify/assert"
)

const (
	logFileName string = "/tmp/abb_ia.test.log"
	logLevel           = logger.DEBUG
)

func TestMain(m *testing.M) {
	logger.Init(logFileName, "DEBUG")
	config.Load()
	os.Exit(m.Run())
}

func TestSearch(t *testing.T) {
	ia := ia_client.New(5, false, false)

	res := ia.Search("Old Time Radio Researchers", "Single Episodes", "audio", "date", "asc") // search by author and title
	assert.NotNil(t, res)
	assert.Equal(t, 5, len(res.Response.Docs))

	res = ia.Search("Old Time Radio Researchers", "", "audio", "date", "asc") // search by author only
	assert.NotNil(t, res)
	assert.Equal(t, 5, len(res.Response.Docs))

	res = ia.Search("", "Single Episodes", "audio", "date", "asc") // search by title only
	assert.NotNil(t, res)
	assert.Equal(t, 5, len(res.Response.Docs))

	res = ia.Search("", "https://archive.org/details/OTRR_Frank_Race_Singles", "audio", "date", "asc") // search by item ID
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res.Response.Docs))
}

func TestGetItemById(t *testing.T) {
	ia := ia_client.New(5, false, false)
	item := ia.GetItemDetails("OTRR_Frank_Race_Singles")
	assert.NotNil(t, item)
	assert.GreaterOrEqual(t, 1, len(item.Metadata.Title))
	assert.GreaterOrEqual(t, 1, len(item.Metadata.Creator))
	if logLevel == logger.DEBUG {
		if item != nil {
			fmt.Printf("Title: %s\n", item.Metadata.Title[0])
			fmt.Printf("Server: %s\n", item.Server)
			fmt.Printf("Directory: %s\n", item.Dir)
			// fmt.Printf("Description: %s\n", ia.Html2Text(item.Metadata.Description[0]))
			fmt.Printf("Creator: %s\n", item.Metadata.Creator[0])
			fmt.Printf("Image: %s\n", item.Misc.Image)

			// for file, meta := range item.Files {
			// fmt.Printf("%s -> %s\n", file, meta.Format)
			// }
		}
	}
}

func TestDownloadItem(t *testing.T) {

	server := "ia800303.us.archive.org"
	dir := "/21/items/OTRR_Frank_Race_Singles"
	file := "/Frank_Race_49-xx-xx_ep13_The_Adventure_Of_The_Garrulous_Bartender_spectrogram.png"
	outputDir := "/tmp/abb_ia"

	ia := ia_client.New(5, false, false)
	ia.DownloadFile(outputDir, filepath.Join(outputDir, dir, file), server, dir, file, 1, 1024, UpdateProgress)
}

func UpdateProgress(fileId int, fileName string, size int64, pos int64, percent int) {
	fmt.Printf("Downloading... %d%%\n", percent)
}

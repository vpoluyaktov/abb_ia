package ia_client_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/ia_client"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
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
	ia := ia_client.New(false, false)

	res := ia.Search("Single Episodes", "audio") // search by title
	assert.NotNil(t, res)
	assert.Equal(t, 25, len(res.Response.Docs))

	res = ia.Search("https://archive.org/details/OTRR_Frank_Race_Singles", "audio") // search by item ID
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res.Response.Docs))
}

func TestGetItemById(t *testing.T) {
	ia := ia_client.New(false, false)
	item := ia.GetItemDetails("OTRR_Frank_Race_Singles")
	assert.NotNil(t, item)
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

	ia := ia_client.New(false, false)
	ia.DownloadFile(outputDir, server, dir, file, 1, 1024, UpdateProgress)
}

func UpdateProgress(fileId int, fileName string, size int64, pos int64, percent int) {
	fmt.Printf("Downloading... %d%%\n", percent)
}

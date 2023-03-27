package ia_client_test

import (
	"os"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/vpoluyaktov/audiobook_creator_IA/pkg/ia_client"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/config"
	"github.com/vpoluyaktov/audiobook_creator_IA/internal/logger"
)


const (
  logFileName string = "/tmp/audiobook_creator_IA.test.log"
  logLevel = logger.DEBUG
)

func TestMain(m *testing.M) {
  logger.Init(logFileName, logLevel)
  config.Load()
	os.Exit(m.Run())
}

func TestSearch(t *testing.T) {
	ia := ia_client.New()

  res := ia.Search("NASA", "audio")  // search by title
	assert.NotNil(t, res)
	assert.Equal(t, 25, len(res.Response.Docs))

	res = ia.Search("https://archive.org/details/OTRR_Frank_Race_Singles", "audio") // search by item ID
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res.Response.Docs))
}

func TestGetItemById(t *testing.T) {
	ia := ia_client.New()
	res := ia.GetItemById("OTRR_Frank_Race_Singles")
	assert.NotNil(t, res)
	if logLevel == logger.DEBUG {
		for file, meta := range res.Files {
			fmt.Printf("%s -> %s\n", file, meta.Format)
		} 
	}
}
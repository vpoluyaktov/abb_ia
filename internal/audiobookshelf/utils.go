package audiobookshelf

import (
	"abb_ia/internal/utils"
	"path/filepath"
)

// Calculate Audiobookshelf directory structure (see: https://www.audiobookshelf.org/docs#book-directory-structure)
func GetDestignationPath(outputDir string, series string, author string) string {
	destPath := filepath.Join(outputDir, author)
	if series != "" {
		destPath = filepath.Join(destPath, author+" - "+series)
	}
	return utils.SanitizeFilePath(destPath)
}

func GetDestignationDir(series string, seriesNo string, title string, narrator string) string {
	abTitle := ""
	if series != "" && seriesNo != "" {
		abTitle = seriesNo + ". "
	}
	abTitle += title
	if narrator != "" {
		abTitle += " {" + narrator + "}"
	}
	return utils.SanitizeFilePath(abTitle)
}

package controller

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/ffmpeg"
	"github.com/vpoluyaktov/abb_ia/internal/utils"

	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type BuildController struct {
	mq        *mq.Dispatcher
	ab        *dto.Audiobook
	startTime time.Time
	stopFlag  bool

	// progress tracking arrays
	filesBuild []fileBuild
	filesCopy  []fileCopy
}

type fileBuild struct {
	fileId   int
	progress int
}

type fileCopy struct {
	fileId      int
	fileSize    int64
	bytesCopied int64
	progress    int
}

// Progress Reader for file copy progress
type Fn func(fileId int, fileName string, size int64, pos int64, percent int)
type ProgressReader struct {
	FileId   int
	FileName string
	Reader   io.Reader
	Size     int64
	Pos      int64
	Percent  int
	Callback Fn
}

func NewBuildController(dispatcher *mq.Dispatcher) *BuildController {
	dc := &BuildController{}
	dc.mq = dispatcher
	dc.mq.RegisterListener(mq.BuildController, dc.dispatchMessage)
	return dc
}

func (c *BuildController) checkMQ() {
	m := c.mq.GetMessage(mq.BuildController)
	if m != nil {
		c.dispatchMessage(m)
	}
}

func (c *BuildController) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.BuildCommand:
		go c.startBuild(dto)
	case *dto.StopCommand:
		go c.stopBuild(dto)
	default:
		m.UnsupportedTypeError(mq.BuildController)
	}
}

func (c *BuildController) stopBuild(cmd *dto.StopCommand) {
	c.stopFlag = true
	logger.Debug(mq.BuildController + ": Received StopBuild command")
}

func (c *BuildController) startBuild(cmd *dto.BuildCommand) {
	c.startTime = time.Now()
	logger.Info(mq.BuildController + " received " + cmd.String())

	c.ab = cmd.Audiobook

	// calculate output file names
	for i := range c.ab.Parts {
		part := &c.ab.Parts[i]
		filePath := filepath.Join(config.OutputDir(), c.ab.Author+" - "+c.ab.Title)
		if len(c.ab.Parts) > 1 {
			filePath = filePath + fmt.Sprintf(", Part %02d", i+1)
		}
		part.AACFile = filePath + ".aac"
		part.M4BFile = filePath + ".m4b"
	}

	c.mq.SendMessage(mq.BuildController, mq.BuildPage, &dto.DisplayBookInfoCommand{Audiobook: c.ab}, true)
	c.mq.SendMessage(mq.BuildController, mq.Footer, &dto.UpdateStatus{Message: "Building audiobook..."}, false)
	c.mq.SendMessage(mq.BuildController, mq.Footer, &dto.SetBusyIndicator{Busy: true}, false)

	// prepare .mp3 file list
	c.createFilesLists(c.ab)
	// prepare metadata
	c.createMetadata(c.ab)
	c.downloadCoverImage(c.ab)

	// build audiobook parts
	c.stopFlag = false
	c.filesBuild = make([]fileBuild, len(c.ab.Parts))
	jd := utils.NewJobDispatcher(config.ParallelEncoders())
	for i := range c.ab.Parts {
		jd.AddJob(i, c.buildAudiobookPart, c.ab, i)
	}
	go c.updateTotalBuildProgress()
	// if c.stopFlag {
	// 	break
	// }
	jd.Start()

	// copy book to Audiobookshelf 
	if config.IsCopyToAudiobookshelf() {
		c.filesCopy = make([]fileCopy, len(c.ab.Parts))
		jd := utils.NewJobDispatcher(1)
		for i := range c.ab.Parts {
			jd.AddJob(i, c.copyAudiobookPart, c.ab, i)
		}
		go c.updateTotalCopyProgress()
		jd.Start()
	}

	c.cleanUp(c.ab)

	c.mq.SendMessage(mq.BuildController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.BuildController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
	c.mq.SendMessage(mq.BuildController, mq.BuildPage, &dto.BuildComplete{Audiobook: cmd.Audiobook}, true)
}

func (c *BuildController) createFilesLists(ab *dto.Audiobook) {
	for i := range ab.Parts {
		part := &ab.Parts[i]
		part.FListFile = filepath.Join(ab.OutputDir, fmt.Sprintf("Part %02d Files.txt", part.Number))
		f, err := os.OpenFile(part.FListFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			logger.Error("Can't open FList file for writing: " + err.Error())
		}
		for _, chapter := range part.Chapters {
			for _, file := range chapter.Files {
				f.WriteString("file '" + strings.TrimPrefix(file.FileName, "/") + "'\n")
			}
		}
		f.Close()
	}
}

func (c *BuildController) createMetadata(ab *dto.Audiobook) {
	for i := range ab.Parts {
		part := &ab.Parts[i]
		part.MetadataFile = filepath.Join(ab.OutputDir, fmt.Sprintf("Part %02d Metadata.txt", part.Number))
		f, err := os.OpenFile(part.MetadataFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			logger.Error("Can't open Metadata file for writing: " + err.Error())
		}
		f.WriteString(";FFMETADATA1\n")
		f.WriteString("major_brand=isom\n")
		f.WriteString("minor_version=512\n")
		f.WriteString("compatible_brands=isomiso2mp41\n")
		f.WriteString("title=" + ab.Title + "\n")
		f.WriteString("artist=" + ab.Author + "\n")
		f.WriteString("album=" + ab.Title + "\n")
		f.WriteString("genre=Audiobook" + ab.Genre + "\n") //TODO: remove default genre
		f.WriteString("description=" + strings.ReplaceAll(ab.Description, "\n", "\\\n") + "\n")
		f.WriteString("copyright=" + ab.Copyright + "\n")
		f.WriteString("comment=Downloaded from Internet Archive: " + ab.IaURL + "\n")
		f.WriteString("encoder=This audiobook was created by 'Audiobook Builder Internet Archive version' https://github.com/vpoluyaktov/abb_ia\n")

		for _, chapter := range part.Chapters {
			f.WriteString("[CHAPTER]\n")
			f.WriteString("TIMEBASE=1/1000\n")
			f.WriteString("START=" + strconv.Itoa(int(chapter.Start*1000)) + "\n")
			f.WriteString("END=" + strconv.Itoa(int(chapter.End*1000)) + "\n")
			f.WriteString("title=" + chapter.Name + "\n")
		}
		f.Close()
	}
}

func (c *BuildController) downloadCoverImage(ab *dto.Audiobook) error {
	filePath := filepath.Join(config.OutputDir(), ab.Author+" - "+ab.Title)
	if strings.HasSuffix(ab.CoverURL, ".jpg") {
		ab.CoverFile = filePath + ".jpg"
	} else if strings.HasSuffix(ab.CoverURL, ".png") {
		ab.CoverFile = filePath + ".png"
	} else {
		logger.Debug("Wrong image type: " + ab.CoverURL)
	}
	response, err := http.Get(ab.CoverURL)
	if err != nil {
		logger.Error("Can't download cover image: " + ab.CoverURL + ": " + err.Error())
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(ab.CoverFile)
	if err != nil {
		logger.Error("Can't create a file for cover image: " + ab.CoverURL + ": " + err.Error())
		return err
	}
	defer file.Close()

	// Copy the response body to the output file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		logger.Error("Can't save cover image: " + ab.CoverURL + ": " + err.Error())
		return err
	}
	return nil
}

func (c *BuildController) buildAudiobookPart(ab *dto.Audiobook, partId int) {
	part := &ab.Parts[partId]

	// launch progress listener
	l, port := c.startProgressListener(partId)
	defer l.Close()
	go c.updateFileBuildProgress(partId, part.M4BFile, part.Duration, l)

	// concatenate mp3 files into single .aac file
	_, err := ffmpeg.NewFFmpeg().
		Input(part.FListFile, "-safe 0 -f concat").
		Output(part.AACFile, fmt.Sprintf("-acodec aac -ab %s -ar %s -vn", config.BitRate(), config.SampleRate())).
		Overwrite(true).
		Params("-hide_banner -nostdin -nostats").
		SendProgressTo("http://127.0.0.1:" + strconv.Itoa(port)).
		Run()
	if err != nil {
		logger.Error("FFMPEG Error: " + string(err.Error()))
	} else {
		// add Metadata, cover image and convert to .m4b
		_, err := ffmpeg.NewFFmpeg().
			Input(part.AACFile, "").
			Input(part.MetadataFile, "").
			Input(ab.CoverURL, "").
			Output(part.M4BFile, "-map_metadata 1 -y -acodec copy -y -vf pad='width=ceil(iw/2)*2:height=ceil(ih/2)*2'").
			Overwrite(true).
			Params("-hide_banner -nostdin -nostats").
			SendProgressTo("http://127.0.0.1:" + strconv.Itoa(port)).
			Run()
		if err != nil {
			logger.Error("FFMPEG Error: " + string(err.Error()))
		}
	}
}

func (c *BuildController) copyAudiobookPart(ab *dto.Audiobook, partId int) {

	part := &ab.Parts[partId]

	file, err := os.Open(part.M4BFile)
	if err != nil {
		logger.Error("Can't open .mb4 file: " + err.Error())
		return
	}
	fileReader := bufio.NewReader(file)
	defer file.Close()

	destPath := filepath.Clean(filepath.Join(config.AudiobookshelfDir(), ab.Author, ab.Title, ab.Series, filepath.Base(part.M4BFile)))
  destDir := filepath.Dir(destPath)

	if err := os.MkdirAll(destDir, 0750); err != nil {
		logger.Error("Can't create output directory: " + err.Error())
		return
	}
	f, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatal("Can't create temporary file: " + err.Error())
		return
	}
	defer f.Close()

	progressReader := &ProgressReader{
		FileId:   partId,
		FileName: part.M4BFile,
		Reader:   fileReader,
		Size:     part.Size,
		Callback: c.updateFileCopyProgress,
	}

	if _, err := io.Copy(f, progressReader); err != nil {
		logger.Error("Error while copying .m4b file: " + err.Error())
	}
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if err == nil {
		pr.Pos += int64(n)
		pr.Percent = int(float64(pr.Pos) / float64(pr.Size) * 100)
		pr.Callback(pr.FileId, pr.FileName, pr.Size, pr.Pos, pr.Percent)
	}
	return n, err
}

func (c *BuildController) startProgressListener(fileId int) (net.Listener, int) {

	basePortNumber := 31000
	portNumber := basePortNumber + fileId

	l, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(portNumber))
	if err != nil {
		logger.Error("Build progress: Start listener error: " + err.Error())
	}
	return l, portNumber
}

func (c *BuildController) updateFileBuildProgress(fileId int, fileName string, totalDuration float64, l net.Listener) {

	re := regexp.MustCompile(`out_time_ms=(\d+)`)
	fd, err := l.Accept()
	if err != nil {
		return // listener is closed
	}
	buf := make([]byte, 16)
	data := ""
	percent := 0
	for {
		_, err := fd.Read(buf)
		if err != nil {
			return // listener is closed
		}
		data += string(buf)
		a := re.FindAllStringSubmatch(data, -1)
		p := 0
		pstr := ""
		if len(a) > 0 && len(a[len(a)-1]) > 0 {
			c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
			pstr = fmt.Sprintf("%.2f", float64(c)/totalDuration/1000000)
		}
		if strings.Contains(data, "progress=end") {
			p = 100
		}
		if pstr == "" {
			p = 0
		}
		pflt, err := strconv.ParseFloat(pstr, 64)
		if err != nil {
			p = 0
		} else {
			p = int(pflt * 100)
		}

		if p != percent {
			percent = p
			if percent > 100 {
				percent = 100
			}
			// sent a message only if progress changed
			c.mq.SendMessage(mq.BuildController, mq.BuildPage, &dto.BuildFileProgress{FileId: fileId, FileName: fileName, Percent: percent}, true)
		}
		c.filesBuild[fileId].fileId = fileId
		c.filesBuild[fileId].progress = percent
	}
}

func (c *BuildController) updateTotalBuildProgress() {
	var percent int = -1

	for !c.stopFlag && percent <= 100 {
		var totalPercent int = 0
		files := 0
		for _, f := range c.filesBuild {
			totalPercent += f.progress
			if f.progress == 100 {
				files++
			}
		}
		p := int(totalPercent / len(c.filesBuild))

		if percent != p {
			// sent a message only if progress changed
			percent = p

			elapsed := time.Since(c.startTime).Seconds()
			speed := int64(float64(percent) / elapsed)
			eta := (100 / (float64(percent) / elapsed)) - elapsed
			if eta < 0 || eta > (60*60*24*365) {
				eta = 0
			}

			elapsedH := utils.SecondsToTime(elapsed)
			filesH := fmt.Sprintf("%d/%d", files, len(c.ab.Parts))
			speedH := utils.SpeedToHuman(speed)
			etaH := utils.SecondsToTime(eta)

			c.mq.SendMessage(mq.BuildController, mq.BuildPage, &dto.BuildProgress{Elapsed: elapsedH, Percent: percent, Files: filesH, Speed: speedH, ETA: etaH}, true)
		}
		time.Sleep(mq.PullFrequency)
	}
}

func (c *BuildController) updateFileCopyProgress(fileId int, fileName string, size int64, pos int64, percent int) {
	if c.filesBuild[fileId].progress != percent {
		// sent a message only if progress changed
		c.mq.SendMessage(mq.BuildController, mq.BuildPage, &dto.CopyFileProgress{FileId: fileId, FileName: fileName, Percent: percent}, false)
	}
	c.filesCopy[fileId].fileId = fileId
	c.filesCopy[fileId].fileSize = size
	c.filesCopy[fileId].bytesCopied = pos
	c.filesCopy[fileId].progress = percent
}

func (c *BuildController) updateTotalCopyProgress() {
	var percent int = -1

	item := c.ab.IAItem
	for !c.stopFlag && percent <= 100 {
		var totalSize = item.TotalSize
		var totalBytesDownloaded int64 = 0
		filesDownloaded := 0
		for _, f := range c.filesCopy {
			totalBytesDownloaded += f.bytesCopied
			if f.progress == 100 {
				filesDownloaded++
			}
		}

		var p int = 0
		if totalSize > 0 {
			p = int(float64(totalBytesDownloaded) / float64(totalSize) * 100)
		}

		// fix wrong file size returned by IA metadata
		if filesDownloaded == len(c.filesBuild) {
			p = 100
			totalBytesDownloaded = item.TotalSize
		}

		if percent != p {
			// sent a message only if progress changed
			percent = p

			elapsed := time.Since(c.startTime).Seconds()
			speed := int64(float64(totalBytesDownloaded) / elapsed)
			eta := (100 / (float64(percent) / elapsed)) - elapsed
			if eta < 0 || eta > (60*60*24*365) {
				eta = 0
			}

			elapsedH := utils.SecondsToTime(elapsed)
			bytesH := utils.BytesToHuman(totalBytesDownloaded)
			filesH := fmt.Sprintf("%d/%d", filesDownloaded, len(item.Files))
			speedH := utils.SpeedToHuman(speed)
			etaH := utils.SecondsToTime(eta)

			c.mq.SendMessage(mq.BuildController, mq.BuildPage, &dto.CopyProgress{Elapsed: elapsedH, Percent: percent, Files: filesH, Bytes: bytesH, Speed: speedH, ETA: etaH}, false)
		}
		time.Sleep(mq.PullFrequency)
	}
}

func (c *BuildController) cleanUp(ab *dto.Audiobook) {
	if !(config.IsSaveMock() || config.IsUseMock()) {
		os.RemoveAll(ab.OutputDir)
	}
	for _, part := range ab.Parts {
		os.Remove(part.AACFile)
		os.Remove(part.FListFile)
		os.Remove(part.MetadataFile)
	}
	os.Remove(ab.CoverFile)
}

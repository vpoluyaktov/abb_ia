package controller

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"abb_ia/internal/dto"
	"abb_ia/internal/ffmpeg"
	"abb_ia/internal/utils"

	"abb_ia/internal/logger"
	"abb_ia/internal/mq"
)

type BuildController struct {
	mq        *mq.Dispatcher
	ab        *dto.Audiobook
	startTime time.Time
	stopFlag  bool
	files     []fileBuild
}

// progress tracking arrays
type fileBuild struct {
	fileName         string
	totalDuration    float64
	bytesProcessed   int64
	secondsProcessed float64
	encodingSpeed    float64
	progress         int
	complete         bool
}

func NewBuildController(dispatcher *mq.Dispatcher) *BuildController {
	c := &BuildController{}
	c.mq = dispatcher
	c.mq.RegisterListener(mq.BuildController, c.dispatchMessage)
	return c
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
	logger.Info(mq.BuildController + " received " + cmd.String())
	c.stopFlag = false
	c.startTime = time.Now()
	c.ab = cmd.Audiobook
	c.files = make([]fileBuild, len(c.ab.Parts))

	// calculate output file names
	for i := range c.ab.Parts {
		part := &c.ab.Parts[i]
		filePath := filepath.Join(c.ab.Config.GetTmpDir(), c.ab.Author+" - "+c.ab.Title)
		if len(c.ab.Parts) > 1 {
			filePath = filePath + fmt.Sprintf(", Part %02d", i+1)
		}
		part.AACFile = filePath + ".aac"
		part.M4BFile = filePath + ".m4b"
		c.files[i].fileName = part.M4BFile
		c.files[i].totalDuration = part.Duration
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

	jd := utils.NewJobDispatcher(c.ab.Config.GetConcurrentEncoders())
	for i := range c.ab.Parts {
		jd.AddJob(i, c.buildAudiobookPart, c.ab, i)
	}
	go c.updateTotalProgress()
	jd.Start()

	c.mq.SendMessage(mq.BuildController, mq.Footer, &dto.SetBusyIndicator{Busy: false}, false)
	c.mq.SendMessage(mq.BuildController, mq.Footer, &dto.UpdateStatus{Message: ""}, false)
	if !c.stopFlag {
		c.mq.SendMessage(mq.BuildController, mq.BuildPage, &dto.BuildComplete{Audiobook: cmd.Audiobook}, true)
	}
	c.stopFlag = true
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
		f.WriteString("genre=" + ab.Genre + "\n")
		f.WriteString("description=" + strings.ReplaceAll(ab.Description, "\n", "\\\n") + "\n")
		f.WriteString("copyright=" + ab.LicenseUrl + "\n")
		f.WriteString("comment=This audiobook was created using the 'Audiobook Builder' tool: https://abb_ia\\\n" +
			"The audio files used for this book were obtained from the Internet Archive site: " + ab.IaURL + "\n")

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
	filePath := filepath.Join(ab.Config.GetTmpDir(), ab.Author+" - "+ab.Title)
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
	if c.stopFlag {
		return
	}

	part := &ab.Parts[partId]

	// launch progress listener
	l, port := c.startProgressListener(partId)
	defer l.Close()
	go c.updateFileProgress(partId, l)

	// concatenate mp3 files into single .aac file
	_, err := ffmpeg.NewFFmpeg().
		Input(part.FListFile, "-safe 0 -f concat").
		Output(part.AACFile, "-acodec aac -vn").
		Overwrite(true).
		Params("-hide_banner -nostdin -nostats").
		SendProgressTo("http://127.0.0.1:" + strconv.Itoa(port)).
		Run()
	if err != nil {
		logger.Error("FFMPEG Error: " + string(err.Error()))
	} else {
		// add Metadata, cover image and convert to .m4b
		ffmpeg := ffmpeg.NewFFmpeg().
			Input(part.AACFile, "").
			Input(part.MetadataFile, "").
			Input(ab.CoverURL, "").
			Output(part.M4BFile, "-map_metadata 1 -y -acodec copy -y -vf pad='width=ceil(iw/2)*2:height=ceil(ih/2)*2'").
			Overwrite(true).
			Params("-hide_banner -nostdin -nostats").
			SendProgressTo("http://127.0.0.1:" + strconv.Itoa(port))

		go c.killSwitch(ffmpeg)
		_, err := ffmpeg.Run()
		if err != nil && !c.stopFlag {
			logger.Error("FFMPEG Error: " + string(err.Error()))
		}
	}
}

func (c *BuildController) killSwitch(ffmpeg *ffmpeg.FFmpeg) {
	for !c.stopFlag {
		time.Sleep(mq.PullFrequency)
	}
	ffmpeg.Kill()
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

func (c *BuildController) updateFileProgress(fileId int, l net.Listener) {
	fd, err := l.Accept()
	if err != nil {
		return // listener may be closed already
	}
	buf := make([]byte, 16)
	data := ""
	percent := 0

	for !c.stopFlag {
		_, err := fd.Read(buf)
		if err != nil {
			return // listener is closed
		}
		data += string(buf)
		bytesProcessed, secondsProcessed, encodingSpeed, complete := ffmpeg.ParseFFMPEGProgress(data)
		percent = int(secondsProcessed / c.files[fileId].totalDuration * 100)
		// wrong calculation protection
		if percent < 0 {
			percent = 0
		} else if percent > 100 {
			percent = 100
		} else if percent < c.files[fileId].progress {
			percent = c.files[fileId].progress
		} else if complete {
			percent = 100
		}

		// sent a message only if progress changed
		if percent != c.files[fileId].progress {
			c.files[fileId].bytesProcessed = bytesProcessed
			c.files[fileId].secondsProcessed = secondsProcessed
			c.files[fileId].encodingSpeed = encodingSpeed
			c.files[fileId].progress = percent
			c.files[fileId].complete = complete
			c.mq.SendMessage(mq.BuildController, mq.BuildPage, &dto.BuildFileProgress{FileId: fileId, FileName: c.files[fileId].fileName, Percent: percent}, true)
		}
	}
}

func (c *BuildController) updateTotalProgress() {
	var p int = -1

	for !c.stopFlag && p <= 100 {
		var totalDuration float64 = 0
		var secondsProcessed float64 = 0
		var totalSpeed float64 = 0
		filesProcessed := 0
		filesComplete := 0
		for _, f := range c.files {
			totalDuration += f.totalDuration
			secondsProcessed += f.secondsProcessed
			totalSpeed += f.encodingSpeed
			if f.complete {
				filesComplete++
			}
			if f.encodingSpeed > 0 {
				filesProcessed++
			}
		}
		percent := int(secondsProcessed / totalDuration * 100)
		// wrong calculation protection
		if percent < 0 {
			percent = 0
		} else if percent > 100 {
			percent = 100
		}

		if percent != p {
			// sent a message only if progress changed
			p = percent

			elapsed := time.Since(c.startTime).Seconds()
			speed := totalSpeed / float64(filesProcessed)
			eta := (100 / (float64(percent) / elapsed)) - elapsed
			if eta < 0 || eta > (60*60*24*365) {
				eta = 0
			}

			elapsedH := utils.SecondsToTime(elapsed)
			filesH := fmt.Sprintf("%d/%d", filesComplete, len(c.ab.Parts))
			speedH := fmt.Sprintf("%.0fx", speed)
			etaH := utils.SecondsToTime(eta)

			c.mq.SendMessage(mq.BuildController, mq.BuildPage, &dto.BuildProgress{Elapsed: elapsedH, Percent: percent, Files: filesH, Speed: speedH, ETA: etaH}, true)
		}
		time.Sleep(mq.PullFrequency)
	}
}

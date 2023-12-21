package ui

import (
	"strings"

	"abb_ia/internal/config"
	"abb_ia/internal/dto"
	"abb_ia/internal/logger"
	"abb_ia/internal/utils"

	"github.com/vpoluyaktov/tview"

	"abb_ia/internal/mq"
)

type ConfigPage struct {
	mq       *mq.Dispatcher
	mainGrid *grid

	configCopy    config.Config
	configSection *grid
	buildSection  *grid
	absSection    *grid

	// Audobookbuilder config section
	logFileNameField *tview.InputField
	logLevelField    *tview.DropDown
	searchCondition  *tview.InputField
	rowsPerPage      *tview.InputField
	useMockField     *tview.Checkbox
	saveMockField    *tview.Checkbox
	outputDir        *tview.InputField
	copyToOutputDir  *tview.Checkbox
	tmpDir           *tview.InputField

	// audiobook build config section
	concurrentDownloaders *tview.InputField
	concurrentEncoders    *tview.InputField
	reEncodeFiles         *tview.Checkbox
	bitRate               *tview.InputField
	sampleRate            *tview.InputField
	maxFileSize           *tview.InputField
	shortenTitles         *tview.Checkbox

	// audiobookshelf config section
	uploadToAudiobookshelf *tview.Checkbox
	audiobookshelfUrl      *tview.InputField
	audiobookshelfUser     *tview.InputField
	audiobookshelfPassword *tview.InputField
	audiobookshelfLibrary  *tview.InputField
	scanAudiobookshelf     *tview.Checkbox

	saveConfigButton *tview.Button
	cancelButton     *tview.Button
}

func newConfigPage(dispatcher *mq.Dispatcher) *ConfigPage {
	p := &ConfigPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.ConfigPage, p.dispatchMessage)

	p.mainGrid = newGrid()
	p.mainGrid.SetRows(-1, -1, -1)
	p.mainGrid.SetColumns(0)

	// Audobookbuilder config section
	p.configSection = newGrid()
	p.configSection.SetColumns(-2, -2, -1)
	p.configSection.SetBorder(true)
	p.configSection.SetTitle(" Audiobook Builder Configuration: ")
	p.configSection.SetTitleAlign(tview.AlignLeft)

	configFormLeft := newForm()
	configFormLeft.SetHorizontal(false)
	p.searchCondition = configFormLeft.AddInputField("Default Search condition:", "", 40, nil, func(t string) { p.configCopy.SetSearchCondition(t) })
	p.rowsPerPage = configFormLeft.AddInputField("Rows per a page in the search result:", "", 4, acceptInt, func(t string) { p.configCopy.SetRowsPerPage(utils.ToInt(t)) })
	p.useMockField = configFormLeft.AddCheckbox("Use mock?", false, func(t bool) { p.configCopy.SetUseMock(t) })
	p.saveMockField = configFormLeft.AddCheckbox("Save mock?", false, func(t bool) { p.configCopy.SetSaveMock(t) })
	p.configSection.AddItem(configFormLeft.Form, 0, 0, 1, 1, 0, 0, true)

	configFormRight := newForm()
	configFormRight.SetHorizontal(false)
	p.outputDir = configFormRight.AddInputField("Output directory:", "", 40, nil, func(t string) { p.configCopy.SetOutputdDir(t) })
	p.copyToOutputDir = configFormRight.AddCheckbox("Copy audiobook to the output directory?", false, func(t bool) { p.configCopy.SetCopyToOutputDir(t) })
	p.tmpDir = configFormRight.AddInputField("Temporary (working) directory:", "", 40, nil, func(t string) { p.configCopy.SetTmpDir(t) })
	p.logFileNameField = configFormRight.AddInputField("Log file name:", "", 40, nil, func(t string) { p.configCopy.SetLogfileName(t) })
	p.logLevelField = configFormRight.AddDropdown("Log level:", utils.AddSpaces(logger.LogLeves()), 1, func(o string, i int) { p.configCopy.SetLogLevel(strings.TrimSpace(o)) })
	p.configSection.AddItem(configFormRight.Form, 0, 1, 1, 1, 0, 0, true)

	buttonsForm := newForm()
	buttonsForm.SetHorizontal(false)
	buttonsForm.SetButtonsAlign(tview.AlignRight)
	p.saveConfigButton = buttonsForm.AddButton("Save Settings", p.SaveConfig)
	p.cancelButton = buttonsForm.AddButton("Cancel", p.Cancel)
	p.configSection.AddItem(buttonsForm.Form, 0, 2, 1, 1, 0, 0, false)

	p.mainGrid.AddItem(p.configSection.Grid, 0, 0, 1, 1, 0, 0, true)

	// audiobook build configuration section
	p.buildSection = newGrid()
	p.buildSection.SetColumns(-1, -1)
	p.buildSection.SetBorder(true)
	p.buildSection.SetTitle(" Audiobook Build Configuration: ")
	p.buildSection.SetTitleAlign(tview.AlignLeft)

	buildFormLeft := newForm()
	buildFormLeft.SetHorizontal(false)
	p.concurrentDownloaders = buildFormLeft.AddInputField("Concurrent Downloaders:", "", 4, acceptInt, func(t string) { p.configCopy.SetConcurrentDownloaders(utils.ToInt(t)) })
	p.concurrentEncoders = buildFormLeft.AddInputField("Concurrent Encoders:", "", 4, acceptInt, func(t string) { p.configCopy.SetConcurrentEncoders(utils.ToInt(t)) })
	p.reEncodeFiles = buildFormLeft.AddCheckbox("Re-encode .mp3 files to the same Bit Rate?", false, func(t bool) { p.configCopy.SetReEncodeFiles(t) })
	p.bitRate = buildFormLeft.AddInputField("Bit Rate (Kbps):", "", 4, acceptInt, func(t string) { p.configCopy.SetBitRate(utils.ToInt(t)) })
	p.sampleRate = buildFormLeft.AddInputField("Sample Rate (Hz):", "", 6, acceptInt, func(t string) { p.configCopy.SetSampleRate(utils.ToInt(t)) })
	p.buildSection.AddItem(buildFormLeft.Form, 0, 0, 1, 1, 0, 0, true)

	buildFormRight := newForm()
	buildFormRight.SetHorizontal(false)
	p.maxFileSize = buildFormRight.AddInputField("Audiobook part max file size (Mb):", "", 6, acceptInt, func(t string) { p.configCopy.SetMaxFileSizeMb(utils.ToInt(t)) })
	p.shortenTitles = buildFormRight.AddCheckbox("Shorten titles (for ex. Old Time Radio -> OTRR)?", false, func(t bool) { p.configCopy.SetShortenTitles(t) })
	p.buildSection.AddItem(buildFormRight.Form, 0, 1, 1, 1, 0, 0, true)

	p.mainGrid.AddItem(p.buildSection.Grid, 1, 0, 1, 1, 0, 0, true)

	// audiobookshelf config section
	p.absSection = newGrid()
	p.absSection.SetColumns(-1)
	p.absSection.SetBorder(true)
	p.absSection.SetTitle(" Audiobookshelf Integration: ")
	p.absSection.SetTitleAlign(tview.AlignLeft)

	absFormLeft := newForm()
	absFormLeft.SetHorizontal(false)
	p.uploadToAudiobookshelf = absFormLeft.AddCheckbox("Upload the audiobook to Audiobookshelf server?", false, func(t bool) { p.configCopy.SetUploadToAudiobookshelf(t) })
	p.audiobookshelfUrl = absFormLeft.AddInputField("Audiobookshelf Server URL:", "", 40, nil, func(t string) { p.configCopy.SetAudiobookshelfUrl(t) })
	p.audiobookshelfUser = absFormLeft.AddInputField("Audiobookshelf Server User:", "", 40, nil, func(t string) { p.configCopy.SetAudiobookshelfUser(t) })
	p.audiobookshelfPassword = absFormLeft.AddPasswordField("Audiobookshelf Server Password:", "", 40, 0, func(t string) { p.configCopy.SetAudiobookshelfPassword(t) })
	p.audiobookshelfLibrary = absFormLeft.AddInputField("Audiobookshelf destination Library:", "", 40, nil, func(t string) { p.configCopy.SetAudiobookshelfLibrary(t) })
	p.scanAudiobookshelf = absFormLeft.AddCheckbox("Scan the Audiobookshelf library after copy/upload?", false, func(t bool) { p.configCopy.SetScanAudiobookshelf(t) })
	p.absSection.AddItem(absFormLeft.Form, 0, 0, 1, 1, 0, 0, true)

	// absFormRight := newForm()
	// absFormRight.SetHorizontal(false)
	// p.absSection.AddItem(absFormRight.f, 0, 1, 1, 1, 0, 0, true)

	p.mainGrid.AddItem(p.absSection.Grid, 2, 0, 1, 1, 0, 0, true)

	// screen navigation order
	p.mainGrid.SetNavigationOrder(
		p.searchCondition,
		p.rowsPerPage,
		p.useMockField,
		p.saveMockField,
		p.outputDir,
		p.copyToOutputDir,
		p.tmpDir,
		p.logFileNameField,
		p.logLevelField,
		p.concurrentDownloaders,
		p.concurrentEncoders,
		p.reEncodeFiles,
		p.bitRate,
		p.sampleRate,
		p.maxFileSize,
		p.shortenTitles,
		p.uploadToAudiobookshelf,
		p.audiobookshelfUrl,
		p.audiobookshelfUser,
		p.audiobookshelfPassword,
		p.audiobookshelfLibrary,
		p.scanAudiobookshelf,
		p.saveConfigButton,
		p.cancelButton,
	)

	return p
}

func (p *ConfigPage) checkMQ() {
	m := p.mq.GetMessage(mq.ConfigPage)
	if m != nil {
		p.dispatchMessage(m)
	}
}

func (p *ConfigPage) dispatchMessage(m *mq.Message) {
	switch dto := m.Dto.(type) {
	case *dto.DisplayConfigCommand:
		go p.displayConfig(dto)
	default:
		m.UnsupportedTypeError(mq.ConfigPage)
	}
}

func (p *ConfigPage) displayConfig(c *dto.DisplayConfigCommand) {
	p.configCopy = c.Config
	p.outputDir.SetText(p.configCopy.GetOutputDir())
	p.copyToOutputDir.SetChecked(p.configCopy.IsCopyToOutputDir())
	p.tmpDir.SetText(p.configCopy.GetTmpDir())

	p.logFileNameField.SetText(p.configCopy.GetLogFileName())
	p.logLevelField.SetCurrentOption(utils.GetIndex(logger.LogLeves(), p.configCopy.GetLogLevel()))
	p.searchCondition.SetText(p.configCopy.GetSearchCondition())
	p.rowsPerPage.SetText(utils.ToString(p.configCopy.GetRowsPerPage()))
	p.useMockField.SetChecked(p.configCopy.IsUseMock())
	p.saveMockField.SetChecked(p.configCopy.IsSaveMock())

	p.concurrentDownloaders.SetText(utils.ToString(p.configCopy.GetConcurrentDownloaders()))
	p.concurrentEncoders.SetText(utils.ToString(p.configCopy.GetConcurrentEncoders()))
	p.reEncodeFiles.SetChecked(p.configCopy.IsReEncodeFiles())
	p.bitRate.SetText(utils.ToString(p.configCopy.GetBitRate()))
	p.sampleRate.SetText(utils.ToString(p.configCopy.GetSampleRate()))
	p.maxFileSize.SetText(utils.ToString(p.configCopy.GetMaxFileSizeMb()))
	p.shortenTitles.SetChecked(p.configCopy.IsShortenTitle())

	p.uploadToAudiobookshelf.SetChecked(p.configCopy.IsUploadToAudiobookshef())
	p.audiobookshelfUrl.SetText(p.configCopy.GetAudiobookshelfUrl())
	p.audiobookshelfLibrary.SetText(p.configCopy.GetAudiobookshelfLibrary())
	p.scanAudiobookshelf.SetChecked(p.configCopy.IsScanAudiobookshef())
	p.audiobookshelfUser.SetText(p.configCopy.GetAudiobookshelfUser())
	p.audiobookshelfPassword.SetText(p.configCopy.GetAudiobookshelfPassword())

	ui.Draw()
	ui.SetFocus(p.configSection.Grid)
}

func (p *ConfigPage) SaveConfig() {
	p.mq.SendMessage(mq.ConfigPage, mq.ConfigController, &dto.SaveConfigCommand{Config: p.configCopy}, true)
	p.mq.SendMessage(mq.ConfigPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}

func (p *ConfigPage) Cancel() {
	p.mq.SendMessage(mq.ConfigPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}

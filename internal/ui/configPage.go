package ui

import (
	"strings"

	"github.com/rivo/tview"
	"github.com/vpoluyaktov/abb_ia/internal/config"
	"github.com/vpoluyaktov/abb_ia/internal/dto"
	"github.com/vpoluyaktov/abb_ia/internal/logger"
	"github.com/vpoluyaktov/abb_ia/internal/utils"

	"github.com/vpoluyaktov/abb_ia/internal/mq"
)

type ConfigPage struct {
	mq   *mq.Dispatcher
	grid *tview.Grid

	configCopy    config.Config
	configSection *tview.Grid
	buildSection  *tview.Grid
	absSection    *tview.Grid

	// Audobookbuilder config section
	logFileNameField *tview.InputField
	logLevelField    *tview.DropDown
	searchCondition  *tview.InputField
	maxSearchRows    *tview.InputField
	useMockField     *tview.Checkbox
	saveMockField    *tview.Checkbox
	outputDir        *tview.InputField

	// audiobook build config section
	concurrentDownloaders *tview.InputField
	concurrentEncoders    *tview.InputField
	reEncodeFiles         *tview.Checkbox
	bitRate               *tview.InputField
	sampleRate            *tview.InputField
	maxFileSize           *tview.InputField
	shortenTitles         *tview.Checkbox

	// audiobookshelf config section
	copyToAudiobookshelf   *tview.Checkbox
	audiobookshelfUrl      *tview.InputField
	audiobookshelfDir      *tview.InputField
	audiobookshelfUser     *tview.InputField
	audiobookshelfPassword *tview.InputField
	audiobookshelfLibrary  *tview.InputField

	saveConfigButton *tview.Button
	cancelButton     *tview.Button
}

func newConfigPage(dispatcher *mq.Dispatcher) *ConfigPage {
	p := &ConfigPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.ConfigPage, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(-1, -1, -1)
	p.grid.SetColumns(0)

	// Audobookbuilder config section
	p.configSection = tview.NewGrid()
	p.configSection.SetColumns(-2, -2, -1)
	p.configSection.SetBorder(true)
	p.configSection.SetTitle(" Audiobook Builder Configuration: ")
	p.configSection.SetTitleAlign(tview.AlignLeft)

	configFormLeft := newForm()
	configFormLeft.SetHorizontal(false)
	p.outputDir = configFormLeft.AddInputField("Output (working) directory:", "", 40, nil, func(t string) { p.configCopy.SetOutputDir(t) })
	p.logFileNameField = configFormLeft.AddInputField("Log file name:", "", 40, nil, func(t string) { p.configCopy.SetLogfileName(t) })
	p.logLevelField = configFormLeft.AddDropdown("Log level:", utils.AddSpaces(logger.LogLeves()), 1, func(o string, i int) { p.configCopy.SetLogLevel(strings.TrimSpace(o)) })
	p.useMockField = configFormLeft.AddCheckbox("Use mock:", false, func(t bool) { p.configCopy.SetUseMock(t) })
	p.saveMockField = configFormLeft.AddCheckbox("Save mock:", false, func(t bool) { p.configCopy.SetSaveMock(t) })
	p.configSection.AddItem(configFormLeft.f, 0, 0, 1, 1, 0, 0, true)

	configFormRight := newForm()
	configFormRight.SetHorizontal(false)
	p.searchCondition = configFormRight.AddInputField("Default Search condition", "", 40, nil, func(t string) { p.configCopy.SetSearchCondition(t) })
	p.maxSearchRows = configFormRight.AddInputField("Maximum rows in the search result:", "", 4, acceptInt, func(t string) { p.configCopy.SetSearchRowsMax(utils.ToInt(t)) })
	p.configSection.AddItem(configFormRight.f, 0, 1, 1, 1, 0, 0, true)

	buttonsForm := newForm()
	buttonsForm.SetHorizontal(false)
	buttonsForm.SetButtonsAlign(tview.AlignRight)
	p.saveConfigButton = buttonsForm.AddButton("Save Settings", p.SaveConfig)
	p.cancelButton = buttonsForm.AddButton("Cancel", p.Cancel)
	p.configSection.AddItem(buttonsForm.f, 0, 2, 1, 1, 0, 0, false)

	p.grid.AddItem(p.configSection, 0, 0, 1, 1, 0, 0, true)

	// audiobook build configuration section
	p.buildSection = tview.NewGrid()
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
	p.buildSection.AddItem(buildFormLeft.f, 0, 0, 1, 1, 0, 0, true)

	buildFormRight := newForm()
	buildFormRight.SetHorizontal(false)
	p.maxFileSize = buildFormRight.AddInputField("Audiobook part max file size (Mb):", "", 6, acceptInt, func(t string) { p.configCopy.SetMaxFileSizeMb(utils.ToInt(t)) })
	p.shortenTitles = buildFormRight.AddCheckbox("Shorten titles (for ex. Old Time Radio -> OTRR)?", false, func(t bool) { p.configCopy.SetShortenTitles(t) })
	p.buildSection.AddItem(buildFormRight.f, 0, 1, 1, 1, 0, 0, true)

	p.grid.AddItem(p.buildSection, 1, 0, 1, 1, 0, 0, true)

	// audiobookshelf config section
	p.absSection = tview.NewGrid()
	p.absSection.SetColumns(-1)
	p.absSection.SetBorder(true)
	p.absSection.SetTitle(" Audiobookshelf Integration: ")
	p.absSection.SetTitleAlign(tview.AlignLeft)

	absFormLeft := newForm()
	absFormLeft.SetHorizontal(false)
	p.copyToAudiobookshelf = absFormLeft.AddCheckbox("Copy the audiobook to Audiobookshelf directory?", false, func(t bool) { p.configCopy.SetCopyToAudiobookshelf(t) })
	p.audiobookshelfUrl = absFormLeft.AddInputField("Audiobookshelf Server URL:", "", 40, nil, func(t string) { p.configCopy.SetAudiobookshelfUrl(t) })
	p.audiobookshelfDir = absFormLeft.AddInputField("Audiobookshelf Server Directory:", "", 60, nil, func(t string) { p.configCopy.SetAudiobookshelfDir(t) })
	p.audiobookshelfLibrary = absFormLeft.AddInputField("Audiobookshelf destination Library:", "", 40, nil, func(t string) { p.configCopy.SetAudiobookshelfLibrary(t) })
	p.audiobookshelfUser = absFormLeft.AddInputField("Audiobookshelf Server User:", "", 40, nil, func(t string) { p.configCopy.SetAudiobookshelfUser(t) })
	p.audiobookshelfPassword = absFormLeft.AddPasswordField("Audiobookshelf Server Password:", "", 40, 0, func(t string) { p.configCopy.SetAudiobookshelfPassword(t) })
	p.absSection.AddItem(absFormLeft.f, 0, 0, 1, 1, 0, 0, true)

	// absFormRight := newForm()
	// absFormRight.SetHorizontal(false)
	// p.absSection.AddItem(absFormRight.f, 0, 1, 1, 1, 0, 0, true)

	p.grid.AddItem(p.absSection, 2, 0, 1, 1, 0, 0, true)

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
	p.logFileNameField.SetText(p.configCopy.GetLogFileName())
	p.logLevelField.SetCurrentOption(utils.GetIndex(logger.LogLeves(), p.configCopy.GetLogLevel()))
	p.searchCondition.SetText(p.configCopy.GetSearchCondition())
	p.maxSearchRows.SetText(utils.ToString(p.configCopy.GetSearchRowsMax()))
	p.useMockField.SetChecked(p.configCopy.IsUseMock())
	p.saveMockField.SetChecked(p.configCopy.IsSaveMock())

	p.concurrentDownloaders.SetText(utils.ToString(p.configCopy.GetConcurrentDownloaders()))
	p.concurrentEncoders.SetText(utils.ToString(p.configCopy.GetConcurrentEncoders()))
	p.reEncodeFiles.SetChecked(p.configCopy.IsReEncodeFiles())
	p.bitRate.SetText(utils.ToString(p.configCopy.GetBitRate()))
	p.sampleRate.SetText(utils.ToString(p.configCopy.GetSampleRate()))
	p.maxFileSize.SetText(utils.ToString(p.configCopy.GetMaxFileSizeMb()))
	p.shortenTitles.SetChecked(p.configCopy.IsShortenTitle())

	p.copyToAudiobookshelf.SetChecked(p.configCopy.IsCopyToAudiobookshelf())
	p.audiobookshelfUrl.SetText(p.configCopy.GetAudiobookshelfUrl())
	p.audiobookshelfDir.SetText(p.configCopy.GetAudiobookshelfDir())
	p.audiobookshelfLibrary.SetText(p.configCopy.GetAudiobookshelfLibrary())
	p.audiobookshelfUser.SetText(p.configCopy.GetAudiobookshelfUser())
	p.audiobookshelfPassword.SetText(p.configCopy.GetAudiobookshelfPassword())

	p.mq.SendMessage(mq.ConfigPage, mq.TUI, &dto.DrawCommand{Primitive: nil}, true)
	p.mq.SendMessage(mq.ConfigPage, mq.TUI, &dto.SetFocusCommand{Primitive: p.configSection}, true)
}

func (p *ConfigPage) SaveConfig() {
	p.mq.SendMessage(mq.ConfigPage, mq.ConfigController, &dto.SaveConfigCommand{Config: p.configCopy}, true)
	p.mq.SendMessage(mq.ConfigPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}

func (p *ConfigPage) Cancel() {
	p.mq.SendMessage(mq.ConfigPage, mq.Frame, &dto.SwitchToPageCommand{Name: "SearchPage"}, false)
}

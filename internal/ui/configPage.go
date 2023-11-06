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

	logFileNameField       *tview.InputField
	logLevelField          *tview.DropDown
	useMockField           *tview.Checkbox
	saveMockField          *tview.Checkbox
	audiobookshelfUrl      *tview.InputField
	audiobookshelfDir      *tview.InputField
	audiobookshelfUser     *tview.InputField
	audiobookshelfPassword *tview.InputField

	saveConfigButton *tview.Button
	cancelButton     *tview.Button
}

func newConfigPage(dispatcher *mq.Dispatcher) *ConfigPage {
	p := &ConfigPage{}
	p.mq = dispatcher
	p.mq.RegisterListener(mq.ConfigPage, p.dispatchMessage)

	p.grid = tview.NewGrid()
	p.grid.SetRows(-1)
	p.grid.SetColumns(0)

	// config section
	p.configSection = tview.NewGrid()
	p.configSection.SetColumns(-2, -1)
	p.configSection.SetBorder(true)
	p.configSection.SetTitle(" Audiobook Builder Default Configuration")
	p.configSection.SetTitleAlign(tview.AlignLeft)
	f := newForm()
	f.SetHorizontal(false)
	p.logFileNameField = f.AddInputField("Log file name:", "", 40, nil, func(t string) { p.configCopy.SetLogfileName(t) })
	p.logLevelField = f.AddDropdown("Log level:", utils.AddSpaces(logger.LogLeves()), 1, func(o string, i int) { p.configCopy.SetLogLevel(strings.TrimSpace(o)) })
	p.useMockField = f.AddCheckbox("Use mock:", false, func(t bool) { p.configCopy.SetUseMock(t) })
	p.saveMockField = f.AddCheckbox("Save mock:", false, func(t bool) { p.configCopy.SetSaveMock(t) })

	p.audiobookshelfUrl = f.AddInputField("Audiobookshelf Server URL:", "", 40, nil, func(t string) { p.configCopy.SetAudiobookshelfUrl(t) })
	p.audiobookshelfDir = f.AddInputField("Audiobookshelf Server Directory:", "", 60, nil, func(t string) { p.configCopy.SetAudiobookshelfDir(t) })
	p.audiobookshelfUser = f.AddInputField("Audiobookshelf Server User:", "", 40, nil, func(t string) { p.configCopy.SetAudiobookshelfUser(t) })
	p.audiobookshelfPassword = f.AddInputField("Audiobookshelf Server Password:", "", 40, nil, func(t string) { p.configCopy.SetAudiobookshelfPassword(t) })

	p.configSection.AddItem(f.f, 0, 0, 1, 1, 0, 0, true)
	f = newForm()
	f.SetHorizontal(false)
	f.f.SetButtonsAlign(tview.AlignRight)
	p.saveConfigButton = f.AddButton("Save Settings", p.SaveConfig)
	p.cancelButton = f.AddButton("Cancel", p.Cancel)
	p.configSection.AddItem(f.f, 0, 1, 1, 1, 0, 0, false)
	p.grid.AddItem(p.configSection, 0, 0, 1, 1, 0, 0, true)
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
	p.logFileNameField.SetText(p.configCopy.GetLogFileName())
	p.logLevelField.SetCurrentOption(utils.GetIndex(logger.LogLeves(), p.configCopy.GetLogLevel()))
	p.useMockField.SetChecked(p.configCopy.IsUseMock())
	p.saveMockField.SetChecked(p.configCopy.IsSaveMock())
	p.audiobookshelfUrl.SetText(p.configCopy.GetAudiobookshelfUrl())
	p.audiobookshelfDir.SetText(p.configCopy.GetAudiobookshelfDir())
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

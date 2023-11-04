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

	logFileNameField *tview.InputField
	logLevelField    *tview.DropDown
	useMockField     *tview.Checkbox
	saveMockField    *tview.Checkbox

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

	// config section
	p.configSection = tview.NewGrid()
	p.configSection.SetColumns(-2, -1)
	p.configSection.SetBorder(true)
	p.configSection.SetTitle(" Audiobook Builder Default Configuration")
	p.configSection.SetTitleAlign(tview.AlignLeft)
	f := newForm()
	f.SetHorizontal(false)
	p.logFileNameField = f.AddInputField("Log file name:", "", 40, nil, func(t string) { p.configCopy.LogFileName = t })
	p.logLevelField = f.AddDropdown("Log level:", utils.AddSpaces(logger.LogLeves()), 1, func(o string, i int) { p.configCopy.LogLevel = strings.TrimSpace(o) })
	p.useMockField = f.AddCheckbox("Use mock:", false, func(t bool) { p.configCopy.UseMock = t })
	p.saveMockField = f.AddCheckbox("Save mock:", false, func(t bool) { p.configCopy.SaveMock = t })

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
	p.logFileNameField.SetText(p.configCopy.LogFileName)
	p.logLevelField.SetCurrentOption(utils.GetIndex(logger.LogLeves(), p.configCopy.LogLevel))
	p.useMockField.SetChecked(p.configCopy.UseMock)
	p.saveMockField.SetChecked(p.configCopy.SaveMock)

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

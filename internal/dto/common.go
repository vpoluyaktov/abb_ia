package dto

import "code.rocketnine.space/tslocum/cview"

type Dto interface {
}

const GeneralCommandType = "dto.GeneralCommand"
type GeneralCommand struct {
	Command string
}

const SetFocusCommandType = "dto.SetFocusCommand"
type SetFocusCommand struct {
	Primitive cview.Primitive
}

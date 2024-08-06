package logTypeEnum

import "errors"

type stageType string

const (
	Start stageType = "START"
	End stageType = "END" 
)

func (t stageType) Validate() error {
	switch t {
		case Start, End:
			return nil
		default:
			return errors.New("invalid stage type")
	}
} 
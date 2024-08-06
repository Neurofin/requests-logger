package logTypeEnum

import "errors"

type StageType string

const (
	Start StageType = "START"
	End StageType   = "END" 
)

func (t StageType) Validate() error {
	switch t {
		case Start, End:
			return nil
		default:
			return errors.New("invalid stage type")
	}
} 
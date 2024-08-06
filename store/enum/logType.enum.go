package logTypeEnum

import "errors"

type LogType string

const (
	API LogType = "API"
	Error LogType = "ERROR"
	Debug LogType = "DEBUG"
	Info LogType = "INFO"
)

func (t LogType) Validate() error {
	switch t {
		case API, Error, Debug, Info:
			return nil
		default:
			return errors.New("invalid log type")
	}
} 

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
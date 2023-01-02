package extension

import "io"

var (
	Exported_generateRegisterOutput  = generateRegisterOutput
	Exported_generateEventNextOutput = generateEventNextOutput
)

type Exported_event = events

func (es *Exported_event) Exported_toRequestBody() (io.Reader, error) {
	return es.toRequestBody()
}

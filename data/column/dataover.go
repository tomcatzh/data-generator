package column

// DataOver is a acceptable error cause by column data run over
type DataOver struct {
	reason string
}

func (e *DataOver) Error() string {
	return e.reason
}

// NewDataOver create a new data over error
func NewDataOver(reason string) *DataOver {
	return &DataOver{reason: reason}
}

// IsDataOver check if is data over error
func IsDataOver(err error) bool {
	_, ok := err.(*DataOver)
	return ok
}

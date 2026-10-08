//go:generate goenums -f -c ./encoding_type.go
package encodingtype

type encodingType int8

const (
	unknown encodingType = iota
	json
	text
)

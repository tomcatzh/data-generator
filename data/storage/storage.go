package storage

import (
	"errors"
	"fmt"
	"io"
)

// Storage is a interface of file store, such as local store or S3 store etc
type Storage interface {
	Save(key string, reader io.Reader) (int64, error)
}

// NewStorage returns a Storeage interface from template
func NewStorage(s map[string]interface{}) (Storage, error) {
	stype, ok := s["Type"].(string)
	if !ok || stype == "" {
		return nil, errors.New("Template do not have Storage type")
	}

	switch stype {
	case "Local":
		return newStorageLocal(s), nil
	case "S3":
		return newStorageS3(s)
	default:
		return nil, fmt.Errorf("Unexecepted Storage type: %v", stype)
	}
}

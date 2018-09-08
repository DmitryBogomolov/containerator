package containerator

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type mappingListVar struct {
	list     []Mapping
	sep      string
	allowOne bool
}

// MappingListVar keeps list of Mapping.
type MappingListVar interface {
	flag.Value
	Get() []Mapping
}

func (m *mappingListVar) String() string {
	return fmt.Sprintf("%v", m.list)
}

func (m *mappingListVar) Set(value string) error {
	parts := strings.SplitN(value, m.sep, 2)
	mapping := Mapping{Source: parts[0]}
	if len(parts) > 1 {
		mapping.Target = parts[1]
	} else if !m.allowOne {
		return errors.New("not a pair")
	}
	m.list = append(m.list, mapping)
	return nil
}

func (m *mappingListVar) Get() []Mapping {
	return m.list
}

// NewMappingListVar constructs MappingListVar.
func NewMappingListVar(sep string, allowOne bool) MappingListVar {
	return &mappingListVar{sep: sep, allowOne: allowOne}
}

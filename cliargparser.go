package containerator

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type _MappingListVar struct {
	list     []Mapping
	sep      string
	allowOne bool
}

/*
MappingListVar interface is used to parse command line flags into list of Mapping instances.

Can be used with `flag` package to parse volume, port, environment variables mappings and
then pass them to RunContainerOptions instance.

	-v /src1:/dst1 -v /src2:/dst2
	-p 50001:3001 -p 50002:3002
	-e A=1 -e B
*/
type MappingListVar interface {
	flag.Value
	Get() []Mapping
}

func (m *_MappingListVar) String() string {
	return fmt.Sprintf("%v", m.list)
}

func (m *_MappingListVar) Set(value string) error {
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

func (m *_MappingListVar) Get() []Mapping {
	return m.list
}

/*
NewMappingListVar creates instance that implements MappingListVar interface.

`sep` is separator string, `allowOne` allows providing value without separator.

	volumes := NewMappingListVar(":", false)
	ports := NewMappingListVar(":", false)
	env := NewMappingListVar("=", true)
	flag.Var(volumes, "v", "")
	flag.Var(ports, "p", "")
	flag.Var(env, "e", "")

	options.Volumes = volumes.Get()
	options.Ports = ports.Get()
	options.Env = env.Get()
*/
func NewMappingListVar(sep string, allowOne bool) MappingListVar {
	return &_MappingListVar{sep: sep, allowOne: allowOne}
}

package core

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type _MappingListFlag struct {
	mappings  []Mapping
	separator string
	allowOne  bool
}

/*
MappingListFlag interface is used to parse command line flags into list of Mapping instances.

Can be used with `flag` package to parse volume, port, environment variables mappings and
then pass them to RunContainerOptions instance.

	-v /src1:/dst1 -v /src2:/dst2
	-p 50001:3001 -p 50002:3002
	-e A=1 -e B
*/
type MappingListFlag interface {
	flag.Value
	Get() []Mapping
}

func (mappingListFlag *_MappingListFlag) String() string {
	return fmt.Sprintf("%v", mappingListFlag.mappings)
}

func (mappingListFlag *_MappingListFlag) Set(value string) error {
	parts := strings.SplitN(value, mappingListFlag.separator, 2)
	source := parts[0]
	target := ""
	if len(parts) > 1 {
		target = parts[1]
	} else if !mappingListFlag.allowOne {
		return errors.New("not a pair")
	}
	mappingListFlag.mappings = append(mappingListFlag.mappings, Mapping{source, target})
	return nil
}

func (mappingListFlag *_MappingListFlag) Get() []Mapping {
	return mappingListFlag.mappings
}

/*
NewMappingListFlag creates MappingListFlag instance.

`separator` is separator string, `allowOne` allows providing value without separator.

	volumes := NewMappingListFlag(":", false)
	ports := NewMappingListFlag(":", false)
	env := NewMappingListFlag("=", true)
	flag.Var(volumes, "v", "")
	flag.Var(ports,NewMappingListFlag "p", "")
	flag.Var(env, "e", "")

	options.Volumes = volumes.Get()
	options.Ports = ports.Get()
	options.Env = env.Get()
*/
func NewMappingListFlag(separator string, allowOne bool) MappingListFlag {
	return &_MappingListFlag{separator: separator, allowOne: allowOne}
}

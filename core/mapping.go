package core

import "encoding/json"

// Mapping stores key-value pair. Used for volumes, ports, environment variables.
type Mapping struct {
	Source string
	Target string
}

func (mapping Mapping) toMap() map[string]string {
	ret := map[string]string{}
	ret[mapping.Source] = mapping.Target
	return ret
}

func (mapping *Mapping) fromMap(data map[string]string) {
	for key, val := range data {
		mapping.Source = key
		mapping.Target = val
	}
}

// MarshalJSON implements `json.Marshaler` interface.
func (mapping Mapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(mapping.toMap())
}

// UnmarshalJSON implements `json.Unmarshaler` interface.
func (mapping *Mapping) UnmarshalJSON(data []byte) error {
	var tmp map[string]string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	mapping.fromMap(tmp)
	return nil
}

// MarshalYAML implements `yaml.Marshaler` interface.
func (mapping Mapping) MarshalYAML() (interface{}, error) {
	return mapping.toMap(), nil
}

// UnmarshalYAML implements `yaml.Unmarshaler` interface.
func (mapping *Mapping) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp map[string]string
	if err := unmarshal(&tmp); err != nil {
		return err
	}
	mapping.fromMap(tmp)
	return nil
}

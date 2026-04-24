package schema

import (
	_ "embed"
	"encoding/json"
	"strings"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed config.schema.json
var configSchemaJSON string

type ValidationError struct {
	Location string
	Message  string
}

func (e ValidationError) String() string {
	if e.Location == "" || e.Location == "/" {
		return e.Message
	}
	return e.Location + ": " + e.Message
}

func ValidateConfig(data map[string]any) ([]ValidationError, error) {
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("config.schema.json", strings.NewReader(configSchemaJSON)); err != nil {
		return nil, err
	}
	sch, err := compiler.Compile("config.schema.json")
	if err != nil {
		return nil, err
	}

	// Round-trip through JSON to normalise YAML types to JSON-compatible types.
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var normalised any
	if err := json.Unmarshal(b, &normalised); err != nil {
		return nil, err
	}

	if err := sch.Validate(normalised); err != nil {
		ve, ok := err.(*jsonschema.ValidationError)
		if !ok {
			return nil, err
		}
		basic := ve.BasicOutput()
		var errs []ValidationError
		for _, e := range basic.Errors {
			if e.Error == "" || strings.HasPrefix(e.Error, "doesn't validate with") {
				continue
			}
			errs = append(errs, ValidationError{
				Location: e.InstanceLocation,
				Message:  e.Error,
			})
		}
		return errs, nil
	}
	return nil, nil
}

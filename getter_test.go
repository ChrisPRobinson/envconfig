package envconfig

import (
	"errors"
	"reflect"
	"testing"
)

type PluginTestErrSpecification struct {
	PluginTestGetter

	Token string `error:"ohno"`
}

type PluginTestSpecification struct {
	PluginTestGetter

	Token  string `special:"foo"`
	Token2 string `default:"d" special:"foo"`
}

type PluginTestGetter struct {
	EnvVarGetter

	SpecialValues map[string]map[string]string
}

func (t PluginTestGetter) Provider() string {
	return "PluginTestGetter"
}

func (t *PluginTestGetter) Get(key, alt string, tags reflect.StructTag) (string, bool, error) {
	retError := tags.Get("error")
	if retError != "" {
		return "", false, errors.New(retError)
	}

	special := tags.Get("special")
	if special != "" {
		var valueMap map[string]string
		var value string
		var exists bool
		valueMap, exists = t.SpecialValues[key]
		if exists {
			value, exists = valueMap[special]
			if exists {
				return value, true, nil
			}
		}
	}

	return t.EnvVarGetter.Get(key, alt, tags)
}

func TestGetter(t *testing.T) {
	val := map[string]string{"foo": "bar"}
	mapVal := map[string]map[string]string{"TOKEN": val}
	v := PluginTestSpecification{PluginTestGetter: PluginTestGetter{SpecialValues: mapVal}}
	err := Process("", &v)

	if err != nil {
		t.Errorf("Processing has failed: %s", err.Error())
	}

	if v.Token != "bar" {
		t.Errorf("Expected value of Token to be 'bar' but got '%s'", v.Token)
	}

	if v.Token2 != "d" {
		t.Errorf("Expected value of Token2 to be 'd' but got '%s'", v.Token2)
	}
}

func TestGetterError(t *testing.T) {
	mapVal := map[string]map[string]string{}
	v := PluginTestErrSpecification{PluginTestGetter: PluginTestGetter{SpecialValues: mapVal}}

	err := Process("", &v)

	if err == nil {
		t.Errorf("Expected an error on processing")
	} else {
		if err.Error() != "envconfig.Process: getting value TOKEN from provider PluginTestGetter for field Token has failed: details: ohno" {
			t.Errorf("Expected a different error message than '%s'", err.Error())
		}
	}

}

package mapstructure_hooks

import (
	"net/url"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func decode(i interface{}, data map[string]interface{}) {
	dc := &mapstructure.DecoderConfig{
		Metadata:   nil,
		Result:     i,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(All()...),
	}

	md, err := mapstructure.NewDecoder(dc)

	if err != nil {
		panic(err)
	}

	if err := md.Decode(data); err != nil {
		panic(err)
	}
}

type testStruct struct {
	IsOK   bool
	Name   string
	Count  int
	Ratio  float64
	Target *url.URL
}

func TestStringToBoolHook(t *testing.T) {
	var s testStruct

	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"IsOK": true}) })
	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"IsOK": "true"}) })
	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"IsOK": "false"}) })
	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"IsOK": "True"}) })
	assert.Panics(t, func() { decode(&s, map[string]interface{}{"IsOK": "notaboolean"}) })
}

func TestStringToNumeric(t *testing.T) {
	var s testStruct

	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"Count": 123}) })
	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"Count": "123"}) })
	assert.Panics(t, func() { decode(&s, map[string]interface{}{"Count": true}) })
	assert.Panics(t, func() { decode(&s, map[string]interface{}{"Count": "notanint"}) })

	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"Ratio": 3.14}) })
	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"Ratio": "3.14"}) })
	assert.Panics(t, func() { decode(&s, map[string]interface{}{"Ratio": true}) })
	assert.Panics(t, func() { decode(&s, map[string]interface{}{"Ratio": "notafloat"}) })
}

func TestStringToURL(t *testing.T) {
	var s testStruct

	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"Target": "http://google.fr/coco?123=456"}) })
	assert.Equal(t, "google.fr", s.Target.Hostname())
	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"Target": "file:///tmp/example.json"}) })
	assert.NotPanics(t, func() { decode(&s, map[string]interface{}{"Target": "s3://mu-bucket/deep/nested/file.json"}) })
	assert.Equal(t, "s3", s.Target.Scheme)

	assert.Panics(t, func() { decode(&s, map[string]interface{}{"Target": 1000}) })
}

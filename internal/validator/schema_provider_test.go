package validator

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSchemaProvider(t *testing.T) {
	p := NewSchemaProvider()
	assert.NotNil(t, p)
}

func TestSchemaProvider_loadSchema_NoFile(t *testing.T) {
	s, err := loadSchema("_")
	assert.Nil(t, s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestSchemaProvider_loadSchema_InvalidSchema(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "invalid-schema-")
	assert.NoError(t, err)
	defer os.Remove(f.Name())
	f.Write([]byte("not json"))
	f.Close()

	s, err := loadSchema(f.Name())
	assert.Nil(t, s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestSchemaProvider_prepareUrl(t *testing.T) {
	cases := []struct {
		url    string
		expUrl string
	}{
		{"/", ""},
		{"//", ""},
		{"///", ""},
		{"sOmeUrl", "someurl"},
		{"/sOme/Url/", "some/url"},
	}
	for _, c := range cases {
		actUrl := prepareUrl(c.url)
		assert.Equal(t, c.expUrl, actUrl)
	}
}

func TestSchemaProvider_Register(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "schema-register-")
	assert.NoError(t, err)
	defer os.Remove(f.Name())
	f.Write([]byte("{}"))
	f.Close()

	p := NewSchemaProvider()
	err = p.Register("x", f.Name())
	assert.NoError(t, err)
}

func TestSchemaProvider_Get(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "schema-get-")
	assert.NoError(t, err)
	defer os.Remove(f.Name())
	f.Write([]byte("{}"))
	f.Close()

	p := NewSchemaProvider()
	err = p.Register("url", f.Name())
	assert.NoError(t, err)

	cases := []struct {
		url      string
		expFound bool
	}{
		{"url", true},
		{"URL", true},
		{"/url/", true},
		{"/url", true},
		{"/url/one", false},
		{"/another/url", false},
	}
	for _, c := range cases {
		s, err := p.Get(c.url)
		if c.expFound {
			assert.NoError(t, err)
			assert.NotNil(t, s)
		} else {
			assert.Error(t, err)
			assert.Nil(t, s)
		}
	}
}

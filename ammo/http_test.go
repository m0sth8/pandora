package ammo

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/yandex/pandora/config"
	"github.com/yandex/pandora/utils"
)

const (
	httpTestFilename = "./testdata/ammo.jsonline"
)

func TestNewHttpProvider(t *testing.T) {
	c := &config.AmmoProvider{
		AmmoSource: httpTestFilename,
		AmmoLimit:  10,
	}
	provider, err := NewHttpProvider(c)
	require.NoError(t, err)

	httpProvider, casted := provider.(*HttpProvider)
	require.True(t, casted, "NewHttpProvider should return *HttpProvider type")

	// look at defaults
	assert.Equal(t, 10, httpProvider.ammoLimit)
	assert.Equal(t, 0, httpProvider.passes)
	assert.NotNil(t, httpProvider.sink)
	assert.NotNil(t, httpProvider.BaseProvider.source)
	assert.NotNil(t, httpProvider.BaseProvider.decoder)

	// compare data
	actualData, err := ioutil.ReadAll(httpProvider.ammoFile)
	require.NoError(t, err)

	expectedData, err := ioutil.ReadFile(httpTestFilename)
	require.NoError(t, err)

	assert.Equal(t, expectedData, actualData)

	// test wrong file

	c.AmmoSource = "./testdata/badammo_name.jsonline"
	_, err = NewHttpProvider(c)
	assert.Error(t, err)
}

func TestHttpProvider(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	providerCtx, _ := context.WithCancel(ctx)

	data, err := ioutil.ReadFile(httpTestFilename)
	require.NoError(t, err)

	ammoCh := make(chan Ammo, 128)
	provider := &HttpProvider{
		passes:   2,
		ammoFile: bytes.NewReader(data),
		sink:     ammoCh,
		BaseProvider: NewBaseProvider(
			ammoCh,
			HttpJSONDecode,
		),
	}
	promise := utils.PromiseCtx(providerCtx, provider.Start)

	ammos := Drain(ctx, provider)
	require.Len(t, ammos, 25*2) // two passes

	httpAmmo, casted := (ammos[2]).(*Http)
	require.True(t, casted, "Ammo should have *Http type")

	assert.Equal(t, "example.org", httpAmmo.Host)
	assert.Equal(t, "/02", httpAmmo.Uri)
	assert.Equal(t, "hello", httpAmmo.Tag)
	assert.Equal(t, "GET", httpAmmo.Method)
	assert.Len(t, httpAmmo.Headers, 4)

	// TODO: add test for decoding error

	select {
	case err := <-promise:
		require.NoError(t, err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

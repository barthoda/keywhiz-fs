// Copyright 2015 Square Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package keywhizfs_test

import (
	"testing"

	"github.com/square/keywhizfs"
	"github.com/stretchr/testify/assert"
)

func TestSecretMapOperations(t *testing.T) {
	assert := assert.New(t)

	s, err := keywhizfs.ParseSecret(fixture("secret.json"))
	assert.NoError(err)

	secretMap := keywhizfs.NewSecretMap()
	assert.Equal(0, secretMap.Len())
	assert.Empty(secretMap.Values())

	lookup, ok := secretMap.Get("foo")
	assert.False(ok)

	secretMap.Put("foo", *s)
	assert.Equal(1, secretMap.Len())

	values := secretMap.Values()
	assert.Len(values, 1)
	assert.Equal(*s, values[0].Secret)

	lookup, ok = secretMap.Get("foo")
	assert.True(ok)
	assert.Equal(*s, lookup.Secret)

	put := secretMap.PutIfAbsent("foo", keywhizfs.Secret{})
	assert.False(put)

	lookup, ok = secretMap.Get("foo")
	assert.True(ok)
	assert.Equal(*s, lookup.Secret)

	secretMap.Put("foo", keywhizfs.Secret{})

	lookup, ok = secretMap.Get("foo")
	assert.True(ok)
	assert.NotEqual(*s, lookup.Secret)
}

func TestSecretMapOverwrite(t *testing.T) {
	assert := assert.New(t)

	s, err := keywhizfs.ParseSecret(fixture("secret.json"))
	assert.NoError(err)

	secretMap := keywhizfs.NewSecretMap()
	secretMap.Put("foo", keywhizfs.Secret{})

	newMap := keywhizfs.NewSecretMap()
	newMap.Put("bar", *s)
	secretMap.Overwrite(newMap)

	assert.Equal(1, secretMap.Len())
	lookup, ok := secretMap.Get("bar")
	assert.True(ok)
	assert.Equal(*s, lookup.Secret)
}

func TestSecretMapTimestamp(t *testing.T) {
	assert := assert.New(t)

	secretMap := keywhizfs.NewSecretMap()
	secretMap.Put("foo", keywhizfs.Secret{})

	val, ok := secretMap.Get("foo")
	assert.True(ok)
	earlierTime := val.Time

	secretMap.Put("foo", keywhizfs.Secret{})
	val, ok = secretMap.Get("foo")
	assert.True(ok)
	assert.True(val.Time.After(earlierTime))
}

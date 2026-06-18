package mgquery_test

import (
	"testing"

	"github.com/d3v-friends/go-tools/fnEnv"
	"github.com/stretchr/testify/assert"
)

const envPath = "../.env"

func loadEnv(t *testing.T) {
	assert.NoError(t, fnEnv.Load(envPath))
}

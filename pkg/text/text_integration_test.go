package text_test

import (
	"os"
	"testing"

	"github.com/jlevesy/goats/pkg/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsing(t *testing.T) {
	file, err := os.Open("./assets/test.goats")
	require.NoError(t, err)
	defer file.Close()

	wantTestNames := []string{
		"test parsing integration",
		"test parsing integration 2",
	}

	wantInstructions := [][][]string{
		{
			{"ls", "/foo/bar", "is", "a", "long", "instruction"},
			{"assert_ok"},
		},
		{
			{"ls", "/bar/biz"},
			{"assert_ok"},
		},
	}

	parser := text.NewParser(text.NewLexer(file), mockResolver(spewerResolver))
	suite, err := parser.Parse()
	assert.NoError(t, err)

	gotInstructions := instructionsFromSuite(suite)
	gotTestNames := testNamesFromSuite(suite)

	assert.Equal(t, wantInstructions, gotInstructions)
	assert.Equal(t, wantTestNames, gotTestNames)
}

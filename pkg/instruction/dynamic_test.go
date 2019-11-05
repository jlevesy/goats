package instruction_test

import (
	"testing"

	"github.com/jlevesy/goats/pkg/instruction"
	"github.com/stretchr/testify/assert"
)

func TestLoadDynamic(t *testing.T) {
	builders := make(instruction.Builders)
	err := instruction.LoadDynamic([]string{"./assets"}, builders)
	if !assert.NoError(t, err) {
		return
	}

	bOK, has := builders["assert_ok"]
	if !assert.True(t, has) {
		return
	}

	_, err = bOK([]string{})
	if !assert.NoError(t, err) {
		return
	}

	bOK2, has := builders["assert_ok_2"]
	if !assert.True(t, has) {
		return
	}

	_, err = bOK2([]string{})
	if !assert.NoError(t, err) {
		return
	}

	fail, has := builders["fail"]
	if !assert.True(t, has) {
		return
	}

	_, err = fail([]string{})
	assert.Error(t, err)
}

func TestDiscoverDynamicInstructions(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []instruction.Tag
		wantErr bool
	}{
		{
			name: "discovers all instruction tags",
			content: `
package lol

// @instruction{name=assert_ok,builder=NewOK}
// @instruction{name=assert_failure,builder=NewFailure}
// @instruction{name=assertSomething,builder=NewSomething}
			`,
			want: []instruction.Tag{
				{
					Name:    "assert_ok",
					Builder: "NewOK",
					Package: "lol",
				},
				{
					Name:    "assert_failure",
					Builder: "NewFailure",
					Package: "lol",
				},
				{
					Name:    "assertSomething",
					Builder: "NewSomething",
					Package: "lol",
				},
			},
		},
		{
			name: "raises an error when package name is not found",
			content: `
// @instruction{name=assert_ok,builder=NewOK}
// @instruction{name=assert_failure,builder=NewFailure}
// @instruction{name=assertSomething,builder=NewSomething}
			`,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tags, err := instruction.ParseTags([]byte(test.content))

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.want, tags)
		})
	}
}

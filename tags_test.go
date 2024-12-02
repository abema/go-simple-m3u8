package m3u8

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttributes(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		attr := Attributes{
			"NUMBER":   "12345",
			"STR":      `"abcde"`,
			"NO-VALUE": "",
		}
		assert.Equal(t, `NO-VALUE,NUMBER=12345,STR="abcde"`, attr.String())
	})
}

func TestParseTagAttributes(t *testing.T) {
	testCases := []struct {
		input    string
		expected Attributes
	}{
		{
			input:    "",
			expected: Attributes{},
		},
		{
			input: "HEX1=0x12ab",
			expected: Attributes{
				"HEX1": "0x12ab",
			},
		},
		{
			input: `HEX1=0x12ab,STR="foo",NO-VALUE,HEX2=0x34cd`,
			expected: Attributes{
				"HEX1":     "0x12ab",
				"HEX2":     "0x34cd",
				"STR":      `"foo"`,
				"NO-VALUE": "",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			m, err := ParseTagAttributes(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, m)
		})
	}
}

func TestTags(t *testing.T) {
	tags := Tags{
		"EXT-X-FOO": []string{"foo"},
	}
	tag := tags.Last("EXT-X-FOO")
	assert.Equal(t, &Tag{Name: "EXT-X-FOO", Attributes: "foo"}, tag)
	tags.Set(&Tag{Name: "EXT-X-BAR", Attributes: "bar"})
	assert.Equal(t, Tags{
		"EXT-X-FOO": []string{"foo"},
		"EXT-X-BAR": []string{"bar"},
	}, tags)
	tags.RemoveByName("EXT-X-FOO")
	assert.Equal(t, Tags{"EXT-X-BAR": []string{"bar"}}, tags)
	tags.Add(&Tag{Name: "EXT-X-BAR", Attributes: "bar2"})
	assert.Equal(t, Tags{"EXT-X-BAR": []string{"bar", "bar2"}}, tags)
	assert.Equal(t, &Tag{Name: "EXT-X-BAR", Attributes: "bar"}, tags.First("EXT-X-BAR"))
	assert.Equal(t, &Tag{Name: "EXT-X-BAR", Attributes: "bar2"}, tags.Last("EXT-X-BAR"))
	tags.Set(&Tag{Name: "EXT-X-BAR", Attributes: "bar3"})
	assert.Equal(t, Tags{"EXT-X-BAR": []string{"bar3"}}, tags)
}

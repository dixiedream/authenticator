package utils_test

import (
	"testing"

	"github.com/dixiedream/authenticator/utils"
)

func TestIsAdmin(t *testing.T) {
	cases := []struct {
		value    string
		expected bool
	}{{"123", false}, {"", false}, {"a", false}, {"332", true}}

	for _, c := range cases {
		result := utils.IsAdmin(c.value)
		if result != c.expected {
			t.Logf("Result should be %t, got %t instead", c.expected, result)
			t.Fail()
		}
	}
}

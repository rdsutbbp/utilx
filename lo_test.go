package utilx

import (
	"github.com/samber/lo"
	"testing"
)

func TestLoUniq(t *testing.T) {
	names := lo.Uniq[string]([]string{"Samuel", "Marc", "Samuel"})
	t.Log(names)

	names2 := Uniq[string]([]string{"1", "2", "1"})
	t.Log(names2)
}

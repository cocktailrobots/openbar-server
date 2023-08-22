package functional

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var ts []int
		us := Map(func(t int) int { return t * 2 }, ts)
		if len(us) != 0 {
			t.Errorf("expected empty slice, got %v", us)
		}
	})

	t.Run("non-empty slice", func(t *testing.T) {
		ts := []int{1, 2, 3}
		us := Map(func(t int) int { return t * 2 }, ts)
		require.Equal(t, len(us), len(ts))
		for i := range ts {
			require.Equal(t, us[i], ts[i]*2)
		}
	})

	t.Run("convert slice", func(t *testing.T) {
		ts := []int64{1, 2, 3}
		us := Map(func(t int64) string { return strconv.FormatInt(t, 10) }, ts)
		require.Equal(t, len(us), len(ts))
		for i := range ts {
			n, err := strconv.ParseInt(us[i], 10, 64)
			require.NoError(t, err)
			require.Equal(t, ts[i], n)
		}
	})
}

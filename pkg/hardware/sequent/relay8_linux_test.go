package sequent

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRelay8States(t *testing.T) {
	r8s := Relay8States{}
	require.Equal(t, r8s, fromByte(0))
	require.Equal(t, r8s.toByte(), 0)

	// test that set does not mutate the original
	_ = r8s.Set(0, relayOn)
	require.Equal(t, r8s.toByte(), 0)

	r8s = r8s.Set(0, relayOn)
	require.Equal(t, r8s, fromByte(relayMaskRemap[0]))
	require.Equal(t, r8s.toByte(), relayMaskRemap[0])

	r8s = r8s.Set(4, relayOn)
	require.Equal(t, r8s, fromByte(relayMaskRemap[0]|relayMaskRemap[4]))
	require.Equal(t, r8s.toByte(), relayMaskRemap[0]|relayMaskRemap[4])

	r8s = r8s.Set(0, relayOff)
	require.Equal(t, r8s, fromByte(relayMaskRemap[4]))
	require.Equal(t, r8s.toByte(), relayMaskRemap[4])
}

package wire

import (
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFluids(t *testing.T) {
	fluids := []db.Fluid{
		{
			Idx:   0,
			Fluid: ptr("gin"),
		},
		{
			Idx:   3,
			Fluid: ptr("tonic"),
		},
		{
			Idx:   1,
			Fluid: ptr("sweet_vermouth"),
		},
		{
			Idx:   2,
			Fluid: ptr("campari"),
		},
	}

	fluidsWire := FromDbFluids(fluids)
	data, err := json.Marshal(fluidsWire)
	require.NoError(t, err)
	require.NoError(t, fluidsWire.Validate())

	var decoded Fluids
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	fluids2 := decoded.ToDbFluids()
	require.True(t, fluidsEqual(fluids, fluids2))
}

func fluidsEqual(fluids, fluids2 []db.Fluid) bool {
	if len(fluids) != len(fluids2) {
		return false
	}

	idxToFluid := map[int]string{}
	for _, fluid := range fluids {
		idxToFluid[fluid.Idx] = *fluid.Fluid
	}

	for _, fluid := range fluids2 {
		if idxToFluid[fluid.Idx] != *fluid.Fluid {
			return false
		}
	}

	return true
}

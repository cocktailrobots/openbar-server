package wire

type FluidVolume struct {
	Fluid    string `json:"fluid"`
	VolumeMl uint   `json:"volume_ml"`
}

type MakeRequest struct {
	FluidVolumes []FluidVolume `json:"fluid_volumes"`
}

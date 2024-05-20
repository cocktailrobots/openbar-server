package wire

type ButtonState struct {
	DepressedButtons []int `json:"depressed_buttons"`
	DurationMs       int   `json:"duration_ms"`
	Async            bool  `json:"async"`
	Forward          bool  `json:"forward"`
}

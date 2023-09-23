package buttons

type Buttons interface {
	NumButtons() int
	Update() error
	IsPressed(idx int) bool
	Close() error
}

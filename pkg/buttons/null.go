package buttons

type NullButtons struct{}

func NewNullButtons() NullButtons {
	return NullButtons{}
}

func (n NullButtons) NumButtons() int {
	return 0
}

func (n NullButtons) IsPressed(idx int) bool {
	return false
}

func (n NullButtons) Update() error {
	return nil
}

func (n NullButtons) Close() error {
	return nil
}

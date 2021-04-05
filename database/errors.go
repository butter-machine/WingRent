package database


type EntityAlreadyExists struct{}

func (m *EntityAlreadyExists) Error() string {
	return "Entity already exists"
}

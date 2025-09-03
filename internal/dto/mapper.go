package dto

type Mapper interface {
	ToState(data []any) (*State, error)
}

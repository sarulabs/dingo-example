package garage

type ErrValidation struct {
	error
	PublicMessage string
}

type ErrNotFound struct {
	error
	PublicMessage string
}

type CarManager struct {
	Repo *CarRepository
}

func (m *CarManager) Get(id string) (*Car, error) {
	return nil, nil
}

func (m *CarManager) Create(input *Car) (*Car, error) {
	return nil, nil
}

func (m *CarManager) Update(input *Car) (*Car, error) {
	return nil, nil
}

func (m *CarManager) Delele(id string) error {
	return nil
}

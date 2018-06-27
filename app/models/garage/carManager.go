package garage

type CarManager struct {
	Repo *CarRepository
}

func (m *CarManager) Create(data *Car) (*Car, error) {
	return nil, nil
}

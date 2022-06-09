package usecase

type StoreUseCase struct {
	repo DbRepo
}

func (s StoreUseCase) Find(result interface{}, filter interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s StoreUseCase) FindJson(result interface{}, filter []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s StoreUseCase) AddOrUpdate(key interface{}, data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s StoreUseCase) Get(key interface{}, result interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (s StoreUseCase) Delete(key interface{}, dataType interface{}) error {
	//TODO implement me
	panic("implement me")
}

func New(r DbRepo) *StoreUseCase {
	return &StoreUseCase{repo: r}
}

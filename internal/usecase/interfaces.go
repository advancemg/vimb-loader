package usecase

type DbRepo interface {
	Find(result interface{}, filter interface{}) error
	FindJson(result interface{}, filter []byte) error
	AddOrUpdate(key interface{}, data interface{}) error
	Get(key interface{}, result interface{}) error
	Delete(key interface{}, dataType interface{}) error
}

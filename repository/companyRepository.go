package repository

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	dbError "xm/error"
)

type Repository interface {
	GetAll(uow *UnitOfWork, out interface{}, queryProcessors []QueryProcessor) dbError.DatabaseError
	Get(uow *UnitOfWork, out interface{}, id uuid.UUID) dbError.DatabaseError
	Add(uow *UnitOfWork, out interface{}) dbError.DatabaseError
	Update(uow *UnitOfWork, out interface{}) dbError.DatabaseError
	Delete(uow *UnitOfWork, out interface{}, where ...interface{}) dbError.DatabaseError
}

// GormRepository implements Repository
type GormRepository struct {
}

// NewRepository returns a new repository object
func NewRepository() Repository {
	return &GormRepository{}
}

// QueryProcessor allows to modify the query before it is executed
type QueryProcessor func(db *gorm.DB, out interface{}) (*gorm.DB, dbError.DatabaseError)

// Filter will filter the results
func Filter(condition string, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, dbError.DatabaseError) {
		db = db.Where(condition, args...)
		return db, nil
	}
}

// GetAll retrieves all the records for a specified entity and returns it
func (repository *GormRepository) GetAll(uow *UnitOfWork, out interface{}, queryProcessors []QueryProcessor) dbError.DatabaseError {
	db := uow.DB

	if queryProcessors != nil {
		var err error
		for _, queryProcessor := range queryProcessors {
			db, err = queryProcessor(db, out)
			if err != nil {
				return dbError.NewDatabaseError(err)
			}
		}
	}
	if err := db.Find(out).Error; err != nil {
		return dbError.NewDatabaseError(err)
	}
	return nil
}

// Get a record for specified entity with specific id
func (repository *GormRepository) Get(uow *UnitOfWork, out interface{}, id uuid.UUID) dbError.DatabaseError {
	if err := uow.DB.First(out, "id = ?", id).Error; err != nil {
		return dbError.NewDatabaseError(err)
	}
	return nil
}

// Add specified Entity
func (repository *GormRepository) Add(uow *UnitOfWork, entity interface{}) dbError.DatabaseError {
	if err := uow.DB.Create(entity).Error; err != nil {
		return dbError.NewDatabaseError(err)
	}
	return nil
}

// Update specified Entity
func (repository *GormRepository) Update(uow *UnitOfWork, entity interface{}) dbError.DatabaseError {
	if err := uow.DB.Model(entity).Updates(entity).Error; err != nil {
		return dbError.NewDatabaseError(err)
	}
	return nil
}

// Delete specified Entity
func (repository *GormRepository) Delete(uow *UnitOfWork, entity interface{}, where ...interface{}) dbError.DatabaseError {
	if err := uow.DB.Delete(entity, where...).Error; err != nil {
		return dbError.NewDatabaseError(err)
	}
	return nil
}

package base

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type BaseRepo[T any] struct {
	DB                     *gorm.DB `wire:"auto"`
	DefaultProcessDB       func(db *gorm.DB) *gorm.DB
	DefaultFindProcessDB   func(db *gorm.DB) *gorm.DB
	DefaultCreateProcessDB func(db *gorm.DB) *gorm.DB
	DefaultUpdateProcessDB func(db *gorm.DB) *gorm.DB
	FindOmitColumns        []string
	CreateOmitColumns      []string
	UpdateOmitColumns      []string
	DBStatement            *gorm.Statement
}

type CrudRepoInf[T any] interface {
	Count(processDB func(*gorm.DB) *gorm.DB) (int, error)

	Create(entity *T) (*T, error)
	Update(entity *T) (*T, error)
	Save(entity *T) (*T, error)
	SaveAll(entities []*T) ([]*T, error)

	FindById(id any) (*T, error)
	FindAll(processDB func(*gorm.DB) *gorm.DB) ([]*T, error)

	DeleteById(id any) error
	DeleteByIds(ids []any) error
}

func (br *BaseRepo[T]) MustInitOnce() {
	var t T
	br.DBStatement = &gorm.Statement{DB: br.DB}
	if err := br.DBStatement.Parse(&t); err != nil {
		panic(fmt.Sprintf("statement parse model error :%v", err))
	}

	br.DefaultProcessDB = func(db *gorm.DB) *gorm.DB {
		return db
	}

	br.DefaultFindProcessDB = func(db *gorm.DB) *gorm.DB {
		return db.Omit(append(db.Statement.Omits, br.FindOmitColumns...)...)
	}

	br.DefaultCreateProcessDB = func(db *gorm.DB) *gorm.DB {
		return db.Omit(append(db.Statement.Omits, br.CreateOmitColumns...)...)
	}

	br.UpdateOmitColumns = append(br.UpdateOmitColumns, "CreatedAt")
	br.DefaultUpdateProcessDB = func(db *gorm.DB) *gorm.DB {
		return db.Omit(append(db.Statement.Omits, br.UpdateOmitColumns...)...)
	}
}

func (br *BaseRepo[T]) Count(processDB func(*gorm.DB) *gorm.DB) (int, error) {
	processDB = br.CheckProcessDB(processDB)

	var m T
	var count int64
	err := br.DB.Transaction(func(tx *gorm.DB) error {
		return processDB(tx).Model(&m).Count(&count).Error
	})
	return int(count), err
}

func (br *BaseRepo[T]) CheckProcessDB(processDB func(*gorm.DB) *gorm.DB) func(*gorm.DB) *gorm.DB {
	if processDB == nil {
		if br.DefaultProcessDB == nil {
			panic("BaseRepo's MustInitOnce must be call first.")
		}
		return br.DefaultProcessDB
	}

	return processDB
}

func (br *BaseRepo[T]) Save(entity *T) (*T, error) {
	return br.SaveWithDB(entity, nil)
}

func (br *BaseRepo[T]) Create(entity *T) (*T, error) {
	if idInf, ok := any(entity).(IdInf); ok {
		if err := idInf.SetId(nil); err != nil {
			return nil, err
		}
	}
	return br.SaveWithDB(entity, br.DefaultCreateProcessDB)
}

func (br *BaseRepo[T]) Update(entity *T) (*T, error) {
	return br.SaveWithDB(entity, br.DefaultUpdateProcessDB)
}

func (br *BaseRepo[T]) SaveWithDB(entity *T, processDB func(*gorm.DB) *gorm.DB) (*T, error) {
	processDB = br.CheckProcessDB(processDB)

	err := br.DB.Transaction(func(tx *gorm.DB) error {
		return processDB(tx).Save(entity).Error
	})
	return entity, err
}

func (br *BaseRepo[T]) SaveAll(entities []*T) ([]*T, error) {
	return br.SaveAllWithDB(entities, nil)
}

func (br *BaseRepo[T]) SaveAllWithDB(entities []*T, processDB func(*gorm.DB) *gorm.DB) ([]*T, error) {
	processDB = br.CheckProcessDB(processDB)

	err := br.DB.Transaction(func(tx *gorm.DB) error {
		return processDB(tx).Save(entities).Error
	})
	return entities, err
}

func (br *BaseRepo[T]) FindById(id any) (*T, error) {
	return br.FindByIdWithDB(id, br.DefaultFindProcessDB)
}

func (br *BaseRepo[T]) FindByIdWithDB(id any, processDB func(*gorm.DB) *gorm.DB) (*T, error) {
	processDB = br.CheckProcessDB(processDB)

	return br.FindWithDB(func(db *gorm.DB) *gorm.DB {
		return processDB(db).Where("id = ?", id)
	})
}

func (br *BaseRepo[T]) FindWithDB(processDB func(*gorm.DB) *gorm.DB) (*T, error) {
	processDB = br.CheckProcessDB(processDB)

	var entity T

	err := br.DB.Transaction(func(tx *gorm.DB) error {
		return processDB(tx).First(&entity).Error
	})

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (br *BaseRepo[T]) FindAll(processDB func(*gorm.DB) *gorm.DB) ([]*T, error) {
	processDB = br.CheckProcessDB(processDB)

	var entities []*T
	err := br.DB.Transaction(func(tx *gorm.DB) error {
		return br.DefaultFindProcessDB(processDB(tx)).Find(&entities).Error
	})
	return entities, err
}

func (br *BaseRepo[T]) DeleteById(id any) error {
	return br.DeleteByIdWithDB(id, nil)
}

func (br *BaseRepo[T]) UpdateById(id any, column string, value interface{}) error {
	return br.UpdateWithDB(column, value, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})
}

func (br *BaseRepo[T]) UpdatesWithDB(entity *T, processDB func(*gorm.DB) *gorm.DB) error {
	processDB = br.CheckProcessDB(processDB)

	return br.DB.Transaction(func(tx *gorm.DB) error {
		return processDB(tx).Updates(entity).Error
	})
}

func (br *BaseRepo[T]) UpdateWithDB(column string, value interface{}, processDB func(*gorm.DB) *gorm.DB) error {
	processDB = br.CheckProcessDB(processDB)

	var entity T
	return br.DB.Transaction(func(tx *gorm.DB) error {
		return processDB(tx).Model(&entity).Update(column, value).Error
	})
}

func (br *BaseRepo[T]) DeleteByIdWithDB(id any, processDB func(*gorm.DB) *gorm.DB) error {
	processDB = br.CheckProcessDB(processDB)

	var entity T
	if baseInf, ok := any(&entity).(IdInf); ok {
		// 此处代码配套BeforeDelete，可以在BeforeDelete方法的对象中获取到id
		if err := baseInf.SetId(id); err != nil {
			return err
		}
		return br.DB.Transaction(func(tx *gorm.DB) error {
			result := processDB(tx).Delete(&entity)

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return errors.New("delete error, record does not exist")
			}

			return nil
		})
	} else {
		return br.DB.Transaction(func(tx *gorm.DB) error {
			return processDB(tx).Delete(&entity, "id = ?", id).Error
		})
	}
}

func (br *BaseRepo[T]) DeleteByIds(id []any) error {
	return br.DeleteByIdsWithDB(id, nil)
}

func (br *BaseRepo[T]) DeleteByIdsWithDB(ids []any, processDB func(*gorm.DB) *gorm.DB) error {
	processDB = br.CheckProcessDB(processDB)

	var entity T
	if idInf, ok := any(&entity).(IdInf); ok {
		// 此处代码配套BeforeDelete，可以在BeforeDelete方法的对象中获取到id
		return br.DB.Transaction(func(tx *gorm.DB) error {
			rowsAffected := 0

			for _, id := range ids {
				if err := idInf.SetId(id); err != nil {
					return err
				}

				result := processDB(tx).Delete(&entity)

				if result.Error != nil {
					return result.Error
				}

				rowsAffected += int(result.RowsAffected)
			}

			if rowsAffected == 0 {
				return errors.New("delete error, record does not exist")
			}

			return nil
		})
	} else {
		return br.DB.Transaction(func(tx *gorm.DB) error {
			return processDB(tx).Delete(&entity, "id in ?", ids).Error
		})
	}
}

func (br *BaseRepo[T]) IncrementCounterFieldByIds(ids []any, fieldName string, delta int) error {
	return br.IncrementCounterField(fieldName, delta, func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", ids)
	})
}

func (br *BaseRepo[T]) IncrementCounterField(fieldName string, delta int, processDB func(*gorm.DB) *gorm.DB) error {
	if processDB == nil {
		return errors.New("processDB must be set")
	}

	field := br.DBStatement.Schema.LookUpField(fieldName)
	if field == nil {
		return errors.New(fmt.Sprintf("filed '%s' not exist.", fieldName))
	}

	fieldIncrement := gorm.Expr(field.DBName+" + ?", delta)

	entity := new(T)

	return br.DB.Transaction(func(tx *gorm.DB) error {
		return processDB(tx.Model(entity)).Update(fieldName, fieldIncrement).Error
	})
}

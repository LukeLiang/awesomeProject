package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Transaction 执行事务
func Transaction(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := fn(tx); err != nil {
			Logger.Error("Transaction failed, rolling back",
				zap.Error(err),
			)
			return err
		}
		return nil
	})
}

// TransactionWithResult 执行事务并返回结果
func TransactionWithResult[T any](fn func(tx *gorm.DB) (T, error)) (T, error) {
	var result T
	err := DB.Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = fn(tx)
		if err != nil {
			Logger.Error("Transaction failed, rolling back",
				zap.Error(err),
			)
			return err
		}
		return nil
	})
	return result, err
}

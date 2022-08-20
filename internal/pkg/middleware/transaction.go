package middleware

import (
	"context"

	impaconstant "github.com/thnthien/impa/constant"
	"gorm.io/gorm"
)

func NewGormTransaction(db *gorm.DB, ctx context.Context, f func(ctx context.Context) error) error {
	tx := db.Begin()
	txCtx := context.WithValue(ctx, impaconstant.CtxDBKey, tx)

	err := f(txCtx)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}

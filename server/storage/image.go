package storage

import (
	"context"
	"io"
)

type Image interface {
	// 画像を保存する
	SaveImage(ctx context.Context, name string, image io.Reader) error
	// 一時URLを発行する。
	// 該当する名前のファイルが存在しなかったら ErrImageNotFound を返す。
	GetTemporalyURL(ctx context.Context, name string) (string, error)
}

package storage

import "context"

type Storage interface {
	Set(context.Context, string, []byte) error
	Get(context.Context, string) ([]byte, error)
	Exists(context.Context, string) (bool, error)
}


package file

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type FilesStorage struct {
	cache sync.Map
	dir   string
}

func (s *FilesStorage) Exists(ctx context.Context, key string) (bool, error) {
	cachedFileName := fmt.Sprintf("%s/%s", s.dir, key)
	cacheContent, err := ioutil.ReadFile(cachedFileName)
	if err != nil {
		return false, err
	}
	s.cache.Store(key, cacheContent)
	return true, nil
}

func (s *FilesStorage) Get(ctx context.Context, key string) ([]byte, error) {
	cachedFileName := fmt.Sprintf("%s/%s", s.dir, key)
	v, ok := s.cache.Load(cachedFileName)

	if !ok {
		cacheContent, err := ioutil.ReadFile(cachedFileName)
		s.cache.Store(cachedFileName, cacheContent)
		if err != nil {
			return []byte{}, err
		} else {
			log.Printf("Key: %s in file cache, serving it\n", key)
			return cacheContent, nil
		}
	}
	log.Printf("Key: %s in memory cache, serving it\n", key)
	content, assertOk := v.([]byte)
	if assertOk {
		return content, nil
	}
	return []byte{}, fmt.Errorf("🤷")
}

func (s *FilesStorage) Set(ctx context.Context, key string, value []byte) error {
	cachedFileName := fmt.Sprintf("%s/%s", s.dir, key)
	s.cache.Store(cachedFileName, value)

	file, err := os.Create(cachedFileName)
	if err != nil {
		return fmt.Errorf("error when creating %s, err: %v", cachedFileName, err)
	}
	defer file.Close()
	_, err = file.Write(value)
	if err != nil {
		return fmt.Errorf("error when writing to %s, err: %v", cachedFileName, err)
	}
	return nil
}

func NewFileStorage(dir string) (*FilesStorage, error) {
	return &FilesStorage{dir: dir}, nil
}

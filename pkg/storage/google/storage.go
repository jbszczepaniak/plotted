package google

import (
	"context"

	"cloud.google.com/go/firestore"
)

type GoogleStorage struct {
	collection *firestore.CollectionRef
}

func NewGoogleStorage(ctx context.Context, projectID, collectionName string) (*GoogleStorage, error) {
	// tylko w developmencie
	//sa := option.WithCredentialsFile("/Users/jedrzejszczepaniak/.glcloud-secret/plotted-207513-5b6b79013df9.json")
	//client, err := firestore.NewClient(ctx, projectID, sa)

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	collection := client.Collection(collectionName)
	return &GoogleStorage{collection: collection}, nil
}

func (g *GoogleStorage) Set(ctx context.Context, key string, value []byte) error {
	data := make(map[string][]byte)
	data["data"] = value

	doc := g.collection.Doc(key)
	_, err := doc.Set(ctx, data)
	return err
}

func (g *GoogleStorage) Get(ctx context.Context, key string) ([]byte, error) {
	doc := g.collection.Doc(key)
	docSnapshot, err := doc.Get(ctx)
	if err != nil {
		return []byte{}, err
	}
	data := docSnapshot.Data()
	return data["data"].([]byte), nil
}

func (g *GoogleStorage) Exists(ctx context.Context, key string) (bool, error) {
	doc := g.collection.Doc(key)
	docSnapshot, err := doc.Get(ctx)
	if err != nil {
		return false, err
	}
	return docSnapshot.Exists(), nil
}

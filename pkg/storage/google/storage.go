package google

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type GoogleStorage struct {
	collection *firestore.CollectionRef
}

func NewGoogleStorage(ctx context.Context, gaeCredentials, projectID, collectionName string) (*GoogleStorage, error) {
	clientOpts := []option.ClientOption{}
	if gaeCredentials != "" {
		// in loval development this allows connection with Google Storage
		clientOpts = append(clientOpts, option.WithCredentialsFile(gaeCredentials))
	}

	client, err := firestore.NewClient(ctx, projectID, clientOpts...)

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

package mappings

import (
	"context"
	"errors"

	"urlshortener/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repo - wrapper for storage access
type Repo struct {
	collection *mongo.Collection
}

// NewRepo ...
func NewRepo(collection *mongo.Collection) *Repo {
	return &Repo{
		collection: collection,
	}
}

var (
	// ErrNoURL - No mapping for given short url
	ErrNoURL = errors.New("No url found")
	// ErrAlreadyExists - Mapping for given short url already exists
	ErrAlreadyExists = errors.New("Given short url is already used in another mapping")
)

// getByFilter returns slice of Mappings with URL info
// filter - mongo bson
func (repo *Repo) getByFilter(filter interface{}) ([]*Mapping, error) {
	mappings := []*Mapping{}
	ctx := context.Background()
	cur, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result Mapping
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, &result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return mappings, nil
}

// GetByShortURL returns Mapping with URL info
func (repo *Repo) GetByShortURL(shortURL string) (*Mapping, error) {
	filter := bson.M{
		"short_url": shortURL,
	}
	mappings, err := repo.getByFilter(filter)
	if err != nil {
		return nil, err
	}
	if len(mappings) == 0 {
		return nil, ErrNoURL
	}
	return mappings[0], nil
}

// AddMapping adds given mapping. If no short url specified then generates new short url
func (repo *Repo) AddMapping(mapping *Mapping) (*Mapping, error) {
	mapping.Views = 0
	if mapping.ShortURL == "" {
		shortURL, err := utils.GenerateRandomString(6)
		if err != nil {
			return nil, err
		}
		mapping.ShortURL = shortURL
	}
	_, err := repo.GetByShortURL(mapping.ShortURL)
	if err == nil {
		return nil, ErrAlreadyExists
	}
	if err != nil && err != ErrNoURL {
		return nil, err
	}
	ctx := context.Background()
	_, err = repo.collection.InsertOne(ctx, mapping)
	if err != nil {
		return nil, err
	}
	return mapping, nil
}

// IncrementMappingViews increments views of given mapping
func (repo *Repo) IncrementMappingViews(mapping *Mapping) (*Mapping, error) {
	filter := bson.M{
		"short_url": mapping.ShortURL,
	}
	mapping.Views++
	ctx := context.Background()
	_, err := repo.collection.ReplaceOne(ctx, filter, mapping)
	if err != nil {
		return nil, err
	}
	return mapping, nil
}

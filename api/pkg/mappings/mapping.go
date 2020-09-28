package mappings

// Mapping - struct with info about (ShortUrl -> LongUrl) mapping
type Mapping struct {
	LongURL  string `json:"long_url" bson:"long_url"`
	ShortURL string `json:"short_url" bson:"short_url"`
	Views    uint64 `json:"views" bson:"views"`
}

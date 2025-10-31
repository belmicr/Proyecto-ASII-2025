package domain_hotels

type Hotel struct {
	ID            string   `json:"id" bson:"_id,omitempty"`
	Name          string   `json:"name" bson:"name"`
	City          string   `json:"city" bson:"city"`
	PricePerNight float64  `json:"price_per_night" bson:"price_per_night"`
	Stars         int      `json:"stars" bson:"stars"`
	Amenities     []string `json:"amenities" bson:"amenities"`
	OwnerID       string   `json:"owner_id" bson:"owner_id"`
}

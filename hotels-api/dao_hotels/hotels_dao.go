package dao_hotels

import "hotels/domain_hotels"

// Representación del documento en MongoDB
type Hotel struct {
	ID            string   `bson:"_id,omitempty"`
	Name          string   `bson:"name"`
	City          string   `bson:"city"`
	PricePerNight float64  `bson:"price_per_night"`
	Stars         int      `bson:"stars"`
	Amenities     []string `bson:"amenities"`
	OwnerID       string   `bson:"owner_id"`
}

type Hotels []Hotel

// Helpers de mapeo entre dominio y DAO (evitan repetir código en el repo)
func FromDomain(d domain_hotels.Hotel) Hotel {
	return Hotel{
		ID:            d.ID,
		Name:          d.Name,
		City:          d.City,
		PricePerNight: d.PricePerNight,
		Stars:         d.Stars,
		Amenities:     d.Amenities,
		OwnerID:       d.OwnerID,
	}
}

func (h Hotel) ToDomain() domain_hotels.Hotel {
	return domain_hotels.Hotel{
		ID:            h.ID,
		Name:          h.Name,
		City:          h.City,
		PricePerNight: h.PricePerNight,
		Stars:         h.Stars,
		Amenities:     h.Amenities,
		OwnerID:       h.OwnerID,
	}
}

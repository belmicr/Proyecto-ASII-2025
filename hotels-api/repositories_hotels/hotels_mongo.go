//go:build ignore

package repositories_hotels

import (
	"context"
	"fmt"
	"log"

	"hotels/dao_hotels"
	"hotels/domain_hotels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuración para conectar a Mongo
type MongoConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Database   string
	Collection string
}

type Mongo struct {
	client     *mongo.Client
	database   string
	collection string
}

const connectionURI = "mongodb://%s:%s"

// NewMongo: crea el cliente y lo deja listo
func NewMongo(cfg MongoConfig) *Mongo {
	cred := options.Credential{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	cli, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(fmt.Sprintf(connectionURI, cfg.Host, cfg.Port)).SetAuth(cred),
	)
	if err != nil {
		log.Panicf("error connecting to mongo DB: %v", err)
	}
	return &Mongo{
		client:     cli,
		database:   cfg.Database,
		collection: cfg.Collection,
	}
}

// helper: colección
func (m *Mongo) col() *mongo.Collection {
	return m.client.Database(m.database).Collection(m.collection)
}

// Create: inserta un hotel y devuelve el dominio con ID
func (m *Mongo) Create(ctx context.Context, h domain_hotels.Hotel) (domain_hotels.Hotel, error) {
	// Permitimos _id string (Hex) para que el dominio sea string
	dao := dao_hotels.FromDomain(h)
	if dao.ID == "" {
		// Generamos un ID simple string (podés cambiar a hex o UUID si querés)
		dao.ID = fmt.Sprintf("h_%d", newID())
	}
	_, err := m.col().InsertOne(ctx, dao)
	if err != nil {
		return domain_hotels.Hotel{}, fmt.Errorf("error creating document: %w", err)
	}
	return dao.ToDomain(), nil
}

// Update: actualiza por id (solo setea los campos que vienen)
func (m *Mongo) Update(ctx context.Context, id string, h domain_hotels.Hotel) (domain_hotels.Hotel, error) {
	update := bson.M{}
	if h.Name != "" {
		update["name"] = h.Name
	}
	if h.City != "" {
		update["city"] = h.City
	}
	if h.PricePerNight != 0 {
		update["price_per_night"] = h.PricePerNight
	}
	if h.Stars != 0 {
		update["stars"] = h.Stars
	}
	if h.Amenities != nil {
		update["amenities"] = h.Amenities
	}
	if h.OwnerID != "" {
		update["owner_id"] = h.OwnerID
	}
	if len(update) == 0 {
		// nada para actualizar; devolvemos el actual
		return m.GetByID(ctx, id)
	}

	res, err := m.col().UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return domain_hotels.Hotel{}, fmt.Errorf("error updating document: %w", err)
	}
	if res.MatchedCount == 0 {
		return domain_hotels.Hotel{}, fmt.Errorf("hotel %s not found", id)
	}
	return m.GetByID(ctx, id)
}

// GetByID: busca por _id string
func (m *Mongo) GetByID(ctx context.Context, id string) (domain_hotels.Hotel, error) {
	var dao dao_hotels.Hotel
	err := m.col().FindOne(ctx, bson.M{"_id": id}).Decode(&dao)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain_hotels.Hotel{}, fmt.Errorf("not found")
		}
		return domain_hotels.Hotel{}, fmt.Errorf("error finding document: %w", err)
	}
	return dao.ToDomain(), nil
}

// List: lista todos (o filtra por nombre/ciudad con q)
func (m *Mongo) List(ctx context.Context, q string) ([]domain_hotels.Hotel, error) {
	filter := bson.M{}
	if q != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": q, "$options": "i"}},
				{"city": bson.M{"$regex": q, "$options": "i"}},
			},
		}
	}

	cur, err := m.col().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting documents: %w", err)
	}
	defer cur.Close(ctx)

	var list []domain_hotels.Hotel
	for cur.Next(ctx) {
		var dao dao_hotels.Hotel
		if err := cur.Decode(&dao); err != nil {
			return nil, fmt.Errorf("error decoding document: %w", err)
		}
		list = append(list, dao.ToDomain())
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return list, nil
}

// newID: contador simple en memoria (para demo). Cambialo por uuid o similar si querés.
var _idCounter int64

func newID() int64 {
	_idCounter++
	return _idCounter
}

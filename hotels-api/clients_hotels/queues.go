package clients_hotels

type RabbitConfig struct {
	Host, Port, Username, Password, QueueName string
}

// Stub que cumple la interfaz y no falla si no hay Rabbit
type Rabbit struct{}

func NewRabbit(_ RabbitConfig) *Rabbit { return &Rabbit{} }
func (r *Rabbit) Publish(_ any) error  { return nil }

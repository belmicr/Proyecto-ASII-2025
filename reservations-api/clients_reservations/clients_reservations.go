package clients_reservations

type RabbitConfig struct {
	Host, Port, Username, Password, QueueName string
}

type Rabbit struct{}

func NewRabbit(_ RabbitConfig) *Rabbit { return &Rabbit{} }
func (r *Rabbit) Publish(_ any) error  { return nil }

package twixter

import (
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Nats struct {
	*nats.EncodedConn
}

func NewNats(config *viper.Viper, log *logrus.Logger) *Nats {
	nc, err := nats.Connect(config.GetString("nats.uri"))
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to nats")
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to nats")
	}

	return &Nats{ec}
}

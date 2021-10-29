package twixter

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Nats struct {
	nats *nats.EncodedConn
	log  *logrus.Logger
	subs []*nats.Subscription
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

	return &Nats{ec, log, make([]*nats.Subscription, 0)}
}

func (n *Nats) RegisterWorker(subject, queue string, handler nats.Handler) error {
	sub, err := n.nats.QueueSubscribe(subject, queue, handler)
	if err != nil {
		n.log.WithError(err).Fatal("Failed to subscribe to nats")
		return err
	}
	n.subs = append(n.subs, sub)

	return nil
}

func (n *Nats) Run() {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	n.log.Info("Shutting Down...")

	for _, sub := range n.subs {
		if err := sub.Unsubscribe(); err != nil {
			n.log.WithError(err).Error("Failed to unsubscribe")
		}
	}
}

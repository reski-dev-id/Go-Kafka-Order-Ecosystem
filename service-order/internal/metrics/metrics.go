package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	KafkaMessagesPublishedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_published_total",
			Help: "Total kafka messages published",
		},
	)

	KafkaMessagesConsumedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_consumed_total",
			Help: "Total kafka messages consumed",
		},
	)

	KafkaMessagesFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_failed_total",
			Help: "Total kafka messages failed",
		},
	)
)

func Init() {
	prometheus.MustRegister(
		KafkaMessagesPublishedTotal,
		KafkaMessagesConsumedTotal,
		KafkaMessagesFailedTotal,
	)
}

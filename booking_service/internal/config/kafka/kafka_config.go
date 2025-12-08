package kafkaconfig

type KafkaConfig struct {
	Address string `yaml:"address"`
	Topic   string `yaml:"topic"`
}

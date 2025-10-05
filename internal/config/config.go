package config

import (
	"log"

	"github.com/dandyZicky/opensky-collector/pkg/events"
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host    string `mapstructure:"host"`
		Port    string `mapstructure:"port"`
		User    string `mapstructure:"user"`
		Pass    string `mapstructure:"pass"`
		Name    string `mapstructure:"name"`
		SSLMode string `mapstructure:"sslmode"`
	} `mapstructure:"database"`
	Kafka struct {
		Producer struct {
			ClientID string `mapstructure:"producer_client_id"`
		} `mapstructure:"producer"`
		Consumer struct {
			GroupID      string `mapstructure:"consumer_group_id"`
			AutoOffReset string `mapstructure:"auto_offset_reset"`
		} `mapstructure:"consumer"`
		BootstrapServers string `mapstructure:"bootstrap_servers"`
		ClientID         string `mapstructure:"client_id"`
		Acks             string `mapstructure:"acks"`
		TopicRaw         string `mapstructure:"topic_raw"`
		TopicEnriched    string `mapstructure:"topic_enriched"`
	} `mapstructure:"kafka"`
	SSE struct {
		Port           string   `mapstructure:"port"`
		AllowedOrigins []string `mapstructure:"allowed_origins"`
	} `mapstructure:"sse"`
	OpenSky struct {
		BaseURL         string `mapstructure:"base_url"`
		AuthURL         string `mapstructure:"auth_url"`
		CredentialsFile string `mapstructure:"credentials_file"`
		TickerInterval  int    `mapstructure:"ticker_interval_ms"`
	} `mapstructure:"opensky"`
}

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables and defaults.")
		} else {
			log.Fatalf("Fatal error config file: %v \n", err)
		}
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v \n", err)
	}

	if AppConfig.Database.Host == "" {
		AppConfig.Database.Host = "localhost"
	}
	if AppConfig.Database.Port == "" {
		AppConfig.Database.Port = "5432"
	}
	if AppConfig.Database.User == "" {
		AppConfig.Database.User = "opensky-collector"
	}
	if AppConfig.Database.Pass == "" {
		AppConfig.Database.Pass = "collector-pass"
	}
	if AppConfig.Database.Name == "" {
		AppConfig.Database.Name = "telemetry"
	}
	if AppConfig.Database.SSLMode == "" {
		AppConfig.Database.SSLMode = "disable"
	}

	if AppConfig.Kafka.BootstrapServers == "" {
		AppConfig.Kafka.BootstrapServers = "localhost:29092"
	}
	if AppConfig.Kafka.Consumer.GroupID == "" {
		AppConfig.Kafka.Consumer.GroupID = "opensky-processor"
	}
	if AppConfig.Kafka.Producer.ClientID == "" {
		AppConfig.Kafka.Producer.ClientID = "openskyCollector"
	}
	if AppConfig.Kafka.Acks == "" {
		AppConfig.Kafka.Acks = "all"
	}
	if AppConfig.Kafka.TopicRaw == "" {
		AppConfig.Kafka.TopicRaw = "telemetry.raw"
	}
	if AppConfig.Kafka.TopicEnriched == "" {
		AppConfig.Kafka.TopicEnriched = "telemetry.enriched"
	}
	if AppConfig.Kafka.Consumer.AutoOffReset == "" {
		AppConfig.Kafka.Consumer.AutoOffReset = "earliest"
	}
	if AppConfig.SSE.Port == "" {
		AppConfig.SSE.Port = "8081"
	}
	if len(AppConfig.SSE.AllowedOrigins) == 0 {
		AppConfig.SSE.AllowedOrigins = []string{"http://localhost:3000"}
	}

	if AppConfig.OpenSky.BaseURL == "" {
		AppConfig.OpenSky.BaseURL = "https://opensky-network.org/api"
	}
	if AppConfig.OpenSky.AuthURL == "" {
		AppConfig.OpenSky.AuthURL = "https://auth.opensky-network.org/auth/realms/opensky-network/protocol/openid-connect/token"
	}
	if AppConfig.OpenSky.CredentialsFile == "" {
		AppConfig.OpenSky.CredentialsFile = "credentials.json"
	}
	if AppConfig.OpenSky.TickerInterval == 0 {
		AppConfig.OpenSky.TickerInterval = 21600
	}

	events.InitTopics(AppConfig.Kafka.TopicRaw, AppConfig.Kafka.TopicEnriched)
}

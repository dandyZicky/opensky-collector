package pg

const (
	DB       = "POSTGRES_DB"
	PASSWORD = "POSTGRES_PASSWORD"
	USER     = "POSTGRES_USER"
	PORT     = "POSTGRES_PORT"
	HOST     = "POSTGRES_HOST"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

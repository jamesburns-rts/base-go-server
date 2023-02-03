package config

import (
	"github.com/Netflix/go-env"
	"net/http"
	"os"
)

type (
	Application struct {
		Port            int    `env:"APPLICATION_PORT"`
		LocalHost       string `env:"APPLICATION_LOCAL_HOST"`
		LogLevel        string `env:"APPLICATION_LOG_LEVEL"`
		AuthEnabled     bool   `env:"APPLICATION_AUTH_ENABLED"`
		CORSOrigins     string `env:"APPLICATION_CORS_ORIGINS"`
		InitialPassword string `env:"APPLICATION_INITIAL_PASSWORD"`
		Database        Database
		JWT             JWT
		ExampleClient   ExampleClient
	}

	Database struct {
		Host       string `env:"DB_HOST"`
		Port       int    `env:"DB_PORT"`
		Database   string `env:"DB_DATABASE"`
		Username   string `env:"DB_USERNAME" json:"-"`
		Password   string `env:"DB_PASSWORD" json:"-"`
		SSLDisable bool   `env:"DB_DISABLE_SSL"`
	}

	JWT struct {
		RSAPrivateKey    string `env:"JWT_PRIVATE_KEY" json:"-"`
		RASPublicKey     string `env:"JWT_PUBLIC_KEY"`
		Lifespan         string `env:"JWT_LIFESPAN"`
		ExpirationLeeway string `env:"JWT_EXPIRATION_LEEWAY"`
	}

	ExampleClient struct {
		Url        string `env:"EXAMPLE_CLIENT_URL"`
		DebugLevel string `env:"EXAMPLE_DEBUG_LEVEL"`
		HttpClient *http.Client
	}
)

var Defaults = Application{
	Port:        8080,
	LocalHost:   "0.0.0.0",
	LogLevel:    "INFO",
	AuthEnabled: false,
	CORSOrigins: "localhost:*",
	Database: Database{
		Host:       "localhost",
		Port:       5432,
		Database:   "mydb",
		Username:   "localuser",
		Password:   "supersecret",
		SSLDisable: false,
	},
	JWT: JWT{
		RSAPrivateKey:    TestPrivateKey,
		RASPublicKey:     TestPublicKey,
		Lifespan:         "1h",
		ExpirationLeeway: "2m",
	},
	ExampleClient: ExampleClient{
		Url:        "https://pokeapi.co",
		HttpClient: http.DefaultClient,
	},
}

// ReadProperties Reads the properties from the environment variables
func ReadProperties() (app Application, err error) {

	app = Defaults

	_, err = env.UnmarshalFromEnviron(&app)
	if err != nil {
		return app, err
	}

	// if ssl wasn't explicitly set, and we are running db locally, just disable it for convenience
	if os.Getenv("DB_DISABLE_SSL") == "" {
		app.Database.SSLDisable = app.Database.Host == "127.0.0.1"
	}

	return app, err
}

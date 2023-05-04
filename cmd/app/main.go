package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/handler"
	"github.com/burenotti/rtu-it-lab-recruit/pkg/httpserver"
	"github.com/burenotti/rtu-it-lab-recruit/repositories"
	"github.com/burenotti/rtu-it-lab-recruit/services"
	"github.com/burenotti/rtu-it-lab-recruit/usecases"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Host               string
	Port               string
	AppName            string
	DbDsn              string
	SmtpHost           string
	SmtpPort           int
	SmtpUser           string
	SmtpPassword       string
	LoginCodeTTL       time.Duration
	AuthTokenTTL       time.Duration
	ActivationTokenTTL time.Duration
	PrivateKey         *rsa.PrivateKey
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *Config) HandlerConfig() *handler.Config {
	return &handler.Config{
		Name: c.AppName,
	}
}

func ParseConfig() *Config {
	viper.AutomaticEnv()
	viper.SetDefault("LOGIN_CODE_TTL", 2*time.Minute)
	viper.SetDefault("AUTH_TOKEN_TTL", 24*time.Hour)
	viper.SetDefault("ACTIVATION_TOKEN_TTL", 10*time.Minute)
	viper.SetDefault("PRIVATE_KEY_PATH", "private.pem")

	privateKey := ReadPrivateKeyFromFile(viper.GetString("PRIVATE_KEY_PATH"))

	cfg := Config{
		AppName:            "RTUITLab recruitment",
		DbDsn:              viper.GetString("DB_DSN"),
		LoginCodeTTL:       viper.GetDuration("LOGIN_CODE_TTL"),
		AuthTokenTTL:       viper.GetDuration("AUTH_TOKEN_TTL"),
		ActivationTokenTTL: viper.GetDuration("ACTIVATION_TOKEN_TTL"),
		SmtpHost:           viper.GetString("SMTP_HOST"),
		SmtpUser:           viper.GetString("SMTP_USER"),
		SmtpPort:           viper.GetInt("SMTP_PORT"),
		SmtpPassword:       viper.GetString("SMTP_PASSWORD"),
		PrivateKey:         privateKey,
	}
	flag.StringVar(&cfg.Host, "host", "0.0.0.0", "Server host")
	flag.StringVar(&cfg.Port, "port", "80", "Server port")
	flag.Parse()

	return &cfg
}

func getLogger() *logrus.Logger {
	logger := logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.JSONFormatter{
			PrettyPrint: true,
		},
		Level: logrus.InfoLevel,
	}
	return &logger
}

func ReadPrivateKeyFromFile(filename string) *rsa.PrivateKey {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		logrus.
			WithError(err).
			WithField("file_path", filename).
			Fatalf("can't read private key from file")
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		logrus.
			WithError(err).
			Fatalf("can't parse private key")
	}
	return privateKey
}

func getDatabase(dbDsn string, logger *logrus.Logger) *repositories.Database {
	db, err := sqlx.Connect("pgx", dbDsn)
	if err != nil {
		logger.WithError(err).Fatalf("can't connect to database: %v", err)
	}
	return repositories.NewDatabase(db)
}

func main() {
	logger := getLogger()
	cfg := ParseConfig()

	db := getDatabase(cfg.DbDsn, logger)

	userStore := repositories.NewUserRepository(db)
	loginCodeStore := &repositories.LoginCodeRepository{
		Db:      db,
		CodeTTL: cfg.LoginCodeTTL,
	}
	dialer := gomail.NewDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUser, cfg.SmtpPassword)
	dialer.SSL = true
	//delivery := &services.MailingService{
	//	Dialer: dialer,
	//}

	delivery := &services.ConsoleDelivery{Logger: logger}

	auth := &services.AuthService{
		TokenTTL:   cfg.ActivationTokenTTL,
		PrivateKey: cfg.PrivateKey,
	}

	activationRepo := &repositories.UserActivationRepository{
		PrivateKey: cfg.PrivateKey,
		TokenTTL:   cfg.ActivationTokenTTL,
	}

	orgRepo := repositories.NewOrganizationRepository(db)

	ucase := handler.UseCases{
		EmailSignInUseCase: usecases.EmailSignInUseCase{
			UserStore:      userStore,
			LoginCodeStore: loginCodeStore,
			Delivery:       delivery,
			Transactioner:  db,
			Auth:           auth,
		},
		SignUpUseCase: usecases.SignUpUseCase{
			UserRepo:       userStore,
			ActivationRepo: activationRepo,
			Delivery:       delivery,
			Logger:         logger,
		},
		OrganizationUseCase: usecases.OrganizationUseCase{
			Transactioner:       db,
			OrganizationStorage: orgRepo,
		},
		AuthService: services.AuthService{
			TokenTTL:   cfg.AuthTokenTTL,
			PrivateKey: cfg.PrivateKey,
		},
	}
	http := handler.New(ucase, cfg.HandlerConfig())
	logger.Infof("Server run on %s", cfg.Addr())
	srv := httpserver.New(cfg.Addr(), http.Handler(), logger)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-srv.Notify():
		if err != nil {
			logger.WithError(err).Errorf("Server exited with error: %v", err)
		} else {
			logger.Info("Server exited without errors")
		}
	case s := <-interrupt:
		logger.WithField("signal", s.String()).Infof("%s Signal caught. Shutdown", s.String())
		if err := srv.Shutdown(); err != nil {
			logger.WithError(err).Errorf("Server shutdowned with error: %v", err)
		}
	}

}

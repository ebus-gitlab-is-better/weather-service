package data

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"
	"weather-service/internal/biz"
	"weather-service/internal/conf"

	slog "log"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewKeyCloakAPI,
	NewKeycloak,
	NewWeatherRepo,
)

// Data структура для работы с базой данных
type Data struct {
	db       *gorm.DB //Реализация работы с базой данной через библиотеку gorm
	keycloak *KeycloakAPI
}

// NewData создания экземпляра для работы с базой данных
func NewData(
	c *conf.Data,
	logger log.Logger,
	db *gorm.DB,
	keycloak *KeycloakAPI,
) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db, keycloak: keycloak}, cleanup, nil
}

type contextTxKey struct{}

func NewTransaction(d *Data) biz.Transaction {
	return d
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func NewKeycloak(c *conf.Data) *gocloak.GoCloak {
	client := gocloak.NewClient(c.Keycloak.Hostname)
	restyClient := client.RestyClient()
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return client
}

// NewDB Подключаемся к бд и создаем экземпляр его
func NewDB(c *conf.Data) *gorm.DB {
	newLogger := logger.New(
		slog.New(os.Stdout, "\r\n", slog.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)
	log.Info("opening database connection ")
	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
			c.Database.Host,
			c.Database.User,
			c.Database.Database,
			c.Database.Password,
			c.Database.Port)), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	//Вызывается ошибка и краш, если соединения с бд не установлено
	if err != nil {
		log.Errorf("failed opening connection to postgres: %v", err)
		panic("failed to connect database")
	}
	return db
}

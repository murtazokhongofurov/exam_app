package db

import (
	"fmt"
	"exam/customer_service/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDBForSuite(cfg config.Config) (*sqlx.DB, func()) {
    psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.PostgresHost,
        cfg.PostgresPort,
        cfg.PostgresUser,
        cfg.PostgresPassword,
        cfg.PostgresDatabase,
    )

    connDb, err := sqlx.Connect("postgres", psqlString)
    if err != nil {
        return nil, func() {}
    }
    cleanUpfunc := func ()  {
        connDb.Close()
    }

    return connDb, cleanUpfunc
}


func ConnectToDB(cfg config.Config) (*sqlx.DB, error) {
    psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.PostgresHost,
        cfg.PostgresPort,
        cfg.PostgresUser,
        cfg.PostgresPassword,
        cfg.PostgresDatabase,
    )

    connDb, err := sqlx.Connect("postgres", psqlString)
    if err != nil {
        return nil, err
    }

    return connDb, nil
}
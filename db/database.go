package db

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"os"
)

var DB *sqlx.DB
var schema = `
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    login text,
    password text,
	created_at timestamp DEFAULT now()
);

CREATE TABLE wallets (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	currency text,
	money int,
	user_id uuid REFERENCES users(id) ON DELETE CASCADE,
	created_at timestamp DEFAULT now()
);

CREATE TABLE transactions (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	sender uuid REFERENCES wallets(id) ON DELETE CASCADE,
	recipient uuid REFERENCES wallets(id) ON DELETE CASCADE,
	money int,
	created_at timestamp DEFAULT now()
);
`

func Init() {
	// Connecting to the database and error handling
	var err error
	DB, err = sqlx.Connect("pgx", os.Getenv("DATABASE_URI"))
	if err != nil {
		panic(err)
	}

	// Creating a schema
	_, err = DB.Exec(schema)
	if err != nil {
		fmt.Println("An error occurred during the creation of the schema. All tables have probably already been created")
	}
}

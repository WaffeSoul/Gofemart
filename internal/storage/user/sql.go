package user

var (
	createSQL     = `insert into users(username, password) values ($1, $2)`
	findByNameSQL = "select * from users where username=$1"
	findByIDSQL   = "select * from users where id=$1"
	migrateSQL    = `CREATE TABLE IF NOT EXISTS users (
id serial PRIMARY KEY,
username VARCHAR(255) UNIQUE,
password VARCHAR(255)
);`
)

package withdrawal

var (
	createSQL       = `insert into withdrawals(user_id, order_number,sum,processed_at) values ($1, $2,$3, $4)`
	findByOrderSQL =  "select * from withdrawals where order_number=$1"
	findByUserIDSQL ="select * from withdrawals where user_id=$1"
	migrateSQL      = `CREATE TABLE IF NOT EXISTS withdrawals(
		order_number VARCHAR(255) PRIMARY KEY,
		user_id INTEGER,
		sum double precision,
		processed_at VARCHAR(255)
	);`
	dropSQL = `DROP TABLE IF EXISTS withdrawals;`
)

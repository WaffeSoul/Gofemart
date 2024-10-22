package order

var (
	createSQL       = `insert into orders(number, user_id, status,accrual, uploaded_at) values ($1, $2, $3,$4, $5)`
	updateSQL       = `UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`
	findByIDUserSQL = "select * from orders where user_id=$1"
	findByNumberSQL = "select * from orders where number=$1"
	migrateSQL      = `CREATE TABLE IF NOT EXISTS orders(
		number VARCHAR(255)  PRIMARY KEY,
		user_id INTEGER,
		status  VARCHAR(16),
		accrual double precision,
		uploaded_at VARCHAR(255)
	);`
	dropSQL = `DROP TABLE IF EXISTS orders;`
)

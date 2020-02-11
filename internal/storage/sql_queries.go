package storage

const (

	// todo stmt (при реконекте stmt слетают, так что я решил не нагружать код пока ими, при больших селектах можно писать хранимые процедуры, но stmt все равно вещь)

	getWallet      = `SELECT id,balance FROM "test".wallets WHERE id =$1;`
	getTransaction = `SELECT id,state,amount FROM "test".transactions WHERE id =$1;`

	//updateBalance = `BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ;	SELECT "test".transfer($1,$2,$3);COMMIT;`
	updateBalance = `SELECT "test".transfer($1,$2,$3);`
)

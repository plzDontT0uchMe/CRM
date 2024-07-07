package test_test

import (
	_ "github.com/lib/pq"
)

const goRoutines = 50
const iterations = 1000

/*func TestSQLDB(t *testing.T) {
	connectionString := fmt.Sprint("host=localhost port=5432 user=postgres password=0000 dbname=crm-authService sslmode=disable")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	var user models.Account
	user.Login = "123"
	user.Password = hash.GenerateHash("123$ecr3t")

	var wg sync.WaitGroup
	for i := 1; i <= goRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			for j := 1; j <= iterations; j++ {
				err = db.QueryRow("SELECT id FROM accounts WHERE login = $1 AND password = $2 LIMIT 1", user.Login, user.Password).Scan(&user.Id)
				if err != nil {
					t.Fatalf("Failed to query database: %v %v %v", i, j, err)
				}
			}
		}(i)
	}
	wg.Wait()
}*/

/*func TestPGXDB(t *testing.T) {
	connectionString := fmt.Sprint("host=localhost port=5432 user=postgres password=0000 dbname=crm-authService sslmode=disable")
	db, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	var user models.User
	user.Login = "123"
	user.Password = hash.GenerateHash("123$ecr3t")

	var wg sync.WaitGroup
	for i := 1; i <= goRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			for j := 1; j <= iterations; j++ {
				err = db.QueryRow(context.Background(), "SELECT id FROM users WHERE login = $1 AND password = $2 LIMIT 1", user.Login, user.Password).Scan(&user.Id)
				if err != nil {
					t.Fatalf("Failed to query database: %v %v %v", i, j, err)
				}
			}
		}(i)
	}
	wg.Wait()
}*/

/*func TestPGXPoolDB(t *testing.T) {
	connectionString := fmt.Sprint("host=localhost port=5432 user=postgres password=0000 dbname=crm-authService sslmode=disable")

	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	var user models.Account
	user.Login = "123"
	user.Password = hash.GenerateHash("123$ecr3t")

	var wg sync.WaitGroup
	for i := 1; i <= goRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			for j := 1; j <= iterations; j++ {
				err = pool.QueryRow(context.Background(), "SELECT id FROM accounts WHERE login = $1 AND password = $2 LIMIT 1", user.Login, user.Password).Scan(&user.Id)
				if err != nil {
					t.Fatalf("Failed to query database: %v %v %v", i, j, err)
				}
			}
		}(i)
	}
	wg.Wait()
}*/

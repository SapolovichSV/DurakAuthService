package postgre

import (
	"context"
	"os"
	"testing"

	"github.com/SapolovichSV/durak/auth/internal/config"
	"github.com/SapolovichSV/durak/auth/internal/logger"
	postgre "github.com/SapolovichSV/durak/auth/internal/storage/postgre/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

//docker run --name test-postgres --rm -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postges -e POSTGRES_USER=postgres -p 5432:5432 -d  postgres

func TestRepoPostgre_AddUser_TestContainer_OkCases(t *testing.T) {
	if os.Getenv("DOCKERTESTDB") == "" {
		t.Skip("Skip test: no docker test db available")
	}
	sqlQueryCountOfRows := `SELECT COUNT(*) FROM users;`
	testConnString := "postgres://postgres:postgres@localhost:5432/postgres"

	ctx := t.Context()

	testPgPool, err := pgxpool.New(ctx, testConnString)
	assert.NoError(t, err, "can't connect to test database")
	defer dropTableAndTypes(ctx, testPgPool)

	testLogger := logger.New(config.Config{LogLevel: -4})

	if err := testPgPool.Ping(ctx); err != nil {
		t.Fatalf("can't ping db with err: %s", err)
	}

	createEmptyTableUsers(t, ctx, testPgPool)
	defer testPgPool.Close()

	type fields struct {
		pgPool *pgxpool.Pool
		hasher Hasher
		logger logger.Logger
	}
	type args struct {
		ctx      context.Context
		email    string
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "DefaultCase",
			fields: func() fields {
				mockHasher := postgre.NewMockHasher(t)
				mockHasher.EXPECT().Hash("defaultPass").Return("defaultHashedPass", nil)
				return fields{
					pgPool: testPgPool,
					hasher: mockHasher,
					logger: testLogger.WithGroup("repoPostgre"),
				}
			}(),
			args: args{
				ctx:      ctx,
				email:    "default@mail.com",
				username: "defaultUsername",
				password: "defaultPass",
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPostgre := New(tt.fields.pgPool, tt.fields.hasher, tt.fields.logger)
			assert.NoErrorf(t, repoPostgre.AddUser(
				tt.args.ctx, tt.args.email, tt.args.username, tt.args.password,
			), "unexpected error at AddUser() :%w", err)
			res, err := testPgPool.Query(ctx, sqlQueryCountOfRows)
			if err != nil {
				t.Logf("error after query %s", err)
			}
			assert.NoErrorf(t, err, "unexpected error at test.Query()", err)
			var countOfRows int
			res.Scan(&countOfRows)
			if i != countOfRows {
				assert.Equal(t, countOfRows, i, "COUNT(*) must be equal count of table tests")
			}
		})
	}

}
func createEmptyTableUsers(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	// IF NOT EXISTS because in multiply test using old roles not deleting,(too lazy to delete them for me)
	sqlEnumUserRole := `CREATE TYPE user_role AS ENUM ('admin', 'user');`
	sqlEnumUserStatus := `CREATE TYPE user_status AS ENUM ('offline', 'online');`

	sqlCreateTableUsers := `CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    passwordHASH VARCHAR(255) NOT NULL,
    status user_status DEFAULT 'offline',
    user_role user_role DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
	batch := &pgx.Batch{}
	batch.Queue(sqlEnumUserRole)
	batch.Queue(sqlEnumUserStatus)
	batch.Queue(sqlCreateTableUsers)
	tx, err := pool.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		t.Fatal("can't start transaction to create users table")
	}

	if batchResults := tx.SendBatch(ctx, batch); batchResults.Close() != nil {
		t.Fatal("can't exec sql statements with err" + batchResults.Close().Error())
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal("can't commit transaction to create users table")
	}
}
func dropTableAndTypes(ctx context.Context, pool *pgxpool.Pool) error {
	sqlDropTable := `DROP TABLE users;`
	sqlDropUserStatus := `DROP TYPE user_status;`
	sqlDropUserRole := `DROP TYPE user_role;`
	batch := pgx.Batch{}
	batch.Queue(sqlDropTable)
	batch.Queue(sqlDropUserRole)
	batch.Queue(sqlDropUserStatus)
	results := pool.SendBatch(ctx, &batch)
	return results.Close()
}

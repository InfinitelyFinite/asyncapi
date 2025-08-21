package store_test

import (
	"asyncapi/fixtures"
	"asyncapi/store"
	"context"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestUserStore(t *testing.T) {
	env := fixtures.NewTestEnv(t)
	cleanup := env.SetupDb(t)
	t.Cleanup(func() {
		cleanup(t)
	})

	now := time.Now()
	ctx := context.Background()
	userStore := store.NewUserStore(env.Db)
	user, err := userStore.CreateUser(context.Background(), "test@test.com", "testingpassword")
	require.NoError(t, err)

	require.Equal(t, "test@test.com", user.Email)
	require.NoError(t, user.ComparePassword("testingpassword"))
	require.Less(t, now.UnixNano(), user.CreatedAt.UnixNano())

	user2, err := userStore.ById(ctx, user.Id)
	require.NoError(t, err)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Id, user2.Id)
	require.Equal(t, user.HashedPasswordBase64, user2.HashedPasswordBase64)
	require.Equal(t, user.CreatedAt.UnixNano(), user2.CreatedAt.UnixNano())

	user3, err := userStore.ByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, user.Email, user3.Email)
	require.Equal(t, user.Id, user3.Id)
	require.Equal(t, user.HashedPasswordBase64, user3.HashedPasswordBase64)
	require.Equal(t, user.CreatedAt.UnixNano(), user3.CreatedAt.UnixNano())

}

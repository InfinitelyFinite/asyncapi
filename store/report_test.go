package store_test

import (
	"asyncapi/fixtures"
	"asyncapi/store"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReportStore(t *testing.T) {
	env := fixtures.NewTestEnv(t)
	cleanup := env.SetupDb(t)
	t.Cleanup(func() { cleanup(t) })

	ctx := context.Background()
	reportStore := store.NewReportStore(env.Db)
	userStore := store.NewUserStore(env.Db)
	user, err := userStore.CreateUser(ctx, "test@test.com", "secretpswd")
	require.NoError(t, err)

	now := time.Now()
	report, err := reportStore.Create(ctx, user.Id, "monsters")
	require.NoError(t, err)
	require.Equal(t, user.Id, report.UserId)
	require.Equal(t, "monsters", report.ReportType)
	timeDiff := report.CreatedAt.Sub(now).Abs()
	require.Less(t, timeDiff, time.Second) // Ensure created within 1 second of now

	startedAt := report.CreatedAt.Add(time.Second)
	completedAt := report.CreatedAt.Add(2 * time.Second)
	failedAt := report.CreatedAt.Add(3 * time.Second)
	errorMsg := "there was a failure"
	downloadUrl := "http://localhost:8080/reports"
	outputPath := "s3://reports-test/reports"
	downloadUrlExpiresAt := report.CreatedAt.Add(4 * time.Second)

	report.ReportType = "food"
	report.StartedAt = &startedAt
	report.CompletedAt = &completedAt
	report.FailedAt = &failedAt
	report.DownloadUrl = &downloadUrl
	report.ErrorMessage = &errorMsg
	report.OutputFilePath = &outputPath
	report.DownloadUrlExpiresAt = &downloadUrlExpiresAt

	report2, err := reportStore.Update(ctx, report)
	require.NoError(t, err)

	require.Equal(t, report.UserId, report2.UserId)
	require.Equal(t, report.Id, report2.Id)
	require.Equal(t, "monsters", report2.ReportType)
	require.Equal(t, report.CreatedAt.UnixNano(), report2.CreatedAt.UnixNano())
	require.Equal(t, report.StartedAt.UnixNano(), report2.StartedAt.UnixNano())
	require.Equal(t, report.CompletedAt.UnixNano(), report2.CompletedAt.UnixNano())
	require.Equal(t, report.FailedAt.UnixNano(), report2.FailedAt.UnixNano())
	require.Equal(t, report.ErrorMessage, report2.ErrorMessage)
	require.Equal(t, report.DownloadUrl, report2.DownloadUrl)
	require.Equal(t, report.OutputFilePath, report2.OutputFilePath)
	require.Equal(t, report.DownloadUrlExpiresAt, report2.DownloadUrlExpiresAt)

	report3, err := reportStore.ByPrimaryKey(ctx, user.Id, report.Id)
	require.NoError(t, err)
	require.Equal(t, report.UserId, report3.UserId)
	require.Equal(t, report.Id, report3.Id)
	require.Equal(t, "monsters", report3.ReportType)
	require.Equal(t, report.CreatedAt.UnixNano(), report3.CreatedAt.UnixNano())
	require.Equal(t, report.StartedAt.UnixNano(), report3.StartedAt.UnixNano())
	require.Equal(t, report.CompletedAt.UnixNano(), report3.CompletedAt.UnixNano())
	require.Equal(t, report.FailedAt.UnixNano(), report3.FailedAt.UnixNano())
	require.Equal(t, report.ErrorMessage, report3.ErrorMessage)
	require.Equal(t, report.DownloadUrl, report3.DownloadUrl)
	require.Equal(t, report.OutputFilePath, report3.OutputFilePath)
	require.Equal(t, report.DownloadUrlExpiresAt, report3.DownloadUrlExpiresAt)
}

package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ReportStore struct {
	db *sqlx.DB
}

func NewReportStore(db *sql.DB) *ReportStore {
	return &ReportStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

type Report struct {
	UserId               uuid.UUID  `db:"user_id"`
	Id                   uuid.UUID  `db:"id"`
	ReportType           string     `db:"report_type"`
	OutputFilePath       *string    `db:"output_file_path"`
	DownloadUrl          *string    `db:"download_url"`
	DownloadUrlExpiresAt *time.Time `db:"download_url_expires_at"`
	ErrorMessage         *string    `db:"error_message"`
	CreatedAt            time.Time  `db:"created_at"`
	StartedAt            *time.Time `db:"started_at"`
	CompletedAt          *time.Time `db:"completed_at"`
	FailedAt             *time.Time `db:"failed_at"`
}

func (s *ReportStore) Create(ctx context.Context, userId uuid.UUID, reportType string) (*Report, error) {
	const insert = `INSERT INTO reports (user_id, report_type) VALUES ($1, $2) RETURNING *`
	var report Report
	if err := s.db.GetContext(ctx, &report, insert, userId, reportType); err != nil {
		return nil, fmt.Errorf("failed to insert report for user %s: %w", userId, err)
	}

	return &report, nil
}

func (s *ReportStore) Update(ctx context.Context, report *Report) (*Report, error) {
	const update = `UPDATE reports SET
                   output_file_path = $1,
                   download_url = $2,
                   download_url_expires_at = $3,
                   error_message = $4,
                   started_at = $5,
                   completed_at = $6,
                   failed_at = $7
                   WHERE user_id = $8 AND id = $9 RETURNING *`

	var updated Report

	if err := s.db.GetContext(ctx, &updated, update,
		report.OutputFilePath,
		report.DownloadUrl,
		report.DownloadUrlExpiresAt,
		report.ErrorMessage,
		report.StartedAt,
		report.CompletedAt,
		report.FailedAt,
		report.UserId,
		report.Id); err != nil {
		return nil, fmt.Errorf("failed to update report %s for user %s: %w", report.Id, report.UserId, err)
	}

	return &updated, nil
}

func (s *ReportStore) ByPrimaryKey(ctx context.Context, userId uuid.UUID, id uuid.UUID) (*Report, error) {
	const query = `SELECT * FROM reports WHERE user_id = $1 AND id = $2`
	var report Report
	if err := s.db.GetContext(ctx, &report, query, userId, id); err != nil {
		return nil, fmt.Errorf("failed to query report for user %s: %w", userId, err)
	}

	return &report, nil
}

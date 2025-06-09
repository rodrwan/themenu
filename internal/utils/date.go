package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ToPgDate convierte un time.Time a pgtype.Date
func ToPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

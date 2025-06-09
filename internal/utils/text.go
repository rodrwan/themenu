package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// ToPgText convierte un string a pgtype.Text
func ToPgText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}

package utils

import (
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

// ToFloat64 convierte un pgtype.Numeric a float64
func ToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	return float64(n.Int.Int64()) / math.Pow10(int(-n.Exp))
}

// ToPgNumeric convierte un float64 a pgtype.Numeric
func ToPgNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(f)
	return n
}

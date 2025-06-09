package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ToPgUUID convierte un uuid.UUID a pgtype.UUID
func ToPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

// FromPgUUID convierte un pgtype.UUID a uuid.UUID
func FromPgUUID(id pgtype.UUID) uuid.UUID {
	return id.Bytes
}

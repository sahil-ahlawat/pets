package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"pet-platform/backend/internal/domain"
)

type VetRepository interface {
	SearchVets(ctx context.Context, query, city string, emergencyOnly bool) ([]domain.Vet, error)
	FindNearbyVets(ctx context.Context, lat, lng, radiusKm float64) ([]domain.Vet, error)
}

type sqlVetRepo struct {
	db *DB
}

func NewVetRepository(db *DB) VetRepository {
	return &sqlVetRepo{db: db}
}

func (r *sqlVetRepo) SearchVets(ctx context.Context, query, city string, emergencyOnly bool) ([]domain.Vet, error) {
	baseQuery := `
		SELECT id, clinic_name, doctor_name, license_number, phone_number, emergency_phone, email, website_url,
		       address_line1, city, state, postal_code, country_code, latitude, longitude, is_24h_emergency,
		       specializations, rating, review_count
		FROM vets
		WHERE ($1 = '' OR clinic_name ILIKE '%' || $1 || '%' OR doctor_name ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR city ILIKE '%' || $2 || '%')
		  AND (NOT $3 OR is_24h_emergency = TRUE)
		ORDER BY rating DESC, review_count DESC
		LIMIT 50
	`

	rows, err := r.db.QueryContext(ctx, baseQuery, query, city, emergencyOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Vet
	for rows.Next() {
		var v domain.Vet
		var doc, lic, emergPhone, email, web, state, postal sql.NullString
		var lat, lng sql.NullFloat64
		var specs []string

		if err := rows.Scan(
			&v.ID, &v.ClinicName, &doc, &lic, &v.PhoneNumber, &emergPhone, &email, &web,
			&v.AddressLine1, &v.City, &state, &postal, &v.CountryCode, &lat, &lng,
			&v.Is24hEmergency, pq.Array(&specs), &v.Rating, &v.ReviewCount,
		); err != nil {
			return nil, err
		}

		if doc.Valid {
			v.DoctorName = doc.String
		}
		if lic.Valid {
			v.LicenseNumber = lic.String
		}
		if emergPhone.Valid {
			v.EmergencyPhone = emergPhone.String
		}
		if email.Valid {
			v.Email = email.String
		}
		if web.Valid {
			v.WebsiteURL = web.String
		}
		if state.Valid {
			v.State = state.String
		}
		if postal.Valid {
			v.PostalCode = postal.String
		}
		if lat.Valid {
			v.Latitude = &lat.Float64
		}
		if lng.Valid {
			v.Longitude = &lng.Float64
		}
		v.Specializations = specs
		list = append(list, v)
	}
	return list, nil
}

func (r *sqlVetRepo) FindNearbyVets(ctx context.Context, lat, lng, radiusKm float64) ([]domain.Vet, error) {
	if radiusKm <= 0 {
		radiusKm = 25
	}

	// Haversine formula calculation in SQL (works universally with standard Postgres, seamlessly enhanced with PostGIS)
	query := fmt.Sprintf(`
		SELECT id, clinic_name, doctor_name, license_number, phone_number, emergency_phone, email, website_url,
		       address_line1, city, state, postal_code, country_code, latitude, longitude, is_24h_emergency,
		       specializations, rating, review_count,
		       (6371 * acos(least(1.0, greatest(-1.0, cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + sin(radians($1)) * sin(radians(latitude)))))) AS distance_km
		FROM vets
		WHERE latitude IS NOT NULL AND longitude IS NOT NULL
		  AND (6371 * acos(least(1.0, greatest(-1.0, cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + sin(radians($1)) * sin(radians(latitude)))))) <= $3
		ORDER BY distance_km ASC
		LIMIT 50
	`)

	rows, err := r.db.QueryContext(ctx, query, lat, lng, radiusKm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Vet
	for rows.Next() {
		var v domain.Vet
		var doc, lic, emergPhone, email, web, state, postal sql.NullString
		var vlat, vlng sql.NullFloat64
		var dist float64
		var specs []string

		if err := rows.Scan(
			&v.ID, &v.ClinicName, &doc, &lic, &v.PhoneNumber, &emergPhone, &email, &web,
			&v.AddressLine1, &v.City, &state, &postal, &v.CountryCode, &vlat, &vlng,
			&v.Is24hEmergency, pq.Array(&specs), &v.Rating, &v.ReviewCount, &dist,
		); err != nil {
			return nil, err
		}

		if doc.Valid {
			v.DoctorName = doc.String
		}
		if lic.Valid {
			v.LicenseNumber = lic.String
		}
		if emergPhone.Valid {
			v.EmergencyPhone = emergPhone.String
		}
		if email.Valid {
			v.Email = email.String
		}
		if web.Valid {
			v.WebsiteURL = web.String
		}
		if state.Valid {
			v.State = state.String
		}
		if postal.Valid {
			v.PostalCode = postal.String
		}
		if vlat.Valid {
			v.Latitude = &vlat.Float64
		}
		if vlng.Valid {
			v.Longitude = &vlng.Float64
		}
		v.DistanceKm = &dist
		v.Specializations = specs
		list = append(list, v)
	}
	return list, nil
}

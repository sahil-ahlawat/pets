package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
)

type PetRepository interface {
	CreatePet(ctx context.Context, userID uuid.UUID, pet *domain.Pet, req *domain.CreatePetRequest) (*domain.PetDetail, error)
	ListPetsByUserID(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]domain.Pet, string, error)
	GetPetByID(ctx context.Context, petID uuid.UUID) (*domain.PetDetail, error)
	UpdatePet(ctx context.Context, petID uuid.UUID, req *domain.UpdatePetRequest) (*domain.PetDetail, error)
	DeletePet(ctx context.Context, petID uuid.UUID) error
	MarkPetLost(ctx context.Context, petID, userID uuid.UUID, req *domain.MarkLostRequest) error
	MarkPetReunited(ctx context.Context, petID, userID uuid.UUID) error
	UserHasPetAccess(ctx context.Context, userID, petID uuid.UUID) (bool, error)
}

type sqlPetRepo struct {
	db *DB
}

func NewPetRepository(db *DB) PetRepository {
	return &sqlPetRepo{db: db}
}

func generateSecureToken(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *sqlPetRepo) CreatePet(ctx context.Context, userID uuid.UUID, pet *domain.Pet, req *domain.CreatePetRequest) (*domain.PetDetail, error) {
	if pet.ID == uuid.Nil {
		pet.ID = uuid.New()
	}
	now := time.Now().UTC()
	pet.CreatedAt = now
	pet.UpdatedAt = now
	pet.Status = domain.PetStatusActive

	qrToken, err := generateSecureToken(16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secure QR token: %w", err)
	}

	err = r.db.WithTx(ctx, func(tx *sql.Tx) error {
		// 1. Insert into pets
		queryPet := `
			INSERT INTO pets (id, name, species, breed, gender, is_neutered, date_of_birth, weight_kg, avatar_url, status, is_archived, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, FALSE, $11, $12)
		`
		var dobVal interface{} = nil
		if req.DateOfBirth != nil && *req.DateOfBirth != "" {
			dobVal = *req.DateOfBirth
		}

		_, err := tx.ExecContext(ctx, queryPet,
			pet.ID, req.Name, req.Species, req.Breed, req.Gender,
			req.IsNeutered, dobVal, req.WeightKg, req.AvatarURL,
			pet.Status, now, now,
		)
		if err != nil {
			return fmt.Errorf("insert pet error: %w", err)
		}

		// 2. Link User as Primary Owner in pet_users
		queryPetUser := `
			INSERT INTO pet_users (id, pet_id, user_id, role, can_edit, created_at, updated_at)
			VALUES ($1, $2, $3, 'PRIMARY_OWNER', TRUE, $4, $5)
		`
		if _, err := tx.ExecContext(ctx, queryPetUser, uuid.New(), pet.ID, userID, now, now); err != nil {
			return fmt.Errorf("insert pet_user error: %w", err)
		}

		// 3. Create pet_profiles with privacy defaults
		queryProfile := `
			INSERT INTO pet_profiles (id, pet_id, bio, allergies, show_owner_name, show_owner_phone, show_microchip, show_allergies_emergency, show_vet_info, created_at, updated_at)
			VALUES ($1, $2, $3, $4, TRUE, FALSE, TRUE, TRUE, FALSE, $5, $6)
		`
		if _, err := tx.ExecContext(ctx, queryProfile, uuid.New(), pet.ID, req.Bio, req.Allergies, now, now); err != nil {
			return fmt.Errorf("insert pet_profile error: %w", err)
		}

		// 4. Create pet_emergency_profiles
		queryEmergency := `
			INSERT INTO pet_emergency_profiles (id, pet_id, emergency_contact_name, emergency_contact_phone, critical_medical_alert, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		if _, err := tx.ExecContext(ctx, queryEmergency, uuid.New(), pet.ID, req.EmergencyContactName, req.EmergencyContactPhone, req.Allergies, now, now); err != nil {
			return fmt.Errorf("insert emergency profile error: %w", err)
		}

		// 5. Create microchip ID if supplied
		if req.MicrochipNumber != "" {
			queryMicrochip := `
				INSERT INTO pet_identifications (id, pet_id, id_type, id_value, verified, created_at, updated_at)
				VALUES ($1, $2, 'MICROCHIP', $3, TRUE, $4, $5)
			`
			if _, err := tx.ExecContext(ctx, queryMicrochip, uuid.New(), pet.ID, req.MicrochipNumber, now, now); err != nil {
				return fmt.Errorf("insert microchip error: %w", err)
			}
		}

		// 6. Generate active QR Tag
		queryQR := `
			INSERT INTO pet_qr_codes (id, pet_id, token, status, tag_label, created_at, updated_at)
			VALUES ($1, $2, $3, 'ACTIVE', 'Collar Tag #1', $4, $5)
		`
		if _, err := tx.ExecContext(ctx, queryQR, uuid.New(), pet.ID, qrToken, now, now); err != nil {
			return fmt.Errorf("insert qr code error: %w", err)
		}

		// 7. Insert outbox event for async indexing and welcome push
		queryOutbox := `
			INSERT INTO outbox_events (id, event_type, aggregate_type, aggregate_id, payload, created_at)
			VALUES ($1, 'PetCreated', 'PET', $2, json_build_object('petId', $3::text, 'userId', $4::text, 'name', $5::text), $6)
		`
		if _, err := tx.ExecContext(ctx, queryOutbox, uuid.New(), pet.ID, pet.ID.String(), userID.String(), req.Name, now); err != nil {
			return fmt.Errorf("insert outbox error: %w", err)
		}

		// 8. Record initial timeline event
		queryTimeline := `
			INSERT INTO pet_events (id, pet_id, event_category, event_timestamp, title, description, icon_name)
			VALUES ($1, $2, 'PROFILE_CREATED', $3, 'Digital Identity Created', 'Registered digital ID and QR Tag', 'sparkles')
		`
		_, err = tx.ExecContext(ctx, queryTimeline, uuid.New(), pet.ID, now)
		return err
	})

	if err != nil {
		return nil, err
	}

	return r.GetPetByID(ctx, pet.ID)
}

func (r *sqlPetRepo) ListPetsByUserID(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]domain.Pet, string, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `
		SELECT p.id, p.name, p.species, p.breed, p.gender, p.is_neutered, p.date_of_birth,
		       p.weight_kg, p.avatar_url, p.status, p.is_archived, p.created_at, p.updated_at
		FROM pets p
		JOIN pet_users pu ON pu.pet_id = p.id
		WHERE pu.user_id = $1 AND p.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit+1)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var pets []domain.Pet
	for rows.Next() {
		var p domain.Pet
		var breed, gender, avatar, dob sql.NullString
		var weight sql.NullFloat64

		if err := rows.Scan(
			&p.ID, &p.Name, &p.Species, &breed, &gender, &p.IsNeutered, &dob,
			&weight, &avatar, &p.Status, &p.IsArchived, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, "", err
		}

		if breed.Valid {
			p.Breed = breed.String
		}
		if gender.Valid {
			p.Gender = gender.String
		}
		if avatar.Valid {
			p.AvatarURL = avatar.String
		}
		if dob.Valid {
			val := dob.String
			p.DateOfBirth = &val
		}
		if weight.Valid {
			p.WeightKg = &weight.Float64
		}
		pets = append(pets, p)
	}

	var nextCursor string
	if len(pets) > limit {
		nextCursor = pets[limit].ID.String()
		pets = pets[:limit]
	}

	return pets, nextCursor, nil
}

func (r *sqlPetRepo) GetPetByID(ctx context.Context, petID uuid.UUID) (*domain.PetDetail, error) {
	query := `
		SELECT p.id, p.name, p.species, p.breed, p.gender, p.is_neutered, p.date_of_birth,
		       p.weight_kg, p.avatar_url, p.status, p.is_archived, p.created_at, p.updated_at,
		       COALESCE(pp.bio, ''), COALESCE(pp.dietary_preferences, ''), COALESCE(pp.allergies, ''),
		       COALESCE(pp.chronic_conditions, ''), COALESCE(pp.show_owner_name, TRUE),
		       COALESCE(pp.show_owner_phone, FALSE), COALESCE(pp.show_microchip, TRUE),
		       COALESCE(pp.show_allergies_emergency, TRUE), COALESCE(pp.show_vet_info, FALSE),
		       COALESCE(pe.emergency_contact_name, ''), COALESCE(pe.emergency_contact_phone, ''),
		       COALESCE(pe.critical_medical_alert, ''),
		       COALESCE(pi.id_value, ''),
		       COALESCE(pq.token, '')
		FROM pets p
		LEFT JOIN pet_profiles pp ON pp.pet_id = p.id
		LEFT JOIN pet_emergency_profiles pe ON pe.pet_id = p.id
		LEFT JOIN pet_identifications pi ON pi.pet_id = p.id AND pi.id_type = 'MICROCHIP'
		LEFT JOIN pet_qr_codes pq ON pq.pet_id = p.id AND pq.status = 'ACTIVE'
		WHERE p.id = $1 AND p.deleted_at IS NULL
	`

	var d domain.PetDetail
	var breed, gender, avatar, dob sql.NullString
	var weight sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, petID).Scan(
		&d.ID, &d.Name, &d.Species, &breed, &gender, &d.IsNeutered, &dob,
		&weight, &avatar, &d.Status, &d.IsArchived, &d.CreatedAt, &d.UpdatedAt,
		&d.Bio, &d.DietaryPreferences, &d.Allergies, &d.ChronicConditions,
		&d.ShowOwnerName, &d.ShowOwnerPhone, &d.ShowMicrochip, &d.ShowAllergiesEmergency, &d.ShowVetInfo,
		&d.EmergencyContactName, &d.EmergencyContactPhone, &d.CriticalMedicalAlert,
		&d.MicrochipNumber, &d.ActiveQRToken,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if breed.Valid {
		d.Breed = breed.String
	}
	if gender.Valid {
		d.Gender = gender.String
	}
	if avatar.Valid {
		d.AvatarURL = avatar.String
	}
	if dob.Valid {
		val := dob.String
		d.DateOfBirth = &val
	}
	if weight.Valid {
		d.WeightKg = &weight.Float64
	}
	if d.ActiveQRToken != "" {
		d.QRURL = fmt.Sprintf("/p/%s", d.ActiveQRToken)
	}

	return &d, nil
}

func (r *sqlPetRepo) UpdatePet(ctx context.Context, petID uuid.UUID, req *domain.UpdatePetRequest) (*domain.PetDetail, error) {
	err := r.db.WithTx(ctx, func(tx *sql.Tx) error {
		// Update core pet table
		if req.Name != nil || req.Breed != nil || req.Gender != nil || req.IsNeutered != nil || req.WeightKg != nil || req.AvatarURL != nil {
			queryPet := `
				UPDATE pets SET
					name = COALESCE($2, name),
					breed = COALESCE($3, breed),
					gender = COALESCE($4, gender),
					is_neutered = COALESCE($5, is_neutered),
					weight_kg = COALESCE($6, weight_kg),
					avatar_url = COALESCE($7, avatar_url),
					updated_at = CURRENT_TIMESTAMP
				WHERE id = $1
			`
			if _, err := tx.ExecContext(ctx, queryPet, petID, req.Name, req.Breed, req.Gender, req.IsNeutered, req.WeightKg, req.AvatarURL); err != nil {
				return err
			}
		}

		// Update profile & privacy
		queryProfile := `
			UPDATE pet_profiles SET
				bio = COALESCE($2, bio),
				allergies = COALESCE($3, allergies),
				chronic_conditions = COALESCE($4, chronic_conditions),
				show_owner_name = COALESCE($5, show_owner_name),
				show_owner_phone = COALESCE($6, show_owner_phone),
				show_microchip = COALESCE($7, show_microchip),
				show_allergies_emergency = COALESCE($8, show_allergies_emergency),
				show_vet_info = COALESCE($9, show_vet_info),
				updated_at = CURRENT_TIMESTAMP
			WHERE pet_id = $1
		`
		if _, err := tx.ExecContext(ctx, queryProfile, petID, req.Bio, req.Allergies, req.ChronicConditions,
			req.ShowOwnerName, req.ShowOwnerPhone, req.ShowMicrochip, req.ShowAllergiesEmergency, req.ShowVetInfo); err != nil {
			return err
		}

		// Update emergency
		if req.EmergencyContactName != nil || req.EmergencyContactPhone != nil || req.CriticalMedicalAlert != nil {
			queryEmerg := `
				UPDATE pet_emergency_profiles SET
					emergency_contact_name = COALESCE($2, emergency_contact_name),
					emergency_contact_phone = COALESCE($3, emergency_contact_phone),
					critical_medical_alert = COALESCE($4, critical_medical_alert),
					updated_at = CURRENT_TIMESTAMP
				WHERE pet_id = $1
			`
			if _, err := tx.ExecContext(ctx, queryEmerg, petID, req.EmergencyContactName, req.EmergencyContactPhone, req.CriticalMedicalAlert); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return r.GetPetByID(ctx, petID)
}

func (r *sqlPetRepo) DeletePet(ctx context.Context, petID uuid.UUID) error {
	query := `UPDATE pets SET deleted_at = CURRENT_TIMESTAMP, is_archived = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, petID)
	return err
}

func (r *sqlPetRepo) MarkPetLost(ctx context.Context, petID, userID uuid.UUID, req *domain.MarkLostRequest) error {
	return r.db.WithTx(ctx, func(tx *sql.Tx) error {
		// Update pet status to LOST
		if _, err := tx.ExecContext(ctx, `UPDATE pets SET status = 'LOST', updated_at = CURRENT_TIMESTAMP WHERE id = $1`, petID); err != nil {
			return err
		}

		// Insert into lost_pet_reports
		queryReport := `
			INSERT INTO lost_pet_reports (id, pet_id, reported_by, last_seen_address, last_seen_city, circumstances, reward_amount, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE)
		`
		if _, err := tx.ExecContext(ctx, queryReport, uuid.New(), petID, userID, req.LastSeenAddress, req.LastSeenCity, req.Circumstances, req.RewardAmount); err != nil {
			return err
		}

		// Record status history
		if _, err := tx.ExecContext(ctx, `INSERT INTO pet_status_history (id, pet_id, previous_status, new_status, changed_by, notes) VALUES ($1, $2, 'ACTIVE', 'LOST', $3, $4)`,
			uuid.New(), petID, userID, "Pet marked lost by owner"); err != nil {
			return err
		}

		// Record in outbox for push broadcast to nearby community & email confirmation
		queryOutbox := `
			INSERT INTO outbox_events (id, event_type, aggregate_type, aggregate_id, payload)
			VALUES ($1, 'PetMarkedLost', 'PET', $2, json_build_object('petId', $3::text, 'address', $4::text))
		`
		if _, err := tx.ExecContext(ctx, queryOutbox, uuid.New(), petID, petID.String(), req.LastSeenAddress); err != nil {
			return err
		}

		// Timeline event
		queryTimeline := `
			INSERT INTO pet_events (id, pet_id, event_category, title, description, icon_name)
			VALUES ($1, $2, 'STATUS_CHANGE', '🚨 Marked as LOST', $3, 'alert-triangle')
		`
		desc := fmt.Sprintf("Last seen near %s. Lost alert broadcast activated.", req.LastSeenAddress)
		_, err := tx.ExecContext(ctx, queryTimeline, uuid.New(), petID, desc)
		return err
	})
}

func (r *sqlPetRepo) MarkPetReunited(ctx context.Context, petID, userID uuid.UUID) error {
	return r.db.WithTx(ctx, func(tx *sql.Tx) error {
		// Update pet status to ACTIVE
		if _, err := tx.ExecContext(ctx, `UPDATE pets SET status = 'ACTIVE', updated_at = CURRENT_TIMESTAMP WHERE id = $1`, petID); err != nil {
			return err
		}

		// Resolve any active lost reports
		if _, err := tx.ExecContext(ctx, `UPDATE lost_pet_reports SET is_active = FALSE, resolved_at = CURRENT_TIMESTAMP WHERE pet_id = $1 AND is_active = TRUE`, petID); err != nil {
			return err
		}

		// Record status history
		if _, err := tx.ExecContext(ctx, `INSERT INTO pet_status_history (id, pet_id, previous_status, new_status, changed_by, notes) VALUES ($1, $2, 'LOST', 'REUNITED', $3, $4)`,
			uuid.New(), petID, userID, "Pet marked safely reunited with family"); err != nil {
			return err
		}

		// Outbox event
		queryOutbox := `
			INSERT INTO outbox_events (id, event_type, aggregate_type, aggregate_id, payload)
			VALUES ($1, 'PetReunited', 'PET', $2, json_build_object('petId', $3::text))
		`
		if _, err := tx.ExecContext(ctx, queryOutbox, uuid.New(), petID, petID.String()); err != nil {
			return err
		}

		// Timeline event
		queryTimeline := `
			INSERT INTO pet_events (id, pet_id, event_category, title, description, icon_name)
			VALUES ($1, $2, 'STATUS_CHANGE', '🎉 Safely Reunited!', 'Pet returned safely home. Emergency mode deactivated.', 'heart')
		`
		_, err := tx.ExecContext(ctx, queryTimeline, uuid.New(), petID)
		return err
	})
}

func (r *sqlPetRepo) UserHasPetAccess(ctx context.Context, userID, petID uuid.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM pet_users WHERE user_id = $1 AND pet_id = $2`
	err := r.db.QueryRowContext(ctx, query, userID, petID).Scan(&count)
	return count > 0, err
}

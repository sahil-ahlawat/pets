package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
)

type QRRepository interface {
	GetQRCodesByPetID(ctx context.Context, petID uuid.UUID) ([]domain.PetQRCode, error)
	GenerateNewQRCode(ctx context.Context, petID uuid.UUID, tagLabel string) (*domain.PetQRCode, error)
	RevokeQRCode(ctx context.Context, qrID uuid.UUID, reason string) error
	GetPublicPetByToken(ctx context.Context, token string) (*domain.PublicPetProfile, *uuid.UUID, *uuid.UUID, *uuid.UUID, error)
	RecordScan(ctx context.Context, scan *domain.QRScan) error
	RecordFoundReport(ctx context.Context, token string, req *domain.ReportFoundRequest) error
}

type sqlQRRepo struct {
	db *DB
}

func NewQRRepository(db *DB) QRRepository {
	return &sqlQRRepo{db: db}
}

func (r *sqlQRRepo) GetQRCodesByPetID(ctx context.Context, petID uuid.UUID) ([]domain.PetQRCode, error) {
	query := `
		SELECT id, pet_id, token, status, tag_label, physical_tag_type, scan_count, last_scanned_at, created_at
		FROM pet_qr_codes
		WHERE pet_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, petID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.PetQRCode
	for rows.Next() {
		var q domain.PetQRCode
		var lastScanned sql.NullTime
		if err := rows.Scan(&q.ID, &q.PetID, &q.Token, &q.Status, &q.TagLabel, &q.PhysicalTagType, &q.ScanCount, &lastScanned, &q.CreatedAt); err != nil {
			return nil, err
		}
		if lastScanned.Valid {
			t := lastScanned.Time
			q.LastScannedAt = &t
		}
		q.PublicURL = fmt.Sprintf("/p/%s", q.Token)
		list = append(list, q)
	}
	return list, nil
}

func (r *sqlQRRepo) GenerateNewQRCode(ctx context.Context, petID uuid.UUID, tagLabel string) (*domain.PetQRCode, error) {
	token, err := generateSecureToken(16)
	if err != nil {
		return nil, err
	}

	if tagLabel == "" {
		tagLabel = "Collar Tag"
	}

	qr := &domain.PetQRCode{
		ID:              uuid.New(),
		PetID:           petID,
		Token:           token,
		Status:          "ACTIVE",
		TagLabel:        tagLabel,
		PhysicalTagType: "QR_COLLAR",
		CreatedAt:       time.Now().UTC(),
		PublicURL:       fmt.Sprintf("/p/%s", token),
	}

	query := `
		INSERT INTO pet_qr_codes (id, pet_id, token, status, tag_label, physical_tag_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	now := time.Now().UTC()
	_, err = r.db.ExecContext(ctx, query, qr.ID, qr.PetID, qr.Token, qr.Status, qr.TagLabel, qr.PhysicalTagType, now, now)
	if err != nil {
		return nil, err
	}
	return qr, nil
}

func (r *sqlQRRepo) RevokeQRCode(ctx context.Context, qrID uuid.UUID, reason string) error {
	query := `
		UPDATE pet_qr_codes
		SET status = 'REVOKED', revoked_at = CURRENT_TIMESTAMP, revocation_reason = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, qrID, reason)
	return err
}

// GetPublicPetByToken returns safe privacy-filtered DTO, plus petID, ownerID, and qrCodeID for notifications and scan recording
func (r *sqlQRRepo) GetPublicPetByToken(ctx context.Context, token string) (*domain.PublicPetProfile, *uuid.UUID, *uuid.UUID, *uuid.UUID, error) {
	query := `
		SELECT 
			p.id, p.name, p.species, p.breed, p.gender, p.avatar_url, p.status,
			u.id, u.first_name, u.last_name, u.phone_number, u.allow_direct_calls,
			pp.show_owner_name, pp.show_owner_phone, pp.show_microchip, pp.show_allergies_emergency,
			COALESCE(pp.allergies, ''),
			COALESCE(pe.critical_medical_alert, ''), COALESCE(pe.special_care_instructions, ''),
			COALESCE(pe.emergency_contact_name, ''), COALESCE(pe.emergency_contact_phone, ''),
			COALESCE(pi.id_value, ''),
			q.id
		FROM pet_qr_codes q
		JOIN pets p ON p.id = q.pet_id
		JOIN pet_users pu ON pu.pet_id = p.id AND pu.role = 'PRIMARY_OWNER'
		JOIN users u ON u.id = pu.user_id
		LEFT JOIN pet_profiles pp ON pp.pet_id = p.id
		LEFT JOIN pet_emergency_profiles pe ON pe.pet_id = p.id
		LEFT JOIN pet_identifications pi ON pi.pet_id = p.id AND pi.id_type = 'MICROCHIP'
		WHERE q.token = $1 AND q.status = 'ACTIVE' AND p.deleted_at IS NULL
	`

	var petID, ownerID, qrCodeID uuid.UUID
	var name, species string
	var breed, gender, avatar sql.NullString
	var status domain.PetStatus
	var ownerFirst, ownerLast string
	var ownerPhone sql.NullString
	var allowDirectCalls, showOwnerName, showOwnerPhone, showMicrochip, showAllergiesEmergency bool
	var allergies, criticalAlert, specialInstructions, emergencyName, emergencyPhone, microchipVal string

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&petID, &name, &species, &breed, &gender, &avatar, &status,
		&ownerID, &ownerFirst, &ownerLast, &ownerPhone, &allowDirectCalls,
		&showOwnerName, &showOwnerPhone, &showMicrochip, &showAllergiesEmergency,
		&allergies, &criticalAlert, &specialInstructions,
		&emergencyName, &emergencyPhone, &microchipVal, &qrCodeID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, nil, nil
	}
	if err != nil {
		return nil, nil, nil, nil, err
	}

	pub := &domain.PublicPetProfile{
		Token:            token,
		Name:             name,
		Species:          species,
		Status:           status,
		AllowDirectCalls: allowDirectCalls,
	}

	if breed.Valid {
		pub.Breed = breed.String
	}
	if gender.Valid {
		pub.Gender = gender.String
	}
	if avatar.Valid {
		pub.AvatarURL = avatar.String
	}

	// Privacy Controls: Owner Name
	if showOwnerName {
		pub.OwnerName = fmt.Sprintf("%s %s", ownerFirst, ownerLast)
	}

	// Privacy Controls: Direct Phone is only exposed if owner explicitly allowed both phone visibility AND direct calls
	if showOwnerPhone && allowDirectCalls && ownerPhone.Valid {
		pub.OwnerPhone = ownerPhone.String
	}

	// Privacy Controls: Microchip
	if showMicrochip && microchipVal != "" {
		pub.MicrochipNumber = microchipVal
	}

	// Privacy Controls: Allergies and Emergency Alerts
	if showAllergiesEmergency {
		pub.Allergies = allergies
		pub.CriticalMedicalAlert = criticalAlert
		pub.SpecialInstructions = specialInstructions
		pub.EmergencyContactName = emergencyName
		pub.EmergencyContactPhone = emergencyPhone
	}

	// If pet is currently LOST, retrieve lost details
	if status == domain.PetStatusLost {
		queryLost := `
			SELECT last_seen_address, last_seen_city, circumstances, reward_amount, lost_date
			FROM lost_pet_reports
			WHERE pet_id = $1 AND is_active = TRUE
			ORDER BY created_at DESC
			LIMIT 1
		`
		var ld domain.LostDetails
		var city, circum, reward sql.NullString
		if err := r.db.QueryRowContext(ctx, queryLost, petID).Scan(&ld.LastSeenAddress, &city, &circum, &reward, &ld.LostDate); err == nil {
			if city.Valid {
				ld.LastSeenCity = city.String
			}
			if circum.Valid {
				ld.Circumstances = circum.String
			}
			if reward.Valid {
				ld.RewardAmount = reward.String
			}
			pub.LostDetails = &ld
		}
	}

	return pub, &petID, &ownerID, &qrCodeID, nil
}

func (r *sqlQRRepo) RecordScan(ctx context.Context, scan *domain.QRScan) error {
	if scan.ID == uuid.Nil {
		scan.ID = uuid.New()
	}
	scan.ScannedAt = time.Now().UTC()

	return r.db.WithTx(ctx, func(tx *sql.Tx) error {
		// 1. Insert scan record
		queryScan := `
			INSERT INTO qr_scans (id, qr_code_id, pet_id, scanned_at, approx_latitude, approx_longitude, accuracy_meters, approx_address, city, country_code, ip_hash, user_agent)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`
		if _, err := tx.ExecContext(ctx, queryScan,
			scan.ID, scan.QRCodeID, scan.PetID, scan.ScannedAt,
			scan.ApproxLatitude, scan.ApproxLongitude, scan.AccuracyMeters,
			scan.ApproxAddress, scan.City, scan.CountryCode, scan.IPHash, scan.UserAgent,
		); err != nil {
			return err
		}

		// 2. Increment scan count on pet_qr_codes
		queryCount := `
			UPDATE pet_qr_codes 
			SET scan_count = scan_count + 1, last_scanned_at = $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
		if _, err := tx.ExecContext(ctx, queryCount, scan.QRCodeID, scan.ScannedAt); err != nil {
			return err
		}

		// 3. Emit Outbox Event so worker alerts owner via Push + Email
		queryOutbox := `
			INSERT INTO outbox_events (id, event_type, aggregate_type, aggregate_id, payload)
			VALUES ($1, 'QRScanned', 'QR', $2, json_build_object('petId', $3::text, 'scanId', $4::text))
		`
		_, err := tx.ExecContext(ctx, queryOutbox, uuid.New(), scan.QRCodeID, scan.PetID.String(), scan.ID.String())
		return err
	})
}

func (r *sqlQRRepo) RecordFoundReport(ctx context.Context, token string, req *domain.ReportFoundRequest) error {
	var petID, qrID, ownerID uuid.UUID
	queryLookup := `
		SELECT q.id, q.pet_id, pu.user_id
		FROM pet_qr_codes q
		JOIN pet_users pu ON pu.pet_id = q.pet_id AND pu.role = 'PRIMARY_OWNER'
		WHERE q.token = $1 AND q.status = 'ACTIVE'
	`
	err := r.db.QueryRowContext(ctx, queryLookup, token).Scan(&qrID, &petID, &ownerID)
	if err != nil {
		return fmt.Errorf("invalid or revoked QR tag: %w", err)
	}

	return r.db.WithTx(ctx, func(tx *sql.Tx) error {
		reportID := uuid.New()
		queryInsert := `
			INSERT INTO found_pet_reports (id, pet_id, qr_code_id, finder_name, finder_phone, message, found_address, found_latitude, found_longitude, photo_url, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP)
		`
		if _, err := tx.ExecContext(ctx, queryInsert, reportID, petID, qrID, req.FinderName, req.FinderPhone, req.Message, req.FoundAddress, req.Latitude, req.Longitude, req.PhotoURL); err != nil {
			return err
		}

		// Insert notification to owner
		notifTitle := "🚨 Someone Reported Finding Your Pet!"
		notifBody := fmt.Sprintf("A finder has scanned your pet's tag: \"%s\". Address: %s", req.Message, req.FoundAddress)
		queryNotif := `
			INSERT INTO notifications (id, user_id, pet_id, notification_type, title, body, is_read, created_at)
			VALUES ($1, $2, $3, 'PET_FOUND', $4, $5, FALSE, CURRENT_TIMESTAMP)
		`
		if _, err := tx.ExecContext(ctx, queryNotif, uuid.New(), ownerID, petID, notifTitle, notifBody); err != nil {
			return err
		}

		// Emit Outbox event
		queryOutbox := `
			INSERT INTO outbox_events (id, event_type, aggregate_type, aggregate_id, payload)
			VALUES ($1, 'PetFound', 'PET', $2, json_build_object('petId', $3::text, 'reportId', $4::text, 'ownerId', $5::text))
		`
		_, err := tx.ExecContext(ctx, queryOutbox, uuid.New(), petID, petID.String(), reportID.String(), ownerID.String())
		return err
	})
}

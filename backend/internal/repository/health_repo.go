package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
)

type HealthRepository interface {
	ListVaccinationTypes(ctx context.Context, species string) ([]domain.VaccinationType, error)
	CreateVaccination(ctx context.Context, v *domain.PetVaccination) error
	ListVaccinations(ctx context.Context, petID uuid.UUID) ([]domain.PetVaccination, error)
	CreateMedicalEvent(ctx context.Context, m *domain.PetMedicalEvent) error
	ListMedicalEvents(ctx context.Context, petID uuid.UUID) ([]domain.PetMedicalEvent, error)
	CreateMedication(ctx context.Context, m *domain.PetMedication) error
	ListMedications(ctx context.Context, petID uuid.UUID) ([]domain.PetMedication, error)
	GetTimeline(ctx context.Context, petID uuid.UUID, limit int) ([]domain.PetTimelineEvent, error)
	CreateDocument(ctx context.Context, doc *domain.PetDocument) error
	ListDocuments(ctx context.Context, petID uuid.UUID) ([]domain.PetDocument, error)
}

type sqlHealthRepo struct {
	db *DB
}

func NewHealthRepository(db *DB) HealthRepository {
	return &sqlHealthRepo{db: db}
}

func (r *sqlHealthRepo) ListVaccinationTypes(ctx context.Context, species string) ([]domain.VaccinationType, error) {
	query := `
		SELECT id, species, name, code, is_core, recommended_interval_months, description
		FROM vaccination_types
		WHERE ($1 = '' OR species = $1)
		ORDER BY is_core DESC, name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, species)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []domain.VaccinationType
	for rows.Next() {
		var vt domain.VaccinationType
		var desc sql.NullString
		if err := rows.Scan(&vt.ID, &vt.Species, &vt.Name, &vt.Code, &vt.IsCore, &vt.RecommendedIntervalMonths, &desc); err != nil {
			return nil, err
		}
		if desc.Valid {
			vt.Description = desc.String
		}
		types = append(types, vt)
	}
	return types, nil
}

func (r *sqlHealthRepo) CreateVaccination(ctx context.Context, v *domain.PetVaccination) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	v.CreatedAt = time.Now().UTC()

	return r.db.WithTx(ctx, func(tx *sql.Tx) error {
		query := `
			INSERT INTO pet_vaccinations (id, pet_id, vaccination_type_id, custom_vaccine_name, administered_at, next_due_date, vet_name, clinic_name, batch_lot_number, manufacturer, document_id, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		`
		var nextDue interface{} = nil
		if v.NextDueDate != nil && *v.NextDueDate != "" {
			nextDue = *v.NextDueDate
		}

		if _, err := tx.ExecContext(ctx, query,
			v.ID, v.PetID, v.VaccinationTypeID, v.VaccineName, v.AdministeredAt,
			nextDue, v.VetName, v.ClinicName, v.BatchLotNumber, v.Manufacturer,
			v.DocumentID, v.Notes, v.CreatedAt, v.CreatedAt,
		); err != nil {
			return err
		}

		// Insert into pet timeline events
		queryTimeline := `
			INSERT INTO pet_events (id, pet_id, event_category, event_timestamp, title, description, icon_name)
			VALUES ($1, $2, 'VACCINE', $3, $4, $5, 'shield-check')
		`
		title := fmt.Sprintf("💉 Vaccination: %s", v.VaccineName)
		desc := fmt.Sprintf("Administered at %s. Vet: %s", v.ClinicName, v.VetName)
		if _, err := tx.ExecContext(ctx, queryTimeline, uuid.New(), v.PetID, v.CreatedAt, title, desc); err != nil {
			return err
		}

		// Emit Outbox event
		queryOutbox := `
			INSERT INTO outbox_events (id, event_type, aggregate_type, aggregate_id, payload)
			VALUES ($1, 'VaccinationAdded', 'VACCINATION', $2, json_build_object('petId', $3::text, 'vaccineName', $4::text))
		`
		_, err := tx.ExecContext(ctx, queryOutbox, uuid.New(), v.ID, v.PetID.String(), v.VaccineName)
		return err
	})
}

func (r *sqlHealthRepo) ListVaccinations(ctx context.Context, petID uuid.UUID) ([]domain.PetVaccination, error) {
	query := `
		SELECT id, pet_id, vaccination_type_id, custom_vaccine_name, administered_at, next_due_date,
		       vet_name, clinic_name, batch_lot_number, manufacturer, document_id, notes, created_at
		FROM pet_vaccinations
		WHERE pet_id = $1
		ORDER BY administered_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, petID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.PetVaccination
	for rows.Next() {
		var v domain.PetVaccination
		var vacTypeID, docID sql.NullString
		var adminAt time.Time
		var nextDue sql.NullTime
		var vet, clinic, batch, mfg, notes sql.NullString

		if err := rows.Scan(
			&v.ID, &v.PetID, &vacTypeID, &v.VaccineName, &adminAt, &nextDue,
			&vet, &clinic, &batch, &mfg, &docID, &notes, &v.CreatedAt,
		); err != nil {
			return nil, err
		}

		v.AdministeredAt = adminAt.Format("2006-01-02")
		if nextDue.Valid {
			val := nextDue.Time.Format("2006-01-02")
			v.NextDueDate = &val
		}
		if vet.Valid {
			v.VetName = vet.String
		}
		if clinic.Valid {
			v.ClinicName = clinic.String
		}
		if batch.Valid {
			v.BatchLotNumber = batch.String
		}
		if mfg.Valid {
			v.Manufacturer = mfg.String
		}
		if notes.Valid {
			v.Notes = notes.String
		}
		if vacTypeID.Valid {
			id, _ := uuid.Parse(vacTypeID.String)
			v.VaccinationTypeID = &id
		}
		if docID.Valid {
			id, _ := uuid.Parse(docID.String)
			v.DocumentID = &id
		}
		list = append(list, v)
	}
	return list, nil
}

func (r *sqlHealthRepo) CreateMedicalEvent(ctx context.Context, m *domain.PetMedicalEvent) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	m.CreatedAt = time.Now().UTC()

	return r.db.WithTx(ctx, func(tx *sql.Tx) error {
		query := `
			INSERT INTO pet_medical_events (id, pet_id, event_type, event_date, title, diagnosis, treatment_summary, vet_name, clinic_name, cost_amount, cost_currency, document_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		`
		if _, err := tx.ExecContext(ctx, query,
			m.ID, m.PetID, m.EventType, m.EventDate, m.Title, m.Diagnosis,
			m.TreatmentSummary, m.VetName, m.ClinicName, m.CostAmount, m.CostCurrency,
			m.DocumentID, m.CreatedAt, m.CreatedAt,
		); err != nil {
			return err
		}

		// Timeline
		queryTimeline := `
			INSERT INTO pet_events (id, pet_id, event_category, event_timestamp, title, description, icon_name)
			VALUES ($1, $2, 'VET_VISIT', $3, $4, $5, 'stethoscope')
		`
		title := fmt.Sprintf("🩺 %s", m.Title)
		desc := fmt.Sprintf("Diagnosis: %s. Clinic: %s", m.Diagnosis, m.ClinicName)
		_, err := tx.ExecContext(ctx, queryTimeline, uuid.New(), m.PetID, m.CreatedAt, title, desc)
		return err
	})
}

func (r *sqlHealthRepo) ListMedicalEvents(ctx context.Context, petID uuid.UUID) ([]domain.PetMedicalEvent, error) {
	query := `
		SELECT id, pet_id, event_type, event_date, title, diagnosis, treatment_summary, vet_name, clinic_name, cost_amount, cost_currency, document_id, created_at
		FROM pet_medical_events
		WHERE pet_id = $1
		ORDER BY event_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, petID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.PetMedicalEvent
	for rows.Next() {
		var m domain.PetMedicalEvent
		var evDate time.Time
		var diag, treat, vet, clinic, docID, curr sql.NullString
		var cost sql.NullFloat64

		if err := rows.Scan(
			&m.ID, &m.PetID, &m.EventType, &evDate, &m.Title, &diag, &treat, &vet, &clinic, &cost, &curr, &docID, &m.CreatedAt,
		); err != nil {
			return nil, err
		}

		m.EventDate = evDate.Format("2006-01-02")
		if diag.Valid {
			m.Diagnosis = diag.String
		}
		if treat.Valid {
			m.TreatmentSummary = treat.String
		}
		if vet.Valid {
			m.VetName = vet.String
		}
		if clinic.Valid {
			m.ClinicName = clinic.String
		}
		if cost.Valid {
			m.CostAmount = &cost.Float64
		}
		if curr.Valid {
			m.CostCurrency = curr.String
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *sqlHealthRepo) CreateMedication(ctx context.Context, m *domain.PetMedication) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	now := time.Now().UTC()

	return r.db.WithTx(ctx, func(tx *sql.Tx) error {
		query := `
			INSERT INTO pet_medications (id, pet_id, medication_name, dosage, frequency, start_date, end_date, prescribed_by, reason, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, TRUE, $10, $11)
		`
		var endDate interface{} = nil
		if m.EndDate != nil && *m.EndDate != "" {
			endDate = *m.EndDate
		}

		if _, err := tx.ExecContext(ctx, query,
			m.ID, m.PetID, m.MedicationName, m.Dosage, m.Frequency,
			m.StartDate, endDate, m.PrescribedBy, m.Reason, now, now,
		); err != nil {
			return err
		}

		// Timeline
		queryTimeline := `
			INSERT INTO pet_events (id, pet_id, event_category, event_timestamp, title, description, icon_name)
			VALUES ($1, $2, 'MEDICATION', $3, $4, $5, 'pill')
		`
		title := fmt.Sprintf("💊 Medication: %s", m.MedicationName)
		desc := fmt.Sprintf("Dosage: %s (%s). Prescribed by %s", m.Dosage, m.Frequency, m.PrescribedBy)
		_, err := tx.ExecContext(ctx, queryTimeline, uuid.New(), m.PetID, now, title, desc)
		return err
	})
}

func (r *sqlHealthRepo) ListMedications(ctx context.Context, petID uuid.UUID) ([]domain.PetMedication, error) {
	query := `
		SELECT id, pet_id, medication_name, dosage, frequency, start_date, end_date, prescribed_by, reason, is_active
		FROM pet_medications
		WHERE pet_id = $1
		ORDER BY is_active DESC, start_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, petID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.PetMedication
	for rows.Next() {
		var m domain.PetMedication
		var start time.Time
		var end sql.NullTime
		var presc, reason sql.NullString

		if err := rows.Scan(
			&m.ID, &m.PetID, &m.MedicationName, &m.Dosage, &m.Frequency, &start, &end, &presc, &reason, &m.IsActive,
		); err != nil {
			return nil, err
		}

		m.StartDate = start.Format("2006-01-02")
		if end.Valid {
			val := end.Time.Format("2006-01-02")
			m.EndDate = &val
		}
		if presc.Valid {
			m.PrescribedBy = presc.String
		}
		if reason.Valid {
			m.Reason = reason.String
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *sqlHealthRepo) GetTimeline(ctx context.Context, petID uuid.UUID, limit int) ([]domain.PetTimelineEvent, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	query := `
		SELECT id, pet_id, event_category, event_timestamp, title, COALESCE(description, ''), COALESCE(icon_name, 'info')
		FROM pet_events
		WHERE pet_id = $1
		ORDER BY event_timestamp DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, petID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.PetTimelineEvent
	for rows.Next() {
		var ev domain.PetTimelineEvent
		if err := rows.Scan(&ev.ID, &ev.PetID, &ev.EventCategory, &ev.Timestamp, &ev.Title, &ev.Description, &ev.IconName); err != nil {
			return nil, err
		}
		events = append(events, ev)
	}
	return events, nil
}

func (r *sqlHealthRepo) CreateDocument(ctx context.Context, doc *domain.PetDocument) error {
	if doc.ID == uuid.Nil {
		doc.ID = uuid.New()
	}
	now := time.Now().UTC()
	doc.CreatedAt = now

	query := `
		INSERT INTO pet_documents (id, pet_id, uploaded_by, title, category, storage_key, file_url, mime_type, file_size_bytes, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, TRUE, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		doc.ID, doc.PetID, doc.UploadedBy, doc.Title, doc.Category,
		doc.StorageKey, doc.FileURL, doc.MimeType, doc.FileSizeBytes,
		now, now,
	)
	return err
}

func (r *sqlHealthRepo) ListDocuments(ctx context.Context, petID uuid.UUID) ([]domain.PetDocument, error) {
	query := `
		SELECT id, pet_id, uploaded_by, title, category, storage_key, file_url, mime_type, file_size_bytes, is_verified, created_at
		FROM pet_documents
		WHERE pet_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, petID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []domain.PetDocument
	for rows.Next() {
		var d domain.PetDocument
		if err := rows.Scan(&d.ID, &d.PetID, &d.UploadedBy, &d.Title, &d.Category, &d.StorageKey, &d.FileURL, &d.MimeType, &d.FileSizeBytes, &d.IsVerified, &d.CreatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

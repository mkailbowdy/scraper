package models

import (
	"database/sql"
	"errors"
	"time"
)

type Shigoto struct {
	ID             int64
	CompanyName    string
	JobTitle       string
	Category       string
	Location       string
	EmploymentType string
	Description    string
	JapaneseLevel  string
	EnglishLevel   string
	Sponsorship    bool
	CreatedAt      time.Time
}

type ShigotoModel struct {
	DB *sql.DB
}

func (m *ShigotoModel) Insert(companyName, jobTitle, category, location, employmentType, description, japaneseLevel, englishLevel string, sponsorship bool) (int, error) {
	var id int // Declare empty variable id here so that we have a place for Scan to copy the id value to.
	// Add RETURNING id at the end so that postgres can return the id of the newly created row.
	stmt := `INSERT INTO shigotos (company_name, job_title, category, location, employment_type, description, japanese_level, english_level, sponsorship, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP) RETURNING id`

	err := m.DB.QueryRow(stmt, companyName, jobTitle, category, location, employmentType, description, japaneseLevel, englishLevel, sponsorship).Scan(&id) // Instead of Exec, use QueryRow to run the statement and return the row. Use Scan to get the value of the id.
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *ShigotoModel) Get(id int) (Shigoto, error) {
	stmt := `SELECT id, company_name, job_title, category, location, employment_type, description, japanese_level, english_level, sponsorship, created_at FROM shigotos WHERE id = $1`
	row := m.DB.QueryRow(stmt, id)

	var s Shigoto

	// Scan() automatically converts the raw output from the row to the required Go types
	err := row.Scan(&s.ID, &s.CompanyName, &s.JobTitle, &s.Category, &s.Location, &s.EmploymentType, &s.Description, &s.JapaneseLevel, &s.EnglishLevel, &s.Sponsorship, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			/* Return a custom ErrNoRecord. The reason is to help encapsulate the model completely, so that our handlers
			aren’t concerned with the underlying datastore or reliant on datastore-specific errors (like sql.ErrNoRows)
			for its behavior. */
			return Shigoto{}, ErrNoRecord
		} else {
			return Shigoto{}, err
		}
	}

	return s, nil
}

func (m *ShigotoModel) Latest() ([]Shigoto, error) {
	stmt := `SELECT id, company_name, job_title, category, location, employment_type, description, japanese_level, english_level, sponsorship, created_at FROM shigotos ORDER BY created_at DESC`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var shigotos []Shigoto

	for rows.Next() {
		var s Shigoto
		err = rows.Scan(&s.ID, &s.CompanyName, &s.JobTitle, &s.Category, &s.Location, &s.EmploymentType, &s.Description, &s.JapaneseLevel, &s.EnglishLevel, &s.Sponsorship, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		shigotos = append(shigotos, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shigotos, nil
}

package repositories

import (
	"database/sql"
	"log"

	"github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
	"github.com/ros3n/hes/api/models"
)

type DBEmailsRepository struct {
	db           *sql.DB
	queryBuilder sq.StatementBuilderType
}

func NewDBEmailsRepository(dbAddr string) (*DBEmailsRepository, error) {
	db, err := sql.Open("postgres", dbAddr)
	if err != nil {
		return nil, err
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &DBEmailsRepository{db: db, queryBuilder: psql}, nil
}

func (der *DBEmailsRepository) Find(userID string, id int64) (*models.Email, error) {
	emailQuery := der.queryBuilder.Select("*").From("emails").Where(sq.Eq{"id": id, "user_id": userID})

	query, args, err := emailQuery.ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	email := models.Email{}

	err = der.db.
		QueryRow(query, args...).
		Scan(
			&email.ID, &email.UserID, &email.Sender, pq.Array(&email.Recipients), &email.Subject, &email.Message,
			&email.Status,
		)

	if err != nil {
		return nil, err
	}

	return &email, nil
}
func (der *DBEmailsRepository) FindByID(id int64) (*models.Email, error) {
	emailQuery := der.queryBuilder.Select("*").From("emails").Where(sq.Eq{"id": id})

	query, args, err := emailQuery.ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	email := models.Email{}

	err = der.db.
		QueryRow(query, args...).
		Scan(
			&email.ID, &email.UserID, &email.Sender, pq.Array(&email.Recipients), &email.Subject, &email.Message,
			&email.Status,
		)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &email, nil
}

func (der *DBEmailsRepository) Create(userID string, email *models.Email) (*models.Email, error) {
	email.UserID = userID
	createEmailQuery := der.queryBuilder.
		Insert("emails").
		Columns("user_id", "sender", "recipients", "subject", "message", "status").
		Values(email.UserID, email.Sender, pq.Array(email.Recipients), email.Subject, email.Message, email.Status).
		Suffix("RETURNING \"id\"")

	query, args, err := createEmailQuery.ToSql()
	if err != nil {
		log.Println(err)
		return nil, ErrCreateFailed
	}

	var emailID int64
	if err = der.db.QueryRow(query, args...).Scan(&emailID); err != nil {
		log.Println(err)
		return nil, ErrCreateFailed
	}

	email.ID = emailID

	return email, nil
}

func (der *DBEmailsRepository) Update(userID string, email *models.Email) (*models.Email, error) {
	// currently updating status is the only requirement
	createEmailQuery := der.queryBuilder.
		Update("emails").
		Set("status", email.Status).
		Where(sq.Eq{"id": email.ID, "user_id": userID})

	query, args, err := createEmailQuery.ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if _, err = der.db.Exec(query, args...); err != nil {
		log.Println(err)
		return nil, err
	}

	return email, nil
}

func (der *DBEmailsRepository) All(userID string) ([]*models.Email, error) {
	emailQuery := der.queryBuilder.Select("*").From("emails").Where(sq.Eq{"user_id": userID})

	query, args, err := emailQuery.ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := der.db.Query(query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	emails := make([]*models.Email, 0)
	for rows.Next() {
		email := models.Email{}
		err = rows.Scan(
			&email.ID, &email.UserID, &email.Sender, pq.Array(&email.Recipients), &email.Subject, &email.Message,
			&email.Status,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		emails = append(emails, &email)
	}

	return emails, nil
}

package pgx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/SlavaShagalov/ds-lab1/internal/models"
	pPersons "github.com/SlavaShagalov/ds-lab1/internal/persons"
	"github.com/SlavaShagalov/ds-lab1/internal/pkg/constants"
	pErrors "github.com/SlavaShagalov/ds-lab1/internal/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type repository struct {
 db  *sql.DB
	log  *zap.Logger
}

func New(db  *sql.DB, log *zap.Logger) pPersons.Repository {
	return &repository{
		db: db,
		log:  log,
	}
}

const createCmd = `
	INSERT INTO persons (name, age, address, work) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, age, address, work;`

func (repo *repository) Create(ctx context.Context, params *pPersons.CreateParams) (*models.Person, error) {
	row := repo.db.QueryRow(createCmd,
		params.Name,
		params.Age,
		params.Address,
		params.Work,
	)

	person := new(models.Person)
	err := scanPerson(row, person)
	if err != nil {
		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", createCmd))
		return nil, errors.Wrap(pErrors.ErrDb, err.Error())
	}

	repo.log.Debug("New person created", zap.Any("person", person))
	return person, nil
}

const getCmd = `
	SELECT id, name, age, address, work
	FROM persons
	WHERE id = $1;`

func (repo *repository) Get(ctx context.Context, id int64) (*models.Person, error) {
	row := repo.db.QueryRow(getCmd, id)

	person := new(models.Person)
	err := scanPerson(row, person)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(pErrors.ErrPersonNotFound, err.Error())
		}

		repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", getCmd),
			zap.Int64("id", id))
		return nil, errors.Wrap(pErrors.ErrDb, err.Error())
	}

	return person, nil
}

const listCmd = `
	SELECT id, name, age, address, work
	FROM persons
	offset $1`

func (repo *repository) List(ctx context.Context, offset, limit int64) ([]models.Person, error) {
	var err error
	var rows *sql.Rows
	query := listCmd
	if limit != 0 {
		query += " limit $2"
		rows, err = repo.db.Query(listCmd, offset, limit)
	} else {
		rows, err = repo.db.Query(listCmd, offset)
	}
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.String("sql_query", listCmd))
		return nil, errors.Wrap(pErrors.ErrDb, err.Error())
	}
	defer rows.Close()

	persons := []models.Person{}
	var person models.Person
	for rows.Next() {
		err = rows.Scan(
			&person.ID,
			&person.Name,
			&person.Age,
			&person.Address,
			&person.Work,
		)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", listCmd))
			return nil, errors.Wrap(pErrors.ErrDb, err.Error())
		}

		persons = append(persons, person)
	}

	return persons, nil
}

const fullUpdateCmd = `
	UPDATE persons
	SET %s
	WHERE id = $%d
	RETURNING id, name, age, address, work;`

func (repo *repository) PartialUpdate(ctx context.Context, params *pPersons.PartialUpdateParams) (*models.Person, error) {
	setValues := make([]string, 0, 4)
	args := make([]any, 0, 5)
	if params.Name != nil {
		setValue := fmt.Sprintf("name = $%d", len(args)+1)
		args = append(args, params.Name)
		setValues = append(setValues, setValue)
	}
	if params.Age != nil {
		setValue := fmt.Sprintf("age = $%d", len(args)+1)
		args = append(args, *params.Age)
		setValues = append(setValues, setValue)
	}
	if params.Address != nil {
		setValue := fmt.Sprintf("address = $%d", len(args)+1)
		args = append(args, *params.Address)
		setValues = append(setValues, setValue)
	}
	if params.Work != nil {
		setValue := fmt.Sprintf("work = $%d", len(args)+1)
		args = append(args, *params.Work)
		setValues = append(setValues, setValue)
	}
	if len(setValues) > 0 {
		setValuesPart := strings.Join(setValues, ", ")
		cmd := fmt.Sprintf(fullUpdateCmd, setValuesPart, len(args)+1)
		args = append(args, params.ID)

		row := repo.db.QueryRow(cmd, args...)
		person := new(models.Person)
		err := scanPerson(row, person)
		if err != nil {
			repo.log.Error(constants.DBScanError, zap.Error(err), zap.String("sql_query", fullUpdateCmd))
			return nil, errors.Wrap(pErrors.ErrDb, err.Error())
		}

		repo.log.Debug("Person partial updated", zap.Any("person", person))
		return person, nil
	}

	return nil, nil
}

const deleteCmd = `
	DELETE FROM persons 
	WHERE id = $1;`

func (repo *repository) Delete(ctx context.Context, id int64) error {
	result, err := repo.db.Exec(deleteCmd, id)
	if err != nil {
		repo.log.Error(constants.DBError, zap.Error(err), zap.Int64("id", id))
		return errors.Wrap(pErrors.ErrDb, err.Error())
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return pErrors.ErrPersonNotFound
	}

	repo.log.Debug("Person deleted", zap.Int64("id", id))
	return nil
}

func scanPerson(row pgx.Row, person *models.Person) error {
	return row.Scan(
		&person.ID,
		&person.Name,
		&person.Age,
		&person.Address,
		&person.Work,
	)
}

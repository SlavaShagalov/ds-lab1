package pgx

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/SlavaShagalov/ds-lab1/internal/models"
	pPersons "github.com/SlavaShagalov/ds-lab1/internal/persons"
	pkgErrors "github.com/SlavaShagalov/ds-lab1/internal/pkg/errors"
)

var err error
var logger *zap.Logger

func init() {
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreate(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		params  pPersons.CreateParams
		Person  models.Person
		err     error
	}

	const createCmd = `
	INSERT INTO persons (name, age, address, work) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, age, address, work;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "address", "work"})
				rows = rows.AddRow(1, "Johnny", 22, "Moscow, Red Square", "Yandex")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs("Johnny", 22, "Moscow, Red Square", "Yandex").
					WillReturnRows(rows)
			},
			params: pPersons.CreateParams{
				Name:    "Johnny",
				Age:     22,
				Address: "Moscow, Red Square",
				Work:    "Yandex",
			},
			Person: models.Person{
				ID:      1,
				Name:    "Johnny",
				Age:     22,
				Address: "Moscow, Red Square",
				Work:    "Yandex",
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(createCmd)).
					WithArgs("Johnny", 22, "Moscow, Red Square", "Yandex").
					WillReturnError(pkgErrors.ErrDb)
			},
			params: pPersons.CreateParams{
				Name:    "Johnny",
				Age:     22,
				Address: "Moscow, Red Square",
				Work:    "Yandex",
			},
			Person: models.Person{},
			err:    pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			repo := New(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}
			println("OK 1")

			Person, err := repo.Create(context.TODO(), &test.params)
			println("OK 2")
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if err == nil && *Person != test.Person {
				t.Errorf("\nExpected: %v\nGot: %v", test.Person, Person)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestList(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		Persons []models.Person
		err     error
	}

	const listCmd = `
	SELECT id, name, age, address, work
	FROM persons
	offset $1`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "address", "work"})
				expect := []models.Person{
					{ID: 1,  Name: "Johnny", Age: 22, Address: "Moscow, Red Square", Work: "Yandex"},
					{ID: 2, Name: "Den", Age: 22, Address: "Moscow, Red Square", Work: "Yandex"},
					{ID: 3, Name: "Ken", Age: 22, Address: "Moscow, Red Square", Work: "Yandex"},
				}
				for _, Person := range expect {
					rows = rows.AddRow(Person.ID, Person.Name, Person.Age, Person.Address, Person.Work)
				}
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(0).
					WillReturnRows(rows)
			},
			Persons: []models.Person{
				{ID: 1,  Name: "Johnny", Age: 22, Address: "Moscow, Red Square", Work: "Yandex"},
					{ID: 2, Name: "Den", Age: 22, Address: "Moscow, Red Square", Work: "Yandex"},
					{ID: 3, Name: "Ken", Age: 22, Address: "Moscow, Red Square", Work: "Yandex"},
			},
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(listCmd)).
					WithArgs(0).
					WillReturnError(fmt.Errorf("db error"))
			},
			Persons: nil,
			err:     pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			repo := New(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			Persons, err := repo.List(context.TODO(), 0, 0)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(Persons, test.Persons) {
				t.Errorf("\nExpected: %v\nGot: %v", test.Persons, Persons)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		id      int64
		Person  models.Person
		err     error
	}

	const getCmd = `
	SELECT id, name, age, address, work
	FROM persons
	WHERE id = $1;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "address", "work"})
				rows = rows.AddRow(1, "Johnny", 22, "Moscow, Red Square", "Yandex")
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnRows(rows)
			},
			id:     3,
			Person: models.Person{
				ID:      1,
				Name:    "Johnny",
				Age:     22,
				Address: "Moscow, Red Square",
				Work:    "Yandex",
			},
			err:    nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectQuery(regexp.QuoteMeta(getCmd)).
					WithArgs(3).
					WillReturnError(fmt.Errorf("db error"))
			},
			id:     3,
			Person: models.Person{},
			err:    pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			repo := New(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			Person, err := repo.Get(context.TODO(), test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if err == nil && *Person != test.Person {
				t.Errorf("\nExpected: %v\nGot: %v", test.Person, Person)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type fields struct {
		mock sqlmock.Sqlmock
	}

	type testCase struct {
		prepare func(f *fields)
		id      int64
		err     error
	}

	const deleteCmd = `
	DELETE FROM persons 
	WHERE id = $1;`

	tests := map[string]testCase{
		"good query": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(deleteCmd)).
					WithArgs(3).
					WillReturnResult(sqlmock.NewResult(3, 1))
			},
			id:  3,
			err: nil,
		},
		"query error": {
			prepare: func(f *fields) {
				f.mock.
					ExpectExec(regexp.QuoteMeta(deleteCmd)).
					WithArgs(3).
					WillReturnError(fmt.Errorf("db error"))
			},
			id:  3,
			err: pkgErrors.ErrDb,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("can't create mock: %s", err)
			}
			defer db.Close()

			repo := New(db, logger)

			f := fields{mock: mock}
			if test.prepare != nil {
				test.prepare(&f)
			}

			err = repo.Delete(context.TODO(), test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("\nThere were unfulfilled expectations: %s", err)
			}
		})
	}
}

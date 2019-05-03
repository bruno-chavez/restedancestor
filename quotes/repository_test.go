package quotes

import (
	"errors"
	"testing"

	"github.com/bruno-chavez/restedancestor/database"
)

type DbDouble struct{}

type StmtDouble struct{}

func (d DbDouble) Prepare(sql string, args ...interface{}) (database.Stmt, error) {
	return StmtDouble{}, nil
}

func (s StmtDouble) Close() error {
	return nil
}

var step int

func (s StmtDouble) Step() (bool, error) {
	step++
	return (step <= 2), nil
}

var exec error

func (s StmtDouble) Exec(...interface{}) error {
	return exec
}

func (s StmtDouble) Scan(dst ...interface{}) error {
	return nil
}

var repo = NewRepository(DbDouble{})

func TestRandomOK(t *testing.T) {
	step = 1
	q := repo.Random()
	if q == nil {
		t.Error("Fail for random")
	}
}

func TestRandomKO(t *testing.T) {
	step = 3
	q := repo.Random()
	if q != nil {
		t.Error("Success for random")
	}
}

func TestAllOK(t *testing.T) {
	step = 1
	qs := repo.All()
	if len(qs) != 1 {
		t.Error("No quote")
	}
}

func TestAllKO(t *testing.T) {
	step = 3
	qs := repo.All()
	if len(qs) != 0 {
		t.Error("There's quotes")
	}
}

func TestFindByUUIDOK(t *testing.T) {
	step = 1
	q := repo.FindByUUID("known")
	if q == nil {
		t.Error("No quote with this uuid")
	}
}

func TestFindByUUIDKO(t *testing.T) {
	step = 3
	q := repo.FindByUUID("unknown")
	if q != nil {
		t.Error("There's quote with this uuid")
	}
}

func TestPreferedOK(t *testing.T) {
	step = 5
	qs := repo.Prefered()
	if len(qs) != 5 {
		t.Error("No prefered quote")
	}
}

func TestPreferedKO(t *testing.T) {
	step = 3
	qs := repo.Prefered()
	if len(qs) != 0 {
		t.Error("There's prefered quotes")
	}
}

func TestIncrementsScoreKO(t *testing.T) {
	step = 1
	exec = errors.New("")

	if err := repo.IncrementsScore("unknown"); err == nil {
		t.Error("Score incrementation succeed")
	}
}

func TestIncrementsScoreOK(t *testing.T) {
	step = 1
	exec = nil
	if err := repo.IncrementsScore("known"); err != nil {
		t.Error("Score incrementation failed")
	}
}

func TestDecrementsScoreKO(t *testing.T) {
	step = 1
	exec = errors.New("")

	if err := repo.DecrementsScore("unknown"); err == nil {
		t.Error("Score decrementation succeed")
	}
}

func TestDecrementsScoreOK(t *testing.T) {
	step = 1
	exec = nil
	if err := repo.DecrementsScore("known"); err != nil {
		t.Error("Score decrementation failed")
	}
}

// func TestPostPeopleNoLast(t *testing.T) {
// 	step = 2
// 	p := People{
// 		ID:   97,
// 		Name: "Boba Fett",
// 	}
// 	if id := repo.PostPeople(p); id != 1 {
// 		t.Error("There's people in storage")
// 	}
// }
//
// func TestPeopleByIDOK(t *testing.T) {
// 	step = 1
// 	if _, err := repo.PeopleByID(5); err != nil {
// 		t.Error("There's no people with this id")
// 	}
// }

// func TestPeopleByIDKO(t *testing.T) {
// 	step = 2
// 	if _, err := repo.PeopleByID(7); err == nil {
// 		t.Error("There's people with this id")
// 	}
// }

// func TestDeletePeopleNoPeople(t *testing.T) {
// 	step = 2
// 	if err := repo.DeletePeople(15); err == nil {
// 		t.Error("Found people with this id")
// 	}
// }
//
// func TestDeletePeopleFail(t *testing.T) {
// 	step = 1
// 	exec = errors.New("")
//
// 	if err := repo.DeletePeople(18); err == nil {
// 		t.Error("Fail exec")
// 	}
// }
//
// func TestDeletePeopleOK(t *testing.T) {
// 	step = 1
// 	exec = nil
// 	if err := repo.DeletePeople(21); err != nil {
// 		t.Error("Delete failed")
// 	}
// }

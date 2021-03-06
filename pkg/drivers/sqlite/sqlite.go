package sqlite

import (
	"database/sql"
	"encoding/json"

	"github.com/PacketFire/immigrant/pkg/core"
)

const (
	stateCreate string = `CREATE TABLE imm_sequence_tracker (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  revisionID VARCHAR(256),
  revisionJSON TEXT
);
`
)

// ERRORS

// errCurrentRemoteState is returned when immigrant is unable to fetch the
// remote state's HEAD.
type errCurrentRemoteState struct{}

func (e errCurrentRemoteState) Error() string {
	return "Unable to fetch remote revision state."
}

type errHeadDoesNotExist struct{}

func (e errHeadDoesNotExist) Error() string {
	return "Remote revision HEAD does not exist."
}

// Type Defs

type stateTrackerRevision struct {
	ID           int
	RevisionID   string
	RevisionJSON string
}

// Driver implements github.com/PacketFire/immigrant/pkg/core.Driver, providing
// CRUD methods for operating on a contained sqlite database.
type Driver struct {
	Db        *sql.DB
	Revisions []core.Revision
}

// Init implements the Init method defined on the
// github.com/PacketFire/pkg/core.Driver interface. This method takes a map of
// configuration data and attempts to open a new sql connection to the
// database. If it is unable to do so, an error is returned.
func (dri *Driver) Init(filepath string) error {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}

	dri.Db = db
	return nil
}

// Migrate implements the Migrate method defined on the
// github.com/PacketFire/pkg/core.Driver interface. This method attempts to
// execute all migrations within a passed Revision against a database.
func (dri *Driver) Migrate(r core.Revision) {
	dri.Revisions = append(dri.Revisions, r)
	tx, err := dri.Db.Begin()
	if err != nil {
		return
	}

	for _, mig := range r.Migrate {
		if _, err = tx.Exec(mig); err != nil {
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
}

// Rollback implements the Rollback method defined on the
// github.com/PacketFire/pkg/core.Driver interface. This method attempts to
// execute all rollbacks within a passed Revision against a database.
func (dri *Driver) Rollback(r core.Revision) {
	if len(dri.Revisions) >= 1 {
		dri.Revisions = dri.Revisions[:len(dri.Revisions)-1]
	}

	tx, err := dri.Db.Begin()
	if err != nil {
		return
	}

	for _, mig := range r.Rollback {
		if _, err = tx.Exec(mig); err != nil {
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
}

// State implements the State method defined on the
// github.com/PacketFire/pkg/core.Driver interface. This method attempts to
// return the last Revision added to the databases state.
func (dri *Driver) State() (*core.Revision, error) {
	rHead := new(core.Revision)

	rows, err := dri.Db.Query("SELECT * FROM imm_sequence_tracker ORDER BY id DESC LIMIT 0, 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		row := new(stateTrackerRevision)
		if err = rows.Scan(row); err != nil {
			return nil, errCurrentRemoteState{}
		}

		if err = json.Unmarshal([]byte(row.RevisionJSON), rHead); err != nil {
			return nil, err
		}

		return rHead, nil
	}

	return nil, errHeadDoesNotExist{}
}

func (dri *Driver) initStateManager() error {
	stmt, err := dri.Db.Prepare(stateCreate)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

// Close implements the Close method defined on the
// github.com/PacketFire/pkg/core.Driver interface. This method tears down the
// connection to the database.
func (dri *Driver) Close() {
	dri.Db.Close()
}

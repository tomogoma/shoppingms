package roach

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/cockroachdb/cockroach-go/crdb"
	crdbH "github.com/tomogoma/crdb"
	"github.com/tomogoma/go-typed-errors"
	"github.com/tomogoma/shoppingms/pkg/config"
)

// Roach is a cockroach db store.
// Use NewRoach() to instantiate.
type Roach struct {
	errors.NotFoundErrCheck
	dsn              string
	dbName           string
	db               *sql.DB
	compatibilityErr error

	isDBInitMutex sync.Mutex
	isDBInit      bool
}

const (
	keyDBVersion = "db.version"
)

// NewRoach creates an instance of *Roach. A db connection is only established
// when InitDBIfNot() or one of the Execute/Query methods is called.
func NewRoach(opts ...Option) *Roach {
	r := &Roach{
		isDBInit:      false,
		isDBInitMutex: sync.Mutex{},
		dbName:        config.CanonicalName(),
	}
	for _, f := range opts {
		f(r)
	}
	return r
}

// InitDBIfNot connects to and sets up the DB; creating it and tables if necessary.
func (r *Roach) InitDBIfNot() error {
	var err error
	r.db, err = crdbH.TryConnect(r.dsn, r.db)
	if err != nil {
		return errors.Newf("connect to db: %v", err)
	}
	return r.instantiate()
}

// ExecuteTx prepares a transaction (with retries) for execution in fn.
// It commits the changes if fn returns nil, otherwise changes are rolled back.
func (r *Roach) ExecuteTx(fn func(*sql.Tx) error) error {
	if err := r.InitDBIfNot(); err != nil {
		return err
	}
	return crdb.ExecuteTx(context.Background(), r.db, nil, fn)
}

// ColDesc returns a string containing cols in the given order separated by ",".
func ColDesc(cols ...string) string {
	desc := ""
	for _, col := range cols {
		if col == "" {
			continue
		}
		desc = desc + col + ", "
	}
	return strings.TrimSuffix(desc, ", ")
}

func (r *Roach) instantiate() error {
	r.isDBInitMutex.Lock()
	defer r.isDBInitMutex.Unlock()
	if r.compatibilityErr != nil {
		return r.compatibilityErr
	}
	if r.isDBInit {
		return nil
	}
	if err := crdbH.InstantiateDB(r.db, r.dbName, AllTableDescs...); err != nil {
		return errors.Newf("instantiating db: %v", err)
	}
	if runningVersion, err := r.validateRunningVersion(); err != nil {
		if !r.IsNotFoundError(err) {
			if err != r.compatibilityErr {
				return fmt.Errorf("check db version: %v", err)
			}
			if err := r.migrate(runningVersion, Version); err != nil {
				return fmt.Errorf("migrate from version %d to %d: %v",
					runningVersion, Version, err)
			}
		}
		if err := r.setRunningVersionCurrent(); err != nil {
			return errors.Newf("set db version: %v", err)
		}
	}
	r.isDBInit = true
	return nil
}

func (r *Roach) validateRunningVersion() (int, error) {
	var runningVersion int
	q := `SELECT ` + ColValue + ` FROM ` + TblConfigurations + ` WHERE ` + ColKey + `=$1`
	var confB []byte
	if err := r.db.QueryRow(q, keyDBVersion).Scan(&confB); err != nil {
		if err == sql.ErrNoRows {
			return -1, errors.NewNotFoundf("config not found")
		}
		return -1, errors.Newf("get conf: %v", err)
	}
	if err := json.Unmarshal(confB, &runningVersion); err != nil {
		return -1, errors.Newf("Unmarshalling config: %v", err)
	}
	if runningVersion != Version {
		r.compatibilityErr = errors.Newf("db incompatible: need db"+
			" version '%d', found '%d'", Version, runningVersion)
		return runningVersion, r.compatibilityErr
	}
	return runningVersion, nil
}

func (r *Roach) setRunningVersionCurrent() error {
	valB, err := json.Marshal(Version)
	if err != nil {
		return errors.Newf("marshal conf: %v", err)
	}
	cols := ColDesc(ColKey, ColValue, ColUpdateDate)
	updCols := ColDesc(ColValue, ColUpdateDate)
	q := `
		INSERT INTO ` + TblConfigurations + ` (` + cols + `)
			VALUES ($1, $2, CURRENT_TIMESTAMP)
			ON CONFLICT (` + ColKey + `)
			DO UPDATE SET (` + updCols + `) = ($2, CURRENT_TIMESTAMP)`
	res, err := r.db.Exec(q, keyDBVersion, valB)
	if err := checkRowsAffected(res, err, 1); err != nil {
		return err
	}
	r.compatibilityErr = nil
	return nil
}

func checkRowsAffected(r sql.Result, err error, expAffected int64) error {
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFound("none found")
		}
		return err
	}
	c, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.NewNotFound("none found for update")
	}
	if c != expAffected {
		return errors.Newf("expected %d affected rows but got %d",
			expAffected, c)
	}
	return nil
}

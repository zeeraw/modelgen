package modelgen

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// TxOptions defines an option type for configuring
// transations. This may only be used with the ExecuteTransaction wrapper.
type TxOptions struct {
	Timeout   time.Duration
	Isolation sql.IsolationLevel
	ReadOnly  bool
}

// ExecuteTransaction closes over a transaction and automatically commits
// or rollbacks depending on whether errors were encountered.
// In the case where nil is passed for opt (*TxOption), the following defaults are used:
//  &TxOptions{
//  	Timeout:   5 * time.Second,
//  	Isolation: sql.LevelSerializable,
//  	ReadOnly:  false,
//  }
func ExecuteTransaction(db *sql.DB, opt *TxOptions, actions func(*sql.Tx) error) (err error) {
	// Provide safe defaults in case none were given.
	if opt == nil {
		opt = &TxOptions{
			Timeout:   5 * time.Second,
			Isolation: sql.LevelSerializable,
			ReadOnly:  false,
		}
	}

	// Build the context with the provided timeout.
	// This will be used to define the total time the transaction may take,
	// past this time, it will be cancelled, rollback, then throw an error.
	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
	defer cancel()

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{
		Isolation: opt.Isolation,
		ReadOnly:  opt.ReadOnly,
	}); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			// Only need to log here because panic won't report whether
			// the rollback was successful or not.
			if txerr := tx.Rollback(); txerr != nil {
				log.Println("db rollback error:", txerr)
			}

			log.Printf("rolled back transaction")
			panic(r)
		} else if err != nil {
			// If we run into issues rolling back, keep track of the error that
			// caused the issue and provide some context on the rollback failure.
			if rerr := tx.Rollback(); rerr != nil {
				err = fmt.Errorf("db error: %v rollback error: %v", err, rerr)
			}
		} else {
			if cerr := tx.Commit(); cerr != nil {
				err = fmt.Errorf("commit error: %v", cerr)
			}
		}
	}()

	err = actions(tx)
	return err
}

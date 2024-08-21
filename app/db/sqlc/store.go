package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var storeDB Store
type Store interface {
	Querier
	ExecWithTransaction(ctx context.Context, fn func(*Queries) error) error
	OnNotifyEvent(ctx context.Context, event string, notify chan interface{}) error
	ExecStatementMany(ctx context.Context, statement string, parameters []interface{}) ([]map[string]interface{}, error)
	ExecStatementOne(ctx context.Context, statement string, parameters []interface{}, outputs []string) (map[string]interface{}, error)
	ExecStatement(ctx context.Context, statement string, parameters []interface{}) error
	ExecWithChannelTransaction(ctx context.Context, fns []func(*Queries) error, wg *sync.WaitGroup, result chan error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new Store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries:  New(connPool),
		connPool: connPool,
	}
}

func InitStore(connPool *pgxpool.Pool) Store {
	storeDB = NewStore(connPool)
	return storeDB
}

func (store *SQLStore) ExecWithTransaction(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

func (store *SQLStore) ExecWithChannelTransaction(ctx context.Context, fns []func(*Queries) error, wg *sync.WaitGroup, result chan error) {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		result <- err
	}
	var errorResults = make(chan error, len(fns))
	for _, fn := range fns {
		wg.Add(1)
		go func(fn2 func(*Queries) error) {
			var err error
			wg.Done()
			// tx, err := store.connPool.Begin(ctx)
			// fmt.Println("transaction 2")
			// if err != nil {
			// 	fmt.Println("transaction loi")
			// 	errorResults <- err
			// }
			q := New(tx)
			if err = fn2(q); err != nil {
				tx.Rollback(ctx)
				errorResults <- err
			}
			errorResults <- nil

		}(fn)

	}
	wg.Wait()
	// if err := tx.Commit(ctx); err != nil {
	// 	fmt.Println("transaction loi o tx")
	// 	result <- err
	// 	return
	// }
	if errorResults != nil {
		for err := range errorResults {
			result <- err
		}
		close(result)
	} else {
		if err := tx.Commit(ctx); err != nil {
			result <- err
			return
		}
	}

}

func (store *SQLStore) OnNotifyEvent(ctx context.Context, event string, notify chan interface{}) error {
	tx, err := store.connPool.Acquire(ctx)
	if err != nil {
		return err
	}
	// defer tx.Release()
	_, err = tx.Exec(ctx, fmt.Sprintf("LISTEN %s", event))
	if err != nil {
		return fmt.Errorf("Unable to execute LISTEN command: %v", err)
	}
	go func() {
		for {
			notification, err := tx.Conn().WaitForNotification(ctx)
			if err != nil {
				log.Printf("Unable to wait for a notification: %v\n", err)
				continue
			}

			notify <- notification.Payload
		}
	}()

	return nil
}

func (store *SQLStore) ExecStatementMany(ctx context.Context, statement string, parameters []interface{}) ([]map[string]interface{}, error) {

	rows, err := store.connPool.Query(ctx, statement, parameters...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		for i, val := range values {
			rowMap[string(rows.FieldDescriptions()[i].Name)] = val
		}
		results = append(results, rowMap)
	}
	return results, nil
}

func (store *SQLStore) ExecStatementOne(ctx context.Context, statement string, parameters []interface{}, outputs []string) (map[string]interface{}, error) {

	row := store.connPool.QueryRow(ctx, statement, parameters...)
	// Prepare a slice to hold the values
	values := make([]interface{}, len(outputs))
	valuePtrs := make([]interface{}, len(outputs))

	for i := 0; i < len(outputs); i++ {
		valuePtrs[i] = &values[i]
	}

	err := row.Scan(valuePtrs...)
	if err != nil {
		return nil, err
	}
	rowMap := make(map[string]interface{})
	for i, val := range values {
		rowMap[outputs[i]] = val
	}
	return rowMap, nil
}

func (store *SQLStore) ExecStatement(ctx context.Context, statement string, parameters []interface{}) error {
	_, err := store.connPool.Exec(ctx, statement, parameters...)
	if err != nil {
		return err
	}
	return nil
}

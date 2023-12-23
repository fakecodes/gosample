package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/fakecodes/gosample/domain"
	"github.com/fakecodes/gosample/task/repository"
)

type postgresTaskRepository struct {
	Conn *sql.DB
}

func NewPostgresTaskRepository(conn *sql.DB) domain.TaskRepository {
	return &postgresTaskRepository{conn}
}

func (m *postgresTaskRepository) Create(ctx context.Context, a *domain.Task) (err error) {
	// Define a prepared statement. You'd typically define the statement
	// elsewhere and save it for use in functions such as this one.
	query := `INSERT INTO tasks (name, description, priority, completed) VALUES ($1, $2, $3, $4)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	res, err := stmt.ExecContext(ctx, a.Name, a.Description, a.Priority, a.Completed)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return
}

func (m *postgresTaskRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Task, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Panic(errRow)
		}
	}()

	result = make([]domain.Task, 0)
	for rows.Next() {
		t := domain.Task{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.CreatedAt,
			&t.DueDate,
			&t.Priority,
			&t.Completed,
		)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (m *postgresTaskRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Task, nextCursor string, err error) {
	query := `SELECT id, name, description, created_at, due_date, priority, completed FROM tasks WHERE created_at > $1 ORDER BY created_at LIMIT $2`

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}
	return
}

func (m *postgresTaskRepository) GetByID(ctx context.Context, id int64) (res domain.Task, err error) {
	query := `SELECT id, name, description, created_at, due_date, priority, completed FROM tasks WHERE ID = $1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Task{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *postgresTaskRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM tasks WHERE id = $1"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}
	return
}

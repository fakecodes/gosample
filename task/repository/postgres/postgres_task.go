package postgres

import (
	"context"
	"database/sql"

	"github.com/fakecodes/gosample/domain"
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

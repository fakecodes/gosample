package mysql

import (
	"context"
	"database/sql"

	"github.com/fakecodes/gosample/domain"
)

type mysqlTaskRepository struct {
	Conn *sql.DB
}

// NewMysqlTaskRepository will create an object that represent the task.Repository interface
func NewMysqlTaskRepository(conn *sql.DB) domain.TaskRepository {
	return &mysqlTaskRepository{conn}
}
func (m *mysqlTaskRepository) Create(ctx context.Context, a *domain.Task) (err error) {
	query := `INSERT INTO tasks (name, description, due_date, priority, completed) VALUES (?, ?, ?, ?, ?)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	// converting date to support inserting in mysql
	dueDateTime := a.DueDate.Format("2006-01-02")

	res, err := stmt.ExecContext(ctx, a.Name, a.Description, dueDateTime, a.Priority, a.Completed)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	return
}

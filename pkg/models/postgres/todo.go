package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
	"toDoGolangProject/pkg/models"
)

const (
	getTodoById     = "select id, title, description, created, updated from todo where id=$1"
	insertTodo      = "insert into todo(title, description, created, updated, user_id) values($1,$2,$3,$4, $5) returning id"
	updateTodo      = "update todo set title=$1, description=$2, updated=$3 where id=$4"
	deleteTodo      = "delete todo where id=$1"
	getTodoByUserId = "select id, title, description, created, updated, user_id from todo where user_id=$1"
)

type ToDoModel struct {
	Pool *pgxpool.Pool
}

func (m *ToDoModel) Get(id int) (*models.ToDo, error) {
	t := &models.ToDo{}
	err := m.Pool.QueryRow(context.Background(), getTodoById, id).
		Scan(&t.ID, &t.Title, &t.Description, &t.Created, &t.Updated)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return t, nil
}

func (m *ToDoModel) Insert(title, description string) (int, error) {
	var id uint64
	row := m.Pool.QueryRow(context.Background(), insertTodo, title, description, time.Now(), time.Now(), UserId)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *ToDoModel) Update(title, description string, id int) error {
	_, err := m.Pool.Query(context.Background(), updateTodo, title, description, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (m *ToDoModel) Delete(id int) error {
	_, err := m.Pool.Query(context.Background(), deleteTodo, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *ToDoModel) AllTodo() ([]*models.ToDo, error) {
	var todoes []*models.ToDo
	fmt.Println(11233)


	rows, err := m.Pool.Query(context.Background(), getTodoByUserId, UserId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t *models.ToDo
		err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.Created, &t.Updated)
		if err != nil {
			return nil, err
		}
		todoes = append(todoes, t)
	}
	return todoes, nil
}

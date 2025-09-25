package repository

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/KostyShatovGO/todo-app"
	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title,description) VALUES ($1,$2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id,list_id) VALUES ($1,$2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)
	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList
	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl
								INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)
	return list, err
}
func (r *TodoListPostgres) Delete(userId, listId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Ensure the list belongs to the user
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE user_id=$1 AND list_id=$2)", usersListsTable)
	if err := tx.QueryRow(checkQuery, userId, listId).Scan(&exists); err != nil {
		tx.Rollback()
		return err
	}
	if !exists {
		tx.Rollback()
		return sql.ErrNoRows
	}

	// Delete items associated with the list
	deleteItemsQuery := fmt.Sprintf("DELETE FROM %s ti USING %s li WHERE ti.id = li.item_id AND li.list_id = $1", todoItemsTable, listsItemsTable)
	if _, err := tx.Exec(deleteItemsQuery, listId); err != nil {
		tx.Rollback()
		return err
	}

	// Delete the list (will cascade remove links in users_lists and lists_items)
	deleteListQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", todoListsTable)
	res, err := tx.Exec(deleteListQuery, listId)
	if err != nil {
		tx.Rollback()
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		tx.Rollback()
		return sql.ErrNoRows
	}

	return tx.Commit()
}
func (r *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ",")
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)

	logrus.Debugf("updateQuery:%s", query)
	logrus.Debugf("args:%s", args)

	_, err := r.db.Exec(query, args...)
	return err
}

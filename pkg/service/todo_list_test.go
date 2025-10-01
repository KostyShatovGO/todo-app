package service

import (
    "errors"
    "testing"

    todo "github.com/KostyShatovGO/todo-app"
)

type mockListRepo struct{
    create func(userId int, list todo.TodoList) (int, error)
    getAll func(userId int) ([]todo.TodoList, error)
    getById func(userId, listId int) (todo.TodoList, error)
    delete func(userId, listId int) error
    update func(userId, listId int, input todo.UpdateListInput) error
}

func (m *mockListRepo) Create(userId int, list todo.TodoList) (int, error) {
    if m.create!=nil { return m.create(userId, list) }
    return 0, errors.New("no impl")
}
func (m *mockListRepo) GetAll(userId int) ([]todo.TodoList, error) {
    if m.getAll!=nil { return m.getAll(userId) }
    return nil, nil
}
func (m *mockListRepo) GetById(userId, listId int) (todo.TodoList, error) {
    if m.getById!=nil { return m.getById(userId, listId) }
    return todo.TodoList{}, nil
}
func (m *mockListRepo) Delete(userId, listId int) error {
    if m.delete!=nil { return m.delete(userId, listId) }
    return nil
}
func (m *mockListRepo) Update(userId, listId int, input todo.UpdateListInput) error {
    if m.update!=nil { return m.update(userId, listId, input) }
    return nil
}

func TestTodoListService_Update_Validate(t *testing.T) {
    s := NewTodoListService(&mockListRepo{})
    err := s.Update(1, 2, todo.UpdateListInput{})
    if err == nil {
        t.Fatal("expected validation error")
    }
}

func TestTodoListService_Create_OK(t *testing.T) {
    called := false
    s := NewTodoListService(&mockListRepo{create: func(userId int, list todo.TodoList) (int, error) {
        called = true
        if userId!=7 || list.Title!="A" { t.Fatalf("bad args: %v %v", userId, list) }
        return 10, nil
    }})
    id, err := s.Create(7, todo.TodoList{Title: "A"})
    if err != nil { t.Fatalf("unexpected: %v", err) }
    if !called || id!=10 { t.Fatalf("not called or bad id: %v %d", called, id) }
}



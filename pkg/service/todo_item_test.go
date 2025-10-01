package service

import (
    "errors"
    "testing"

    todo "github.com/KostyShatovGO/todo-app"
)

type mockItemRepo struct{
    create func(listId int, item todo.TodoItem) (int, error)
    getAll func(userId, listId int) ([]todo.TodoItem, error)
    getById func(userId, itemId int) (todo.TodoItem, error)
    delete func(userId, itemId int) error
    update func(userId, itemId int, input todo.UpdateItemInput) error
}
func (m *mockItemRepo) Create(listId int, item todo.TodoItem) (int, error) { if m.create!=nil { return m.create(listId, item) }; return 0, nil }
func (m *mockItemRepo) GetAll(userId, listId int) ([]todo.TodoItem, error) { if m.getAll!=nil { return m.getAll(userId, listId) }; return nil, nil }
func (m *mockItemRepo) GetById(userId, itemId int) (todo.TodoItem, error) { if m.getById!=nil { return m.getById(userId, itemId) }; return todo.TodoItem{}, nil }
func (m *mockItemRepo) Delete(userId, itemId int) error { if m.delete!=nil { return m.delete(userId, itemId) }; return nil }
func (m *mockItemRepo) Update(userId, itemId int, input todo.UpdateItemInput) error { if m.update!=nil { return m.update(userId, itemId, input) }; return nil }

func TestTodoItemService_Create_CheckListExists(t *testing.T) {
    listCalled := false
    listRepo := &mockListRepo{getById: func(userId, listId int) (todo.TodoList, error) {
        listCalled = true
        if listId != 5 || userId != 3 { t.Fatalf("bad args: %d %d", userId, listId) }
        return todo.TodoList{Id: listId, Title: "L"}, nil
    }}
    itemCalled := false
    itemRepo := &mockItemRepo{create: func(listId int, item todo.TodoItem) (int, error) {
        itemCalled = true
        if listId != 5 { t.Fatalf("bad listId: %d", listId) }
        return 11, nil
    }}
    s := NewTodoItemService(itemRepo, listRepo)
    id, err := s.Create(3, 5, todo.TodoItem{Title: "I"})
    if err != nil { t.Fatalf("unexpected: %v", err) }
    if !listCalled || !itemCalled || id != 11 { t.Fatalf("calls: %v %v id=%d", listCalled, itemCalled, id) }
}

func TestTodoItemService_Create_ListMissing(t *testing.T) {
    listRepo := &mockListRepo{getById: func(userId, listId int) (todo.TodoList, error) {
        return todo.TodoList{}, errors.New("not found")
    }}
    s := NewTodoItemService(&mockItemRepo{}, listRepo)
    _, err := s.Create(1, 2, todo.TodoItem{Title: "I"})
    if err == nil { t.Fatal("expected error when list missing") }
}



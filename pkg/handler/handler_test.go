package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    todo "github.com/KostyShatovGO/todo-app"
    "github.com/KostyShatovGO/todo-app/pkg/service"
    "github.com/gin-gonic/gin"
)

type stubAuth struct{
    createUser func(user todo.User) (int, error)
    genToken func(username, password string) (string, error)
    parse func(token string) (int, error)
}
func (s *stubAuth) CreateUser(user todo.User) (int, error) { return s.createUser(user) }
func (s *stubAuth) GenerateToken(username, password string) (string, error) { return s.genToken(username, password) }
func (s *stubAuth) ParseToken(token string) (int, error) { return s.parse(token) }

type stubLists struct{}
func (s *stubLists) Create(userId int, list todo.TodoList) (int, error) { return 1, nil }
func (s *stubLists) GetAll(userId int) ([]todo.TodoList, error) { return []todo.TodoList{}, nil }
func (s *stubLists) GetById(userId, listId int) (todo.TodoList, error) { return todo.TodoList{Id:listId}, nil }
func (s *stubLists) Delete(userId, listId int) error { return nil }
func (s *stubLists) Update(userId, listId int, input todo.UpdateListInput) error { return nil }

type stubItems struct{}
func (s *stubItems) Create(userId, listId int, item todo.TodoItem) (int, error) { return 1, nil }
func (s *stubItems) GetAll(userId, listId int) ([]todo.TodoItem, error) { return []todo.TodoItem{}, nil }
func (s *stubItems) GetById(userId, itemId int) (todo.TodoItem, error) { return todo.TodoItem{Id:itemId}, nil }
func (s *stubItems) Delete(userId, itemId int) error { return nil }
func (s *stubItems) Update(userId, itemId int, input todo.UpdateItemInput) error { return nil }

func TestAuth_SignUp_OK(t *testing.T) {
    gin.SetMode(gin.TestMode)
    svc := &service.Service{ Authorization: &stubAuth{ createUser: func(user todo.User) (int, error) { return 7, nil }, genToken: func(string,string)(string,error){return "", nil}, parse: func(string)(int,error){return 0,nil} }, TodoList: &stubLists{}, TodoItem: &stubItems{} }
    h := NewHandler(svc)
    r := h.InitRoutes()
    body, _ := json.Marshal(map[string]string{"name":"N","username":"u","password":"p"})
    req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
    }
}

func TestAuth_SignIn_OK(t *testing.T) {
    gin.SetMode(gin.TestMode)
    svc := &service.Service{ Authorization: &stubAuth{ genToken: func(username, password string) (string, error) { return "tok", nil }, createUser: func(user todo.User) (int, error) { return 1, nil }, parse: func(string)(int,error){return 0,nil} }, TodoList: &stubLists{}, TodoItem: &stubItems{} }
    h := NewHandler(svc)
    r := h.InitRoutes()
    body, _ := json.Marshal(map[string]string{"username":"u","password":"p"})
    req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
    }
}



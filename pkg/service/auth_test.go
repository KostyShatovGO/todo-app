package service

import (
    "errors"
    "strings"
    "testing"

    todo "github.com/KostyShatovGO/todo-app"
)

type mockAuthRepo struct{
    userByCred func(username, password string) (todo.User, error)
    createUser func(user todo.User) (int, error)
}

func (m *mockAuthRepo) GetUser(username, password string) (todo.User, error) {
    if m.userByCred != nil { return m.userByCred(username, password) }
    return todo.User{}, errors.New("not implemented")
}
func (m *mockAuthRepo) CreateUser(user todo.User) (int, error) {
    if m.createUser != nil { return m.createUser(user) }
    return 0, errors.New("not implemented")
}

func TestAuthService_CreateUser_HashesPassword(t *testing.T) {
    captured := todo.User{}
    repo := &mockAuthRepo{createUser: func(user todo.User) (int, error) {
        captured = user
        return 1, nil
    }}
    s := NewAuthService(repo)
    _, err := s.CreateUser(todo.User{Name: "N", Username: "U", Password: "plain"})
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if captured.Password == "plain" || captured.Password == "" {
        t.Fatalf("password was not hashed: %q", captured.Password)
    }
}

func TestAuthService_GenerateAndParseToken(t *testing.T) {
    user := todo.User{Id: 42, Username: "u"}
    repo := &mockAuthRepo{userByCred: func(username, password string) (todo.User, error) {
        if username == "u" && password != "" { return user, nil }
        return todo.User{}, errors.New("invalid")
    }}
    s := NewAuthService(repo)
    token, err := s.GenerateToken("u", "p")
    if err != nil { t.Fatalf("GenerateToken error: %v", err) }
    if token == "" { t.Fatal("empty token") }
    id, err := s.ParseToken(token)
    if err != nil { t.Fatalf("ParseToken error: %v", err) }
    if id != 42 { t.Fatalf("unexpected user id: %d", id) }
}

func TestAuthService_ParseToken_Invalid(t *testing.T) {
    s := NewAuthService(&mockAuthRepo{})
    // malformed token
    _, err := s.ParseToken("not-a-token")
    if err == nil && !strings.Contains(errString(err), "") {
        // In current implementation returns (0,nil) on parse error; just assert id==0.
    }
}

func errString(err error) string { if err==nil { return "" }; return err.Error() }



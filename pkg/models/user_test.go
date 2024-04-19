package models

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	errs := User{}.Validate()

	if errs == nil {
		t.Fatalf(`User.Validate() == %q want errors`, errs)
	}

	tests := []struct {
		user User
		attr string
		want string
	}{
		{
			User{Name: "", Email: "text@example.com", Username: "test"},
			"Name",
			"name is empty",
		},
		{
			User{Name: "test", Email: "", Username: "test"},
			"Email",
			"email is empty",
		},
		{
			User{Name: "test", Email: "text@example.com", Username: ""},
			"Username",
			"username is empty",
		},
	}

	for _, tt := range tests {
		errs := tt.user.Validate()
		found := false
		for _, err := range errs {
			found = strings.Contains(err.Error(), tt.want)
		}
		if !found {
			t.Errorf(`User{%s: ""}.Validate() == %v want nil`, tt.attr, errs)
		}
	}

}

func TestHasRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		user User
		role string
		want bool
	}{
		{
			User{Roles: []string{"user"}},
			"user",
			true,
		},
		{
			User{Roles: []string{"user"}},
			"admin",
			false,
		},
		{
			User{Roles: []string{"user", "admin"}},
			"user",
			true,
		},
		{
			User{Roles: []string{"user", "admin"}},
			"none",
			false,
		},
	}

	for _, tt := range tests {
		got := tt.user.HasRole(tt.role)
		if got != tt.want {
			t.Errorf(`User.HasRole(%v) == %v want true`, tt.role, got)
		}
	}
}

func TestHasRoleAnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		user  User
		roles []string
		want  bool
	}{
		{
			User{Roles: []string{}},
			[]string{},
			true,
		},
		{
			User{Roles: []string{"user"}},
			[]string{},
			true,
		},
		{
			User{Roles: []string{}},
			[]string{"user"},
			false,
		},
		{
			User{Roles: []string{"user", "admin"}},
			[]string{"user", "admin"},
			true,
		},
		{
			User{Roles: []string{"user", "admin", "none"}},
			[]string{"user", "admin"},
			true,
		},
		{
			User{Roles: []string{"user", "admin"}},
			[]string{"user", "admin", "none"},
			false,
		},
		{
			User{Roles: []string{"user", "admin"}},
			[]string{"user", "none"},
			false,
		},
	}

	for _, tt := range tests {
		got := tt.user.HasRolesAnd(tt.roles)
		if got != tt.want {
			t.Errorf(`User.HasRolesAnd(%v) == %v want true`, tt.roles, got)
		}
	}
}

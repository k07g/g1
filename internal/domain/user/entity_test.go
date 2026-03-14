package user_test

import (
	"testing"

	domain "github.com/k07g/g1/internal/domain/user"
)

func TestUser_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		user    domain.User
		wantErr error
	}{
		{
			name:    "正常: 名前あり・年齢0",
			user:    domain.User{Name: "テストユーザー", Age: 0},
			wantErr: nil,
		},
		{
			name:    "正常: 名前あり・年齢あり",
			user:    domain.User{Name: "テストユーザー", Age: 30},
			wantErr: nil,
		},
		{
			name:    "異常: 名前が空",
			user:    domain.User{Name: "", Age: 20},
			wantErr: domain.ErrInvalidName,
		},
		{
			name:    "異常: 年齢がマイナス",
			user:    domain.User{Name: "テストユーザー", Age: -1},
			wantErr: domain.ErrInvalidAge,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.user.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

package item_test

import (
	"testing"

	domain "github.com/k07g/g1/internal/domain/item"
)

func TestItem_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		item    domain.Item
		wantErr error
	}{
		{
			name:    "正常: 名前あり・価格0",
			item:    domain.Item{Name: "テスト商品", Price: 0},
			wantErr: nil,
		},
		{
			name:    "正常: 名前あり・価格あり",
			item:    domain.Item{Name: "テスト商品", Price: 1000},
			wantErr: nil,
		},
		{
			name:    "異常: 名前が空",
			item:    domain.Item{Name: "", Price: 100},
			wantErr: domain.ErrInvalidName,
		},
		{
			name:    "異常: 価格がマイナス",
			item:    domain.Item{Name: "テスト商品", Price: -1},
			wantErr: domain.ErrInvalidPrice,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.item.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

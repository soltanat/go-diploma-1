package external

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/soltanat/go-diploma-1/internal/clients/accrual"
	"github.com/soltanat/go-diploma-1/internal/entities"
)

func TestAccrualStorage_Get(t *testing.T) {
	type fields struct {
		client accrual.ClientWithResponsesInterface
	}
	type args struct {
		ctx    context.Context
		number entities.OrderNumber
	}

	client, err := accrual.NewClientWithResponses("http://localhost:8080")
	assert.NoError(t, err)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entities.AccrualOrder
		wantErr bool
	}{
		{
			name: "Get accrual order",
			fields: fields{
				client: client,
			},
			args: args{
				ctx:    context.Background(),
				number: 4561261212345467,
			},
			want: &entities.AccrualOrder{
				Number:  1,
				Status:  entities.AccrualOrderStatusREGISTERED,
				Accrual: &entities.Currency{Whole: 0},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AccrualStorage{
				client: tt.fields.client,
			}
			got, err := s.Get(tt.args.ctx, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

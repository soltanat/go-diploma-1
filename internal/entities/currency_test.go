package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrency_Add(t *testing.T) {
	type fields struct {
		Whole   int
		Decimal int
	}
	type args struct {
		add *Currency
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Currency
	}{
		{
			name: "add",
			fields: fields{
				Whole:   1,
				Decimal: 90,
			},
			args: args{
				add: &Currency{
					Whole:   1,
					Decimal: 90,
				},
			},
			want: &Currency{
				Whole:   3,
				Decimal: 80,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Currency{
				Whole:   tt.fields.Whole,
				Decimal: tt.fields.Decimal,
			}
			c.Add(tt.args.add)
			if got := c; !assert.Equal(t, tt.want, got) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrency_Sub(t *testing.T) {
	type fields struct {
		Whole   int
		Decimal int
	}
	type args struct {
		sub *Currency
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Currency
		wantErr bool
	}{
		{
			name: "zero result",
			fields: fields{
				Whole:   1,
				Decimal: 90,
			},
			args: args{
				sub: &Currency{
					Whole:   1,
					Decimal: 90,
				},
			},
			want: &Currency{
				Whole:   0,
				Decimal: 0,
			},
			wantErr: false,
		},
		{
			name: "sub out of balance",
			fields: fields{
				Whole:   1,
				Decimal: 90,
			},
			args: args{
				sub: &Currency{
					Whole:   2,
					Decimal: 0,
				},
			},
			want: &Currency{
				Whole:   1,
				Decimal: 90,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Currency{
				Whole:   tt.fields.Whole,
				Decimal: tt.fields.Decimal,
			}
			err := c.Sub(tt.args.sub)
			if tt.wantErr {
				assert.ErrorAs(t, err, &OutOfBalanceError{})
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, c)
		})
	}
}

func TestCurrency_Validate(t *testing.T) {
	type fields struct {
		Whole   int
		Decimal int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Whole:   1,
				Decimal: 90,
			},
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				Whole:   -1,
				Decimal: 90,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Currency{
				Whole:   tt.fields.Whole,
				Decimal: tt.fields.Decimal,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

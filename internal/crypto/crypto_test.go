package crypto

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				password: `123sdfa`,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if IsPasswordCorrect(got, tt.args.password) != tt.want {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPasswordCorrect(t *testing.T) {
	type args struct {
		hashedPassword string
		givenPassword  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Ok",
			args: args{
				hashedPassword: "$2a$04$ghmxNXhvUffxxe72TBN9WO1bX8zeIheCHze.ssj/yuVXk0p9Nuf/O",
				givenPassword: `123sdfa`,
			},
			want:    true,
		},
		{
			name: "False",
			args: args{
				hashedPassword: "$2a$04$ghmxNXhvUffxxe72TBN9WO1bX8zeIheCHze.ssj/yuVXk0p9Nuf/O",
				givenPassword: `asdawda`,
			},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPasswordCorrect(tt.args.hashedPassword, tt.args.givenPassword); got != tt.want {
				t.Errorf("IsPasswordCorrect() = %v, want %v", got, tt.want)
			}
		})
	}
}

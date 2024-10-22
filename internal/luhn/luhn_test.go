package luhn

import "testing"

func TestLuhnAlgorithm(t *testing.T) {
	type args struct {
		cardNumber string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "false",
			args: args{
				cardNumber: `2416231740242`,
			},
			want: false,
		},
		{
			name: "true",
			args: args{
				cardNumber: `2416231740`,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LuhnAlgorithm(tt.args.cardNumber); got != tt.want {
				t.Errorf("LuhnAlgorithm() = %v, want %v", got, tt.want)
			}
		})
	}
}

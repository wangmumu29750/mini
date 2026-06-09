package mock

import "testing"

func TestVerifyIdentity(t *testing.T) {
	type args struct {
		realName   string
		idCardNo   string
		phone      string
		bankCardNo string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyIdentity(tt.args.realName, tt.args.idCardNo, tt.args.phone, tt.args.bankCardNo); (err != nil) != tt.wantErr {
				t.Errorf("VerifyIdentity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
		{name: "all valid", args: args{realName: "张三", idCardNo: "110101199001011234", phone: "13800138000", bankCardNo: "6222020202020202020"}, wantErr: false},
		{name: "empty name", args: args{realName: "", idCardNo: "110101199001011234", phone: "13800138000", bankCardNo: "6222020202020202020"}, wantErr: true},
		{name: "whitespace name only", args: args{realName: "   ", idCardNo: "110101199001011234", phone: "13800138000", bankCardNo: "6222020202020202020"}, wantErr: true},
		{name: "invalid id card - too short", args: args{realName: "张三", idCardNo: "123", phone: "13800138000", bankCardNo: "6222020202020202020"}, wantErr: true},
		{name: "invalid id card - 18 digits with X", args: args{realName: "张三", idCardNo: "11010119900101123X", phone: "13800138000", bankCardNo: "6222020202020202020"}, wantErr: false},
		{name: "invalid phone - letters", args: args{realName: "张三", idCardNo: "110101199001011234", phone: "1380013800a", bankCardNo: "6222020202020202020"}, wantErr: true},
		{name: "invalid phone - too short", args: args{realName: "张三", idCardNo: "110101199001011234", phone: "13800", bankCardNo: "6222020202020202020"}, wantErr: true},
		{name: "invalid bank card - too short", args: args{realName: "张三", idCardNo: "110101199001011234", phone: "13800138000", bankCardNo: "6222020"}, wantErr: true},
		{name: "invalid bank card - too long", args: args{realName: "张三", idCardNo: "110101199001011234", phone: "13800138000", bankCardNo: "6222020202020202020123456"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyIdentity(tt.args.realName, tt.args.idCardNo, tt.args.phone, tt.args.bankCardNo); (err != nil) != tt.wantErr {
				t.Errorf("VerifyIdentity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

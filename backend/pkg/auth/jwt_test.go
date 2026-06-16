package auth

import (
	"reflect"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	type args struct {
		secret    string
		ttl       time.Duration
		principal Principal
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid admin token", args: args{secret: "test-secret", ttl: time.Hour, principal: Principal{UserID: 1, Username: "admin", Role: "admin"}}, wantErr: false},
		{name: "valid passenger token", args: args{secret: "another-secret", ttl: 30 * time.Minute, principal: Principal{UserID: 2, Username: "alice", Role: "passenger"}}, wantErr: false},
		{name: "empty secret", args: args{secret: "", ttl: time.Hour, principal: Principal{UserID: 1, Username: "admin", Role: "admin"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.args.secret, tt.args.ttl, tt.args.principal)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("GenerateToken() returned empty token")
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	// Pre-generate valid tokens for testing
	validToken, _ := GenerateToken("test-secret", time.Hour, Principal{UserID: 1, Username: "admin", Role: "admin"})
	expiredToken, _ := GenerateToken("test-secret", -time.Hour, Principal{UserID: 2, Username: "alice", Role: "passenger"})
	wrongSecretToken, _ := GenerateToken("other-secret", time.Hour, Principal{UserID: 1, Username: "admin", Role: "admin"})

	type args struct {
		secret string
		raw    string
	}
	tests := []struct {
		name    string
		args    args
		want    Principal
		wantErr bool
	}{
		{name: "valid token", args: args{secret: "test-secret", raw: validToken}, want: Principal{UserID: 1, Username: "admin", Role: "admin"}, wantErr: false},
		{name: "expired token", args: args{secret: "test-secret", raw: expiredToken}, want: Principal{}, wantErr: true},
		{name: "wrong secret", args: args{secret: "test-secret", raw: wrongSecretToken}, want: Principal{}, wantErr: true},
		{name: "empty token", args: args{secret: "test-secret", raw: ""}, want: Principal{}, wantErr: true},
		{name: "malformed token", args: args{secret: "test-secret", raw: "not.a.jwt"}, want: Principal{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseToken(tt.args.secret, tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

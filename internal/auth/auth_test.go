package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name 		string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr 	bool
	}{
		{
			name: 		 "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr: 	 false,
		},
		{
			name: 		 "invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "wrong secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("VAlidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
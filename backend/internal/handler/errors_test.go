package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_respondError(t *testing.T) {
	type args struct {
		c   *gin.Context
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respondError(tt.args.c, tt.args.err)
		})
	}
}

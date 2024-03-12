package logging

import (
	"context"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_New_DefaultLogLevel(t *testing.T) {
	l := New(context.Background())
	l.SystemError()
	assert.Equal(t, log.InfoLevel, l.Logger.GetLevel())
}

func TestSetLoggingFormat(t *testing.T) {
	SetLoggingFormat(log.ErrorLevel)

	if log.GetLevel() != log.ErrorLevel {
		t.Errorf("SetLoggingFormat() wanted %v, got %v", log.ErrorLevel, log.GetLevel())
	}

	SetLoggingFormat(log.PanicLevel)

	if log.GetLevel() != log.PanicLevel {
		t.Errorf("SetLoggingFormat() wanted %v, got %v", log.PanicLevel, log.GetLevel())
	}
}

func TestFieldsFromCTX(t *testing.T) {
	ctx := context.Background()

	ctx = context.WithValue(ctx, ContextValues, CtxValues{
		RequestID:     "requestId",
		TransactionID: "transactionId",
		CommitInfo:    "commitInfo",
	})

	fields := make(map[string]interface{})
	fields["requestId"] = "requestId"
	fields["transactionId"] = "transactionId"
	fields["commitInfo"] = "commitInfo"

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name       string
		args       args
		wantFields log.Fields
	}{
		{
			name: "gets all values from context",
			args: args{
				ctx: ctx,
			},
			wantFields: fields,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFields := FieldsFromCTX(tt.args.ctx); !reflect.DeepEqual(gotFields, tt.wantFields) {
				t.Errorf("FieldsFromCTX() = %v, want %v", gotFields, tt.wantFields)
			}
		})
	}
}

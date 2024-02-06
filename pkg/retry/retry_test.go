package retry

import (
	"testing"
	"time"
)

func TestRetry_Retry(t *testing.T) {
	type fields struct {
		retryCount   int
		retryDelay   time.Duration
		retryTimeout time.Duration
	}
	type args struct {
		fn func() error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test3",
			fields: fields{retryCount: 3, retryDelay: 1, retryTimeout: 1000},
			args: args{fn: func() error {

				// Do some work here
				//return RetriableError{error: errors.New("foo")}
				return nil
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Retry{
				retryCount: tt.fields.retryCount,
				retryDelay: tt.fields.retryDelay,
			}
			if err := r.Retry(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Retry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package test

import (
	"reflect"
	"testing"
	"gosandbox/proxy"
)

func TestPolicies(t *testing.T) {
	type args struct {
		keys []string
		vals []string
	}
	tests := []struct {
		name    string
		args    args
		want    []proxy.Policy
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proxy.Policies(tt.args.keys, tt.args.vals)
			if (err != nil) != tt.wantErr {
				t.Errorf("Policies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Policies() = %v, want %v", got, tt.want)
			}
		})
	}
}
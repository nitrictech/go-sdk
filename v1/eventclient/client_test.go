package eventclient

import (
	"reflect"
	"testing"
)

func TestNitricEventClient_GetTopics(t *testing.T) {
	tests := []struct {
		name    string
		want    []Topic
		wantErr bool
	}{
		{
			name: "test something",
			want: []Topic{
				&NitricTopic{name: "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NitricEventClient{}
			got, err := e.GetTopics()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTopics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTopics() got = %v, want %v", got, tt.want)
			}
		})
	}
}

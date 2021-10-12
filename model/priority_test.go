package model

import (
	"reflect"
	"testing"
)

func TestPriority_Ptr(t *testing.T) {
	critical := CriticalPriority
	tests := []struct {
		name string
		i    Priority
		want *Priority
	}{
		{name: "converts to pointer", i: critical, want: &critical},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Ptr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriority_String(t *testing.T) {
	tests := []struct {
		name string
		i    Priority
		want string
	}{
		{name: "CriticalPriority.String() yields CriticalPriority", i: CriticalPriority, want: "CriticalPriority"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriority_unmarshal(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		i       Priority
		args    args
		wantErr bool
	}{
		{name: "low", i: LowPriority, args: args{bytes: []byte("low")}, wantErr: false},
		{name: "LowPriority", i: LowPriority, args: args{bytes: []byte("LowPriority")}, wantErr: false},
		{name: "medium", i: MediumPriority, args: args{bytes: []byte("medium")}, wantErr: false},
		{name: "MediumPriority", i: MediumPriority, args: args{bytes: []byte("MediumPriority")}, wantErr: false},
		{name: "high", i: HighPriority, args: args{bytes: []byte("high")}, wantErr: false},
		{name: "HighPriority", i: HighPriority, args: args{bytes: []byte("HighPriority")}, wantErr: false},
		{name: "critical", i: CriticalPriority, args: args{bytes: []byte("critical")}, wantErr: false},
		{name: "CriticalPriority", i: CriticalPriority, args: args{bytes: []byte("CriticalPriority")}, wantErr: false},
		{name: "unknown", i: Priority(0), args: args{bytes: []byte("unknown")}, wantErr: true},

		// special casing, <3 chars => LowPriority
		{name: "lo", i: LowPriority, args: args{bytes: []byte("lo")}, wantErr: false},
		{name: "na", i: LowPriority, args: args{bytes: []byte("na")}, wantErr: false},
		{name: "empty string", i: LowPriority, args: args{bytes: []byte("")}, wantErr: false},

		// different casings
		{name: "mediumpriority", i: MediumPriority, args: args{bytes: []byte("mediumpriority")}, wantErr: false},
		{name: "medium-priority", i: MediumPriority, args: args{bytes: []byte("medium-priority")}, wantErr: false},
		{name: "medium priority", i: MediumPriority, args: args{bytes: []byte("medium priority")}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.i.unmarshal(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package main

import (
	"reflect"
	"testing"
	"davidhancock.com/varchive"
)

func TestScanPaths(t *testing.T) {
	tests := []struct {
		name string
		want map[string]varchive.FilesWithSize
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := varchive.ScanPaths(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

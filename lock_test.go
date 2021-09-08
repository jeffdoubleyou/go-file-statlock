package go_file_statlock

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestLock_Lock(t *testing.T) {
	type fields struct {
		Path     string
		Duration int
		Status   int
		file     *os.File
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Lock{
				Path:     tt.fields.Path,
				Duration: tt.fields.Duration,
				Status:   tt.fields.Status,
				file:     tt.fields.file,
			}
			got, err := l.Lock()
			if (err != nil) != tt.wantErr {
				t.Errorf("Lock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Lock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLock_Unlock(t *testing.T) {
	type fields struct {
		Path     string
		Duration int
		Status   int
		file     *os.File
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Lock{
				Path:     tt.fields.Path,
				Duration: tt.fields.Duration,
				Status:   tt.fields.Status,
				file:     tt.fields.file,
			}
			got, err := l.Unlock()
			if (err != nil) != tt.wantErr {
				t.Errorf("Unlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Unlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func sleepLock(t *testing.T, lock *Lock, seconds int, w *sync.WaitGroup) {
	t.Logf("Sleep for %d seconds", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)
	if _, err := lock.Unlock(); err != nil {
		fmt.Printf("Got error unlocking file: %s\n", err.Error())
		t.Error(err.Error())
	} else {
		t.Logf("Unlocked file: %d", lock.Status)
	}
	w.Done()
}

func TestLock_WaitForLock(t *testing.T) {
	type args struct {
		path    string
		seconds int
		sleep   int
	}
	tests := []struct {
		name string
		args args
		want *Lock
	}{
		{
			name: "InitialLock",
			args: args{
				path:    "testlock.pid",
				seconds: 5,
				sleep:   10,
			},
			want: &Lock{
				Path:            "testlock.pid",
				Duration:        5,
				Status:          0,
				MaxWaitInterval: 10,
			},
		},
		{
			name: "WaitForInitialLock",
			args: args{
				path:    "testlock.pid",
				seconds: 5,
			},
			want: &Lock{
				Path:            "testlock.pid",
				Duration:        5,
				Status:          0,
				MaxWaitInterval: 10,
			},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLock(tt.args.path, tt.args.seconds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLock() = %v, want %v", got, tt.want)
			} else {
				if _, err := got.WaitForLock(); err != nil {
					t.Errorf("Unable to lock file: %s", err.Error())
				} else {
					t.Logf("Lock aquired with status %d", got.Status)
					if tt.args.sleep > 0 {
						t.Logf("Run sleep in background")
						go sleepLock(t, got, tt.args.sleep, &wg)
					} else {
						if _, err := got.Unlock(); err != nil {
							t.Error(err.Error())
						} else {
							t.Logf("Unlocked file: %d", got.Status)
						}
						wg.Done()
						wg.Wait()
					}
				}
			}
		})
	}
	wg.Wait()
}
func TestLock_openLock(t *testing.T) {
	type fields struct {
		Path     string
		Duration int
		Status   int
		file     *os.File
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Lock{
				Path:     tt.fields.Path,
				Duration: tt.fields.Duration,
				Status:   tt.fields.Status,
				file:     tt.fields.file,
			}
			got, err := l.openLock()
			if (err != nil) != tt.wantErr {
				t.Errorf("openLock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("openLock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLock(t *testing.T) {
	type args struct {
		path    string
		seconds int
	}
	tests := []struct {
		name string
		args args
		want *Lock
	}{
		{
			name: "NewLock",
			args: args{
				path:    "testlock.pid",
				seconds: 5,
			},
			want: &Lock{
				Path:     "testlock.pid",
				Duration: 5,
				Status:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLock(tt.args.path, tt.args.seconds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLock() = %v, want %v", got, tt.want)
			} else {
				if _, err := got.Lock(); err != nil {
					t.Errorf("Unable to lock file: %s", err.Error())
				} else {
					t.Logf("Lock aquired with status %d", got.Status)
					if _, err := got.Unlock(); err != nil {
						t.Error(err.Error())
					} else {
						t.Logf("Unlocked file: %d", got.Status)
					}
				}
			}
		})
	}
}

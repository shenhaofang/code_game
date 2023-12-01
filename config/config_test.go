package config

import "testing"

func TestConfig_Load(t *testing.T) {
	type args struct {
		cfgFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test no file",
			args:    args{cfgFile: "./.yml"},
			wantErr: true,
		},
		{
			name:    "test root file",
			args:    args{cfgFile: "../.yml"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := new(Config)
			if err := c.Load(tt.args.cfgFile); (err != nil) != tt.wantErr {
				t.Errorf("Config.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

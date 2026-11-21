package config

import "testing"

func TestValidateAndConvertConfig_AboutTargets(t *testing.T) {
	cfg := Config{
		AboutTargets: []AboutTarget{
			{Name: "gdrive", Remote: "gdrive:"},
		},
	}

	internal, err := cfg.validateAndConvertConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(internal.AboutTargets) != 1 {
		t.Fatalf("expected 1 about target, got %d", len(internal.AboutTargets))
	}
	if internal.AboutTargets[0].Name != "gdrive" || internal.AboutTargets[0].Remote != "gdrive:" {
		t.Errorf("unexpected about target: %+v", internal.AboutTargets[0])
	}
}

func TestValidateAndConvertConfig_AboutTargetsValidationErrors(t *testing.T) {
	tests := []struct {
		name   string
		target AboutTarget
	}{
		{name: "missing name", target: AboutTarget{Remote: "gdrive:"}},
		{name: "missing remote", target: AboutTarget{Name: "gdrive"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{AboutTargets: []AboutTarget{tt.target}}
			if _, err := cfg.validateAndConvertConfig(); err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	t.Setenv("SSM_RCLONE_EXECUTABLE_PATH", "/custom/rclone")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.RcloneExecutablePath != "/custom/rclone" {
		t.Errorf("expected RcloneExecutablePath %q, got %q", "/custom/rclone", cfg.RcloneExecutablePath)
	}
}

package config_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/matansh/envloader/config"
	"github.com/matansh/envloader/errs"
)

func TestLoadEnv(t *testing.T) {
	t.Parallel()
	type testCfg struct {
		Str string `env:"str"`
	}

	t.Run("not a pointer", func(t *testing.T) {
		t.Parallel()
		errArr := config.LoadEnv(testCfg{})
		if len(errArr) != 1 {
			t.Fatalf("expected a single error, got: %d", len(errArr))
		}
		if !errors.Is(errArr[0], errs.ErrNotPtr) {
			t.Error(fmt.Errorf("err is not of type ErrNotPtr: %w", errArr[0]))
		}
	})

	t.Run("not a struct", func(t *testing.T) {
		t.Parallel()
		cfg := "not a struct"
		errArr := config.LoadEnv(&cfg)
		if len(errArr) != 1 {
			t.Fatalf("expected a single error, got: %d", len(errArr))
		}
		if !errors.Is(errArr[0], errs.ErrNotStruct) {
			t.Error(fmt.Errorf("err is not of type ErrNotStruct: %w", errArr[0]))
		}
	})

	t.Run("not in env", func(t *testing.T) {
		t.Parallel()
		os.Clearenv()
		cfg := testCfg{}
		errArr := config.LoadEnv(&cfg)
		if len(errArr) != 1 {
			t.Fatalf("expected a single error, got: %d", len(errArr))
		}
		if !errors.Is(errArr[0], errs.ErrEnvVarNotFound) {
			t.Error(fmt.Errorf("err is not of type ErrEnvVarNotFound: %w", errArr[0]))
		}
	})

	t.Run("existing value overwritten", func(t *testing.T) {
		t.Parallel()
		if err := os.Setenv("str", "test"); err != nil {
			t.Error(err)
		}
		cfg := testCfg{
			Str: "These aren’t the droids you’re looking for",
		}
		errArr := config.LoadEnv(&cfg)
		if len(errArr) != 0 {
			t.Error(errArr)
		}
		if cfg.Str != "test" {
			t.Errorf("unexpected value of struct fields: %s", cfg.Str)
		}
	})
}

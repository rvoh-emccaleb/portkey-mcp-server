package config

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrNoBuildVersion       = errors.New("no version provided at build time")
	ErrInvalidVersionFormat = errors.New(
		"value is not in a semver or git commit hash format, e.g. v1, v1.2, v1.2.3, abcde01, " +
			"abcde01abcde01abcde01abcde01abcde01f2345",
	)
)

type BuildTimeVars struct {
	// AppVersion is the version of this application that is running.
	AppVersion string `envconfig:"-"`
}

func (cfg *BuildTimeVars) Validate() error {
	err := validateAppVersion(cfg.AppVersion)
	if err != nil {
		return fmt.Errorf("error validating application version: %w", err)
	}

	return nil
}

func validateAppVersion(s string) error {
	if s == "" {
		return ErrNoBuildVersion
	}

	// Valid values:
	// - v1
	// - v1.2
	// - v1.2.3
	// - abcde01 (7 hex)
	// - abcde01abcde01abcde01abcde01abcde01f2345 (40 hex)
	semVer := `v(\d+)\.(\d+)?\.(\d+)?`
	gitCommitHash := `[a-f0-9]{7}|[a-f0-9]{40}`
	regex := fmt.Sprintf(`^(%s|%s)$`, semVer, gitCommitHash)

	r, err := regexp.Compile(regex)
	if err != nil {
		return fmt.Errorf("error compiling regex: %w", err)
	}

	if !r.MatchString(s) {
		return ErrInvalidVersionFormat
	}

	return nil
}

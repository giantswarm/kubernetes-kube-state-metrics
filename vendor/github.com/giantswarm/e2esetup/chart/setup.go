package chart

import (
	"context"
	"testing"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/e2esetup/chart/env"
)

func Setup(ctx context.Context, m *testing.M, config Config) error {
	var v int
	var err error
	var errors []error

	if config.HelmClient == nil {
		return microerror.Maskf(invalidConfigError, "%T.HelmClient must not be empty", config)
	}
	if config.Host == nil {
		return microerror.Maskf(invalidConfigError, "%T.Host must not be empty", config)
	}

	err = config.Host.CreateNamespace("giantswarm")
	if err != nil {
		errors = append(errors, err)
		v = 1
	}

	err = config.HelmClient.EnsureTillerInstalled()
	if err != nil {
		errors = append(errors, err)
		v = 1
	}

	if v == 0 {
		v = m.Run()
	}

	if env.KeepResources() != "true" {
		// Only do full teardown when not on CI.
		if env.CircleCI() != "true" {
			err := teardown(config)
			if err != nil {
				errors = append(errors, err)
				v = 1
			}
			// TODO there should be error handling for the framework teardown.
			config.Host.Teardown()
		}
	}

	if len(errors) > 0 {
		return microerror.Mask(errors[0])
	}

	return nil
}

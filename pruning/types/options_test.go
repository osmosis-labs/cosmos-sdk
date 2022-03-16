package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPruningOptions_Validate(t *testing.T) {
	testCases := []struct {
		opts   *PruningOptions
		expectErr  bool
	}{
		{NewPruningOptions(Default), false}, // default
		{NewPruningOptions(Everything), false},     // everything
		{NewPruningOptions(Nothing), false},      // nothing
		{NewCustomPruningOptions(0, 10, 10), false},
		{NewCustomPruningOptions(100, 0, 0), true}, // invalid interval
		{NewCustomPruningOptions(0, 1, 5), true},   // invalid interval
	}

	for _, tc := range testCases {
		err := tc.opts.Validate()
		require.Equal(t, tc.expectErr, err != nil, "options: %v, err: %s", tc.opts, err)
	}
}

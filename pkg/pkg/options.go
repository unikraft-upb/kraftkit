// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <alex@unikraft.io>
//
// Copyright (c) 2022, Unikraft UG.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package pkg

import (
	"context"

	"go.unikraft.io/kit/pkg/log"
	"go.unikraft.io/kit/utils"
)

// PackageOptions contains configuration for the Package
type PackageOptions struct {
	// Access to a logger
	log log.Logger

	// ctx should contain all implementation-specific options, using
	// `context.WithValue`
	ctx context.Context
}

type PackageOption func(opts *PackageOptions) error

// WithLogger defines the log.Logger
func WithLogger(l *log.Logger) PackageOption {
	return func(o *PackageOptions) error {
		o.Log = l
		return nil
	}
}

// NewPackageOptions creates PackageOptions
func NewPackageOptions(opts ...PackageOption) (*PackageOptions, error) {
	options := &PackageOptions{}

	for _, o := range opts {
		err := o(options)
		if err != nil {
			return nil, err
		}
	}

	return options, nil
}

func WithLogger(l log.Logger) PackageOption {
	return func(opts *PackageOptions) error {
		opts.log = l
		return nil
	}
}

type PullPackageOptions struct {
	architectures     []string
	platforms         []string
	calculateChecksum bool
	onProgress        func(progress float64)
	workdir           string
	log               log.Logger
	useCache          bool
}

// OnProgress calls (if set) an embedded progress function which can be used to
// update an external progress bar, for example.
func (ppo *PullPackageOptions) OnProgress(progress float64) {
	if ppo.onProgress != nil {
		ppo.onProgress(progress)
	}
}

// Workdir returns the set working directory as part of the pull request
func (ppo *PullPackageOptions) Workdir() string {
	return ppo.workdir
}

// CalculateChecksum returns whether the pull request should perform a check of
// the resource sum.
func (ppo *PullPackageOptions) CalculateChecksum() bool {
	return ppo.calculateChecksum
}

// Log returns the available logger
func (ppo *PullPackageOptions) Log() log.Logger {
	return ppo.log
}

// UseCache returns whether the pull should redirect to using a local cache if
// available.
func (ppo *PullPackageOptions) UseCache() bool {
	return ppo.useCache
}

type PullPackageOption func(opts *PullPackageOptions) error

// NewPullPackageOptions creates PullPackageOptions
func NewPullPackageOptions(opts ...PullPackageOption) (*PullPackageOptions, error) {
	options := &PullPackageOptions{}

	for _, o := range opts {
		err := o(options)
		if err != nil {
			return nil, err
		}
	}

	return options, nil
}

// WithPullArchitecture requests a given architecture (if applicable)
func WithPullArchitecture(archs ...string) PullPackageOption {
	return func(opts *PullPackageOptions) error {
		for _, arch := range archs {
			if arch == "" {
				continue
			}

			if utils.Contains(opts.architectures, arch) {
				continue
			}

			opts.architectures = append(opts.architectures, archs...)
		}

		return nil
	}
}

// WithPullPlatform requests a given platform (if applicable)
func WithPullPlatform(plats ...string) PullPackageOption {
	return func(opts *PullPackageOptions) error {
		for _, plat := range plats {
			if plat == "" {
				continue
			}

			if utils.Contains(opts.platforms, plat) {
				continue
			}

			opts.platforms = append(opts.platforms, plats...)
		}

		return nil
	}
}

// WithPullProgressFunc set an optional progress function which is used as a
// callback during the transmission of the package and the host
func WithPullProgressFunc(onProgress func(progress float64)) PullPackageOption {
	return func(opts *PullPackageOptions) error {
		opts.onProgress = onProgress
		return nil
	}
}

// WithPullWorkdir set the working directory context of the pull such that the
// resources of the package are placed there appropriately
func WithPullWorkdir(workdir string) PullPackageOption {
	return func(opts *PullPackageOptions) error {
		opts.workdir = workdir
		return nil
	}
}

// WithPullLogger set the use of a logger
func WithPullLogger(l log.Logger) PullPackageOption {
	return func(opts *PullPackageOptions) error {
		opts.log = l
		return nil
	}
}

// WithPullChecksum to set whether to calculate and compare the checksum of the
// package
func WithPullChecksum(calc bool) PullPackageOption {
	return func(opts *PullPackageOptions) error {
		opts.calculateChecksum = calc
		return nil
	}
}

// WithPullCache to set whether use cache if possible
func WithPullCache(cache bool) PullPackageOption {
	return func(opts *PullPackageOptions) error {
		opts.useCache = cache
		return nil
	}
}

// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <alex@unikraft.io>
//
// Copyright (c) 2022, Unikraft GmbH.  All rights reserved.
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

package target

import (
	"fmt"

	"kraftkit.sh/initrd"
	"kraftkit.sh/iostreams"
	"kraftkit.sh/unikraft"
	"kraftkit.sh/unikraft/arch"
	"kraftkit.sh/unikraft/component"
	"kraftkit.sh/unikraft/plat"
)

type TargetConfig struct {
	component.ComponentConfig `yaml:"-" json:"-"`

	Architecture arch.ArchitectureConfig `yaml:",omitempty" json:"architecture,omitempty"`
	Platform     plat.PlatformConfig     `yaml:",omitempty" json:"platform,omitempty"`
	Format       string                  `yaml:",omitempty" json:"format,omitempty"`
	Kernel       string                  `yaml:"-" json:"-"`
	KernelDbg    string                  `yaml:"-" json:"-"`
	Initrd       *initrd.InitrdConfig    `yaml:"-" json:"-"`
	Command      []string                `yaml:",omitempty" json:"commands"`

	Extensions map[string]interface{} `yaml:",inline" json:"-"`
}

type Targets []TargetConfig

func (tc TargetConfig) Name() string {
	return tc.ComponentConfig.Name
}

func (tc TargetConfig) Version() string {
	return tc.ComponentConfig.Version
}

func (tc TargetConfig) Type() unikraft.ComponentType {
	return unikraft.ComponentTypeUnknown
}

// ArchPlatString returns the canonical name for platform architecture string
// combination
func (tc *TargetConfig) ArchPlatString() string {
	return tc.Platform.Name() + "-" + tc.Architecture.Name()
}

func (tc TargetConfig) PrintInfo(io *iostreams.IOStreams) error {
	fmt.Fprint(io.Out, "not implemented: unikraft.plat.TargetConfig.PrintInfo")
	return nil
}

func (tc TargetConfig) MarshalYAML() (interface{}, error) {
	type targetConfigYAML struct {
		component.ComponentConfig `yaml:"-" json:"-"`

		NameYAML     string                  `yaml:"name,omitempty" json:"name,omitempty"`
		Architecture arch.ArchitectureConfig `yaml:",omitempty" json:"architecture,omitempty"`
		Platform     plat.PlatformConfig     `yaml:",omitempty" json:"platform,omitempty"`
		Format       string                  `yaml:",omitempty" json:"format,omitempty"`
		Kernel       string                  `yaml:"-" json:"-"`
		KernelDbg    string                  `yaml:"-" json:"-"`
		Initrd       *initrd.InitrdConfig    `yaml:"-" json:"-"`
		Command      []string                `yaml:",omitempty" json:"commands"`

		Extensions map[string]interface{} `yaml:",inline" json:"-"`
	}

	return targetConfigYAML{
		ComponentConfig: tc.ComponentConfig,
		NameYAML:        tc.ComponentConfig.Name,
		Architecture:    tc.Architecture,
		Platform:        tc.Platform,
		Format:          tc.Format,
		Kernel:          tc.Kernel,
		KernelDbg:       tc.KernelDbg,
		Initrd:          tc.Initrd,
		Command:         tc.Command,
		Extensions:      tc.Extensions,
	}, nil
}

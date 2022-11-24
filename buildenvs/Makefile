# SPDX-License-Identifier: BSD-3-Clause
# Copyright (c) 2022, Unikraft GmbH and The KraftKit Authors.
# Licensed under the BSD-3-Clause License (the "License")
# You may not use this file expect in compliance with the License.

# Meta directory association
_MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
BUILDENVSDIR := $(patsubst %/,%,$(dir $(_MKFILE_PATH)))
WORKDIR      ?= $(BUILDENVSDIR)/..

REG          ?= ghcr.io
ORG          ?= unikraft
APP          ?= kraftkit
IMAGE_TAG    ?= latest

BUILDENVS    ?= gcc \
                qemu \
                myself \
                base

DOCKER       ?= docker

.PHONY: qemu
qemu: QEMU_VERSION ?= 7.1.0
qemu: ENVIRONMENT  ?= qemu
qemu: IMAGE        ?= $(REG)/$(ORG)/$(APP)/$(ENVIRONMENT):$(QEMU_VERSION)
qemu:
	$(DOCKER) build \
		--build-arg QEMU_VERSION=$(QEMU_VERSION) \
		--tag $(IMAGE) \
		--cache-from $(IMAGE) \
		--file $(BUILDENVSDIR)/$(ENVIRONMENT).Dockerfile \
		$(DOCKER_BUILD_EXTRA) $(WORKDIR)

.PHONY: gcc
gcc: UK_ARCH          ?= x86_64
gcc: GCC_VERSION      ?= 12.2.0
gcc: BINUTILS_VERSION ?= 2.39
gcc: GLIB_VERSION     ?= 2.31
gcc: ENVIRONMENT      ?= gcc
gcc: IMAGE            ?= $(REG)/$(ORG)/$(APP)/$(ENVIRONMENT):$(GCC_VERSION)-$(UK_ARCH)
gcc:
	$(DOCKER) build \
		--build-arg UK_ARCH=$(UK_ARCH) \
		--build-arg GCC_VERSION=$(GCC_VERSION) \
		--build-arg BINUTILS_VERSION=$(BINUTILS_VERSION) \
		--build-arg GLIB_VERSION=$(GLIB_VERSION) \
		--tag $(IMAGE) \
		--cache-from $(IMAGE) \
		--file $(BUILDENVSDIR)/$(ENVIRONMENT).Dockerfile \
		$(DOCKER_BUILD_EXTRA) $(WORKDIR)

.PHONY: myself
myself: GO_VERSION         ?= 1.18
myself: ENVIRONMENT        ?= myself
myself: IMAGE              ?= $(REG)/$(ORG)/$(APP)/$(ENVIRONMENT):$(IMAGE_TAG)
myself:
	$(DOCKER) build \
		--build-arg ORG=$(ORG) \
		--build-arg APP=$(APP) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--tag $(IMAGE) \
		--cache-from $(IMAGE) \
		--file $(BUILDENVSDIR)/$(ENVIRONMENT).Dockerfile \
		$(DOCKER_BUILD_EXTRA) $(WORKDIR)

.PHONY: base
base: GO_VERSION         ?= 1.18
base: UK_ARCH            ?= x86_64
base: GCC_VERSION        ?= 12.2.0
base: QEMU_VERSION       ?= 7.1.0
base: ENVIRONMENT        ?= base
base: IMAGE              ?= $(REG)/$(ORG)/$(APP)/$(ENVIRONMENT):$(IMAGE_TAG)
base:
	$(DOCKER) build \
		--build-arg REG=$(REG) \
		--build-arg UK_ARCH=$(UK_ARCH) \
		--build-arg GCC_VERSION=$(GCC_VERSION) \
		--build-arg GCC_PREFIX=$(GCC_PREFIX) \
		--build-arg QEMU_VERSION=$(QEMU_VERSION) \
		--build-arg ORG=$(ORG) \
		--build-arg APP=$(APP) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--tag $(IMAGE) \
		--cache-from $(IMAGE) \
		--file $(BUILDENVSDIR)/$(ENVIRONMENT).Dockerfile \
		$(DOCKER_BUILD_EXTRA) $(WORKDIR)

# SPDX-License-Identifier: BSD-3-Clause
#
# Authors: Alexander Jung <alexander.jung@neclab.eu>
#          Cezar Craciunoiu <cezar.craciunoiu@unikraft.io>
#
# Copyright (c) 2020, NEC Europe Ltd., NEC Corporation. All rights reserved.
# Copyright (c) 2022, Unikraft GmbH. All rights reserved.
#
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions
# are met:
#
# 1. Redistributions of source code must retain the above copyright
#    notice, this list of conditions and the following disclaimer.
# 2. Redistributions in binary form must reproduce the above copyright
#    notice, this list of conditions and the following disclaimer in the
#    documentation and/or other materials provided with the distribution.
# 3. Neither the name of the copyright holder nor the names of its
#    contributors may be used to endorse or promote products derived from
#    this software without specific prior written permission.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
# AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
# IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
# ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
# LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
# CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
# SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
# INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
# CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
# ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
# POSSIBILITY OF SUCH DAMAGE.

ARG REG=

ARG UK_ARCH=x86_64

ARG GCC_VERSION=12.2.0
ARG GCC_SUFFIX=

ARG QEMU_VERSION=7.1.0

ARG DEBIAN_VERSION=bullseye-20221114

FROM ${REG}/unikraft/kraftkit/gcc:${GCC_VERSION}-x86_64${GCC_SUFFIX} AS gcc-x86_64
# FROM ${REG}/unikraft/kraftkit/gcc:${GCC_VERSION}-arm${GCC_SUFFIX}    AS gcc-arm
# FROM ${REG}/unikraft/kraftkit/gcc:${GCC_VERSION}-arm64${GCC_SUFFIX}  AS gcc-arm64

FROM ${REG}/unikraft/kraftkit/qemu:${QEMU_VERSION}                   AS qemu
FROM ${REG}/unikraft/kraftkit/myself:latest                          AS kraft

LABEL MAINTAINER="Alexander Jung <alexander.jung@neclab.eu>"
LABEL MAINTAINER="Cezar Craciunoiu <cezar.craciunoiu@unikraft.io>"

# Copy the kraft tool
COPY . /go/src/kraftkit.sh

# Build the kraft tool - also installs it to the environment
WORKDIR /go/src/kraftkit.sh

RUN set -xe; \
  make kraft;

FROM debian:${DEBIAN_VERSION} AS kraftkit

ARG GCC_PREFIX=x86_64-linux-gnu

# Install kraftkit dependencies
COPY --from=kraft /go/src/kraftkit.sh/dist/kraft /bin
COPY --from=kraft /root/.config/kraftkit/config.yaml /root/.config/kraftkit/config.yaml
COPY --from=kraft /root/.local/share/kraftkit /root/.local/share/kraftkit

# Copy the GCC toolchain
COPY --from=gcc-x86_64 /bin/                                         /bin
COPY --from=gcc-x86_64 /lib/gcc/                                     /lib/gcc
COPY --from=gcc-x86_64 /include/                                     /include
COPY --from=gcc-x86_64 /x86_64-linux-gnu                             /x86_64-linux-gnu
COPY --from=gcc-x86_64 /libexec/gcc/x86_64-linux-gnu/12.2.0/cc1      /libexec/gcc/x86_64-linux-gnu/12.2.0/cc1
COPY --from=gcc-x86_64 /libexec/gcc/x86_64-linux-gnu/12.2.0/collect2 /libexec/gcc/x86_64-linux-gnu/12.2.0/collect2
# COPY the others if needed

# COPY --from=gcc-arm /bin/                                            /bin
# COPY --from=gcc-arm /lib/gcc/                                        /lib/gcc
# COPY --from=gcc-arm /include/                                        /include
# COPY --from=gcc-arm /lib/gcc/                                        /lib/gcc
# COPY --from=gcc-arm /arm-linux-gnueabihf                             /arm-linux-gnueabihf
# COPY --from=gcc-arm /libexec/gcc/arm-linux-gnueabihf/12.2.0/cc1      /libexec/gcc/arm-linux-gnueabihf/12.2.0/cc1
# COPY --from=gcc-arm /libexec/gcc/arm-linux-gnueabihf/12.2.0/collect2 /libexec/gcc/arm-linux-gnueabihf/12.2.0/collect2
# COPY the others if needed

# COPY --from=gcc-arm64 /bin/                                          /bin
# COPY --from=gcc-arm64 /lib/gcc/                                      /lib/gcc
# COPY --from=gcc-arm64 /include/                                      /include
# COPY --from=gcc-arm64 /lib/gcc/                                      /lib/gcc
# COPY --from=gcc-arm64 /aarch64-linux-gnu/                            /aarch64-linux-gnu
# COPY --from=gcc-arm64 /libexec/gcc/aarch64-linux-gnu/12.2.0/cc1      /libexec/gcc/aarch64-linux-gnu/12.2.0/cc1
# COPY --from=gcc-arm64 /libexec/gcc/aarch64-linux-gnu/12.2.0/collect2 /libexec/gcc/aarch64-linux-gnu/12.2.0/collect2
# COPY the others if needed

# Link the GCC toolchain
RUN set -xe; \
    ln -s /bin/${GCC_PREFIX}-as         /bin/as; \
    ln -s /bin/${GCC_PREFIX}-ar         /bin/ar; \
    ln -s /bin/${GCC_PREFIX}-c++        /bin/c++; \
    ln -s /bin/${GCC_PREFIX}-c++filt    /bin/c++filt; \
    ln -s /bin/${GCC_PREFIX}-elfedit    /bin/elfedit; \
    ln -s /bin/${GCC_PREFIX}-gcc        /bin/cc; \
    ln -s /bin/${GCC_PREFIX}-gcc        /bin/gcc; \
    ln -s /bin/${GCC_PREFIX}-gcc-ar     /bin/gcc-ar; \
    ln -s /bin/${GCC_PREFIX}-gcc-nm     /bin/gcc-nm; \
    ln -s /bin/${GCC_PREFIX}-gcc-ranlib /bin/gcc-ranlib; \
    ln -s /bin/${GCC_PREFIX}-gccgo      /bin/gccgo; \
    ln -s /bin/${GCC_PREFIX}-gcov       /bin/gcov; \
    ln -s /bin/${GCC_PREFIX}-gcov-dump  /bin/gcov-dump; \
    ln -s /bin/${GCC_PREFIX}-gcov-tool  /bin/gcov-tool; \
    ln -s /bin/${GCC_PREFIX}-gprof      /bin/gprof; \
    ln -s /bin/${GCC_PREFIX}-ld         /bin/ld; \
    ln -s /bin/${GCC_PREFIX}-nm         /bin/nm; \
    ln -s /bin/${GCC_PREFIX}-objcopy    /bin/objcopy; \
    ln -s /bin/${GCC_PREFIX}-objdump    /bin/objdump; \
    ln -s /bin/${GCC_PREFIX}-ranlib     /bin/ranlib; \
    ln -s /bin/${GCC_PREFIX}-readelf    /bin/readelf; \
    ln -s /bin/${GCC_PREFIX}-size       /bin/size; \
    ln -s /bin/${GCC_PREFIX}-strings    /bin/strings; \
    ln -s /bin/${GCC_PREFIX}-strip      /bin/strip;

# Copy the QEMU emulator
COPY --from=qemu /bin/                  /bin
COPY --from=qemu /share/qemu/           /share/qemu
COPY --from=qemu /lib/x86_64-linux-gnu/ /lib/x86_64-linux-gnu

# COPY --from=qemu /lib/arm-linux-gnueabihf/ /lib/arm-linux-gnueabihf
# COPY --from=qemu /lib/aarch64-linux-gnu/ /lib/aarch64-linux-gnu

# Unikraft workdir and environment
WORKDIR /usr/src/unikraft/apps/app
ENV UK_WORKDIR /usr/src/unikraft
ENV UK_CACHEDIR /usr/src/unikraft/.kraftcache
ENV KRAFTRC /usr/src/unikraft/.kraftrc

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

# Install unikraft dependencies
RUN set -xe; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
      make \
      libncursesw5-dev \
      libncursesw5 \
      libyaml-dev \
      flex \
      git \
      wget \
      patch \
      gawk \
      socat \
      bison \
      bindgen \
      bzip2 \
      unzip \
      uuid-runtime \
      openssh-client \
      autoconf \
      xz-utils \
      python3 \
      libgit2-dev \
      ca-certificates; \
    apt-get clean; \
    rm -Rf /var/cache/apt/*; \
    mkdir -p /usr/src/unikraft/unikraft; \
    mkdir -p /usr/src/unikraft/libs; \
    mkdir -p /usr/src/unikraft/apps/app; \
    kraft -h;

ENTRYPOINT [ "kraft" ]
# SPDX-License-Identifier: BSD-3-Clause
#
# Authors: Alexander Jung <alexander.jung@neclab.eu>
#
# Copyright (c) 2020, NEC Europe Ltd., NEC Corporation. All rights reserved.
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

ARG DEBIAN_VERSION=bullseye-20221114

FROM debian:${DEBIAN_VERSION} AS qemu-build

ARG QEMU_VERSION=7.1.0
ARG WITH_XEN=disable
ARG WITH_KVM=enable

ARG WITH_x86_64=enable
ARG WITH_aarch64=disable
ARG WITH_arm=disable


WORKDIR /out

# Install dependencies
RUN set -ex; \
    apt-get -y update; \
    apt-get install -y \
        bison \
        build-essential \
        curl \
        flex \
        libaio-dev \
        libattr1-dev \
        libcap-dev \
        libcap-ng-dev \
        libglib2.0-dev \
        liblzo2-dev \
        libpixman-1-dev \
        ninja-build \
        pkg-config \
        python \
        texinfo \
        vde2 \
        xz-utils \
        zlib1g-dev; \
    apt-get clean;

# Download and extract QEMU
RUN set -ex; \
    curl -O https://download.qemu.org/qemu-${QEMU_VERSION}.tar.xz; \
    tar xf qemu-${QEMU_VERSION}.tar.xz; \
    apt-get install -y;

# Configure and build QEMU
RUN set -ex; \
    cd qemu-${QEMU_VERSION}; \
    tlist=""; \
    if [ "${WITH_x86_64}" = "enable" ]; then \
        tlist="x86_64-softmmu"; \
    fi; \
    if [ "${WITH_aarch64}" = "enable" ]; then \
        tlist="${tlist},aarch64-softmmu"; \
    fi; \
    if [ "${WITH_arm}" = "enable" ]; then \
        tlist="${tlist},arm-softmmu"; \
    fi; \
    ./configure \
        --target-list=${tlist} \
        --static \
        --prefix=/ \
        --audio-drv-list="" \
        --enable-attr \
        --disable-auth-pam \
        --disable-avx2 \
        --disable-avx512f \
        --disable-bochs \
        --disable-bpf \
        --disable-brlapi \
        --disable-bsd-user \
        --disable-bzip2 \
        --disable-canokey \
        --disable-capstone \
        --disable-cfi \
        --disable-cfi-debug \
        --disable-cloop \
        --disable-cocoa \
        --disable-coreaudio \
        --disable-crypto-afalg \
        --disable-curl \
        --disable-curses \
        --disable-dbus-display \
        --disable-dmg \
        --disable-docs \
        --disable-dsound \
        --disable-fuse \
        --disable-fuse-lseek \
        --disable-gcov \
        --disable-gcrypt \
        --disable-gettext \
        --disable-gio \
        --disable-glusterfs \
        --disable-gnutls \
        --disable-gprof \
        --disable-gtk \
        --disable-guest-agent \
        --disable-guest-agent-msi \
        --disable-hax \
        --disable-hvf \
        --disable-iconv \
        --disable-jack \
        --disable-keyring \
        --${WITH_KVM}-kvm \
        --disable-l2tpv3 \
        --disable-libdaxctl \
        --disable-libiscsi \
        --disable-libnfs \
        --disable-libpmem \
        --disable-libssh \
        --disable-libudev \
        --disable-libusb \
        --disable-libvduse \
        --disable-linux-aio \
        --disable-linux-io-uring \
        --disable-linux-user \
        --disable-live-block-migration \
        --disable-lzfse \
        --enable-lzo \
        --disable-malloc-trim \
        --disable-membarrier \
        --disable-modules \
        --disable-mpath \
        --disable-multiprocess \
        --disable-netmap \
        --disable-nettle \
        --disable-numa \
        --disable-nvmm \
        --disable-opengl \
        --disable-oss \
        --disable-pa \
        --disable-parallels \
        --disable-pie \
        --disable-png \
        --disable-profiler \
        --disable-pvrdma \
        --disable-qcow1 \
        --disable-qed \
        --disable-qga-vss \
        --disable-rbd \
        --disable-rdma \
        --disable-replication \
        --disable-safe-stack \
        --disable-sdl \
        --disable-sdl-image \
        --disable-seccomp \
        --disable-selinux \
        --disable-slirp-smbd \
        --disable-smartcard \
        --disable-snappy \
        --disable-sparse \
        --disable-spice \
        --disable-spice-protocol \
        --disable-tcg \
        --enable-tools \
        --disable-tpm \
        --disable-u2f \
        --disable-usb-redir \
        --disable-user \
        --disable-vde \
        --disable-vdi \
        --disable-vduse-blk-export \
        --disable-vfio-user-server \
        --enable-vhost-crypto \
        --enable-vhost-kernel \
        --enable-vhost-net \
        --enable-vhost-user \
        --enable-vhost-user-blk-server \
        --enable-vhost-vdpa \
        --disable-virglrenderer \
        --disable-virtiofsd \
        --disable-vmnet \
        --disable-vnc \
        --disable-vnc-jpeg \
        --disable-vnc-sasl \
        --disable-vte \
        --disable-vvfat \
        --disable-werror \
        --disable-whpx \
        --${WITH_XEN}-xen \
        --${WITH_XEN}-xen-pci-passthrough \
        --disable-xkbcommon \
        --disable-zstd \
        --enable-virtfs \
        ; \
        make -j$(($(nproc) - 1)); \
        make install;

FROM scratch AS qemu

COPY --from=qemu-build /bin/qemu-system-x86_64 /bin/qemu-g[a] \
    /bin/qemu-system-i38[6] /bin/qemu-system-ar[m] /bin/qemu-system-aarch6[4] \
    /bin/qemu-pr-helpe[r] /bin/qemu-im[g] /bin/qemu-i[o] /bin/qemu-nb[d] /bin/
COPY --from=qemu-build /share/qemu/ /share/qemu/
COPY --from=qemu-build /lib/x86_64-linux-gnu/ /lib/x86_64-linux-gnu/

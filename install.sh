#!/bin/sh

LIBVIRT_HOOKS_DIR=/etc/libvirt/hooks

go build .

mkdir -p "${LIBVIRT_HOOKS_DIR}"

sudo cp libvirt-portfwd-hooks "${LIBVIRT_HOOKS_DIR}"

sudo ln -s ${LIBVIRT_HOOKS_DIR}/hooks ${LIBVIRT_HOOKS_DIR}/qemu
sudo ln -s ${LIBVIRT_HOOKS_DIR}/hooks ${LIBVIRT_HOOKS_DIR}/network
sudo ln -s ${LIBVIRT_HOOKS_DIR}/hooks ${LIBVIRT_HOOKS_DIR}/lxc

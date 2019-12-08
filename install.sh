#!/bin/sh

LIBVIRT_HOOKS_DIR=/etc/libvirt/hooks

go build .

mkdir -p "${LIBVIRT_HOOKS_DIR}"

cp libvirt-portfwd-hooks "${LIBVIRT_HOOKS_DIR}"

ln -s ${LIBVIRT_HOOKS_DIR}/hooks ${LIBVIRT_HOOKS_DIR}/qemu
ln -s ${LIBVIRT_HOOKS_DIR}/hooks ${LIBVIRT_HOOKS_DIR}/network
ln -s ${LIBVIRT_HOOKS_DIR}/hooks ${LIBVIRT_HOOKS_DIR}/lxc

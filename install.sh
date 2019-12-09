#!/bin/sh

LIBVIRT_HOOKS_DIR=/etc/libvirt/hooks

go build .

mkdir -p "${LIBVIRT_HOOKS_DIR}"

sudo cp libvirt-portfwd-hooks "${LIBVIRT_HOOKS_DIR}"

sudo ln -fs ${LIBVIRT_HOOKS_DIR}/libvirt-portfwd-hooks ${LIBVIRT_HOOKS_DIR}/qemu
sudo ln -fs ${LIBVIRT_HOOKS_DIR}/libvirt-portfwd-hooks ${LIBVIRT_HOOKS_DIR}/network
sudo ln -fs ${LIBVIRT_HOOKS_DIR}/libvirt-portfwd-hooks ${LIBVIRT_HOOKS_DIR}/lxc

sudo touch /var/log/libvirt-portfwd-hooks.log

if [ ! -f "${LIBVIRT_HOOKS_DIR}/hooks.json" ]; then
  sudo cp configs/empty.json "${LIBVIRT_HOOKS_DIR}/hooks.json"
fi

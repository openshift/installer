#!/bin/bash
sudo virsh destroy sno-test
sudo virsh undefine sno-test
sudo virsh vol-delete --pool default sno-test.qcow2
#!/bin/bash

sudo apt-get update -qq
sudo apt-get install qemu-utils -y

wget https://aka.ms/downloadazcopy-v10-linux -O azcopy.tar
tar xvf azcopy.tar --strip-components 1
chmod +x azcopy

wget ${vhd_url} -O fcos.vhd.xz
xz -d fcos.vhd.xz

echo "Convert vhd file to raw image..."
qemu-img convert -f vpc -O raw fcos.vhd fcos.raw

# Calculcate new size which is a multiple of 1MByte (rounded up).
rawdisk="fcos.raw"
vhddisk="fcos.vhd"
MB=$((1024*1024))
size=$(qemu-img info -f raw --output json "$rawdisk" | gawk 'match($0, /"virtual-size": ([0-9]+),/, val) {print val[1]}')
rounded_size=$((($size/$MB + 1) * $MB))

echo "Resize raw image..."
qemu-img resize fcos.raw $rounded_size
qemu-img convert -f raw -o subformat=fixed,force_size -O vpc fcos.raw fcos-fixed.vhd

# Upload fixed size VHD file to Azure blob storage.
echo "Upload fixed size VHD image to Azure blob storage..."
./azcopy copy ./fcos-fixed.vhd '${primary_blob_endpoint}${container_name}/fcos.vhd${sas_token}'


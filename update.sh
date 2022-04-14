#!/bin/bash

set -ex

cd /tmp

echo "===> Downloading AV signature updates..."
 wget --progress=bar:force -Nc -r -nd -l1 -A "avvepo?????dat.zip" http://download.nai.com/products/DatFiles/4.x/nai/

for avvepo in `ls avvepo*`; do
  echo " > unzipping $avvepo"
  unzip -o $avvepo
done

for avvdat in `ls avvdat-*`; do
  echo " > unzipping $avvepo"
  unzip -o $avvdat -d /usr/local/uvscan
done

echo "===> Decompressing signatures..."
/usr/local/uvscan/uvscan --decompress

echo " * Clean up unnecessary files"
rm -rf /tmp/*

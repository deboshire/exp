#!/bin/bash

set -ue

mkdir data
cd data
wget http://vasc.ri.cmu.edu/idb/images/face/frontal_images/images.tar
tar -xf images.tar
rm images.tar

#!/usr/bin/env bash

cd testdata &&
  ffmpeg -i sample-mp4-file.mp4 \
    -profile:v baseline \
    -level 3.0 \
    -s 640x360 \
    -start_number 0 \
    -hls_time 10 \
    -hls_list_size 0 \
    -f hls index.m3u8

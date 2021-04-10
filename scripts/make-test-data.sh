#!/bin/env bash

LEN=3
MAXLEN=15
function get_length {
  LEN=$(( ((${LEN} + 1) % ${MAXLEN}) + 1 ))
  echo ${LEN}
}
FREQ=5
function get_freq {
  FREQ=$(( ((${FREQ} + 1) % 9) + 2 ))
  echo ${FREQ}
}

COUNT=0

DEST=testing/test-data/one

mkdir -p ${DEST}
rm -rf ${DEST}/*
rm -rf ${DEST}/*

for i in $(seq 1 30); do
FREQ=$(get_freq)
LEN=$(get_length)
COUNT=$(( ${COUNT}+1 ))
echo "C=${COUNT} F=${FREQ} L=${LEN}"

ffmpeg -f lavfi -i testsrc=duration=${LEN}:size=640x320:rate=25 "${DEST}/video.mpg"

ffmpeg -f lavfi -i "sine=frequency=${FREQ}:duration=${LEN}" "${DEST}/audio.mp3"

ffmpeg -i "${DEST}/video.mpg" -i "${DEST}/audio.mp3" -c copy -map 0:v:0 -map 1:a:0 -shortest "${DEST}/sample ${COUNT}.mpg"

rm -rf ${DEST}/audio.mp3*
rm -rf ${DEST}/video.mpg*
done
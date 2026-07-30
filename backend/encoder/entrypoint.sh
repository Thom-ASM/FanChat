#!/bin/sh

set -eu

: "${YOUTUBE_RTMPS_URL:?YOUTUBE_RTMPS_URL is required}"
: "${YOUTUBE_STREAM_KEY:?YOUTUBE_STREAM_KEY is required}"

VIDEO_PATH="${VIDEO_PATH:-/app/slate.mp4}"

if [ ! -r "$VIDEO_PATH" ]; then
    echo "Video file is not readable: $VIDEO_PATH" >&2
    exit 1
fi

RTMPS_DESTINATION="${YOUTUBE_RTMPS_URL%/}/${YOUTUBE_STREAM_KEY}"

echo "Starting FanChat publisher"
echo "Video: $VIDEO_PATH"

exec ffmpeg \
    -hide_banner \
    -loglevel info \
    -nostdin \
    -re \
    -stream_loop -1 \
    -i "$VIDEO_PATH" \
    -map 0:v:0 \
    -map 0:a:0 \
    -c copy \
    -flvflags no_duration_filesize \
    -f flv \
    "$RTMPS_DESTINATION"

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ffmpeg \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN useradd \
    --create-home \
    --uid 10001 \
    encoder

WORKDIR /app

COPY --chown=encoder:encoder ./assets/slate.mp4 /app/slate.mp4
COPY --chown=encoder:encoder ./entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

USER encoder

STOPSIGNAL SIGTERM

ENTRYPOINT ["/app/entrypoint.sh"]

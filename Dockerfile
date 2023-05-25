FROM python:3.11-alpine

ENV \
    RCLONE_HOST=localhost \
    RCLONE_PORT=5572 \
    RCLONE_PROTOCOL=https \
    BACKUP_SCHEDULE="0 0 * * 0" \
    BACKUP_SOURCE=/data \
    BACKUP_REMOTE=remote \
    BACKUP_DEST=/backup

RUN adduser -D rclone

COPY --chmod=755 entrypoint.sh /

WORKDIR /home/rclone
COPY --chown=rclone:rclone --chmod=0755 . .
RUN pip install -r requirements.txt

ENTRYPOINT [ "/entrypoint.sh" ]
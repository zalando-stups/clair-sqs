#!/bin/sh

# setup configurations from environment variables
echo "CLAIR_DATABASE_SOURCE: $CLAIR_DATABASE_SOURCE"
sed -i 's/CLAIR_DATABASE_SOURCE/'$CLAIR_DATABASE_SOURCE'/' /etc/clair/config.yaml

echo "CLAIR_API_PAGINATIONKEY: $CLAIR_API_PAGINATIONKEY"
sed -i 's/CLAIR_API_PAGINATIONKEY/'$CLAIR_API_PAGINATIONKEY'/' /etc/clair/config.yaml

# let supervisord manage the processes
exec supervisord -c /etc/supervisord.conf

#!/bin/sh

if [ -z "$6" ]; then
    echo "Usage: $0 <sqs-queue> <sqs-region> <registry> <repository> <artifact> <version> [auth]" >&2
    exit 1
fi

sqs_queue=$1
sqs_region=$2
registry=$3
repository=$4
artifact=$5
version=$6
auth=$7

if [ -z "$auth" ]; then
    alias curl="curl -s"
else
    alias curl="curl -s -u '$auth'"
fi

schema="$(curl "https://$registry/v2/$repository/$artifact/manifests/$version" | jq -r '.schemaVersion')"
if [ -z "$schema" ]; then
    echo "Invalid arguments." >&2
    exit 2
fi
if [ "s$schema" != "s1" ] && [ "s$schema" != "s2" ]; then
    echo "Schema '$schema' not supported." >&2
    exit 3
fi

# a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4 is the 0-byte layer of only metadata, no need to look at it

if [ "s$schema" = "s1" ]; then
    layers=$(curl "https://$registry/v2/$repository/$artifact/manifests/$version" | jq -r '.fsLayers[].blobSum' | grep -v 'a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4' | tac)
else
    layers=$(curl "https://$registry/v2/$repository/$artifact/manifests/$version" | jq -r '.layers[].digest' | grep -v 'a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4')
fi

parent=
json=
for layer in $layers; do
    [ -z "$json" ] && json="[" || json="${json}, "
    json="${json}{\"Layer\": {\"Name\": \"$layer\", \"ParentName\": \"$parent\", \"Path\": \"https://$registry/v2/$repository/$artifact/blobs/$layer\", \"Format\": \"Docker\"}}"
    parent=$layer
done
json="${json}]"

echo -n "Sending layers to SQS...   "
aws sqs send-message \
    --queue-url $sqs_queue \
    --region $sqs_region \
    --message-body \
    "$json" | jq -r '.MessageId'
if [ $? -ne 0 ]; then
    echo "FAILED"
    exit 4
fi

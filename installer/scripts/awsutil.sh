#!/bin/bash

function check_aws_creds {
  if [ -z "$AWS_ACCESS_KEY_ID" ] || [ -z "$AWS_SECRET_ACCESS_KEY" ];then
      echo "Must export both \$AWS_ACCESS_KEY_ID and \$AWS_SECRET_ACCESS_KEY"
      return 1
  fi
}

function aws_upload_file {
  check_aws_creds

  file=$1
  dest=$2
  bucket=$3
  contentType=$4

  resource="/${bucket}/${dest}"
  dateValue=$(date -u +"%a, %d %b %Y %H:%M:%S GMT")
  stringToSign="PUT\n\n${contentType}\n${dateValue}\n${resource}"

  signature="$(echo -en "${stringToSign}" | openssl sha1 -hmac "${AWS_SECRET_ACCESS_KEY}" -binary | base64)"

  url="https://${bucket}.s3.amazonaws.com/${dest}"

  echo "uploading $file to $resource"
  curl \
      --fail \
      --upload-file "${file}" \
      -H "Host: ${bucket}.s3.amazonaws.com" \
      -H "Date: ${dateValue}" \
      -H "Content-type: ${contentType}" \
      -H "Authorization: AWS ${AWS_ACCESS_KEY_ID}:${signature}" \
      "${url}"

  echo "uploaded $file to $resource"
}

function aws_download_file {
  check_aws_creds

  file=$1
  dest=$2
  bucket=$3
  contentType=$4

  resource="/${bucket}/${file}"
  dateValue=$(date -u +"%a, %d %b %Y %H:%M:%S GMT")
  stringToSign="GET\n\n${contentType}\n${dateValue}\n${resource}"

  signature="$(echo -en "${stringToSign}" | openssl sha1 -hmac "${AWS_SECRET_ACCESS_KEY}" -binary | base64)"

  url="https://${bucket}.s3.amazonaws.com/${file}"

  echo "downloading $url to $dest"
  curl \
      --fail \
      -H "Host: ${bucket}.s3.amazonaws.com" \
      -H "Date: ${dateValue}" \
      -H "Content-type: ${contentType}" \
      -H "Authorization: AWS ${AWS_ACCESS_KEY_ID}:${signature}" \
      "${url}" -o "${dest}"

  echo "downloaded $url to $dest"
}

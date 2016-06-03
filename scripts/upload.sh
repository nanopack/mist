#!/usr/bin/env bash
set -e

# look through each os/arch/file and generate an md5 for each
echo "Generating md5s..."
for os in $(ls ./build); do
  for arch in $(ls ./build/${os}); do
    for file in $(ls ./build/${os}/${arch}); do
      cat "./build/${os}/${arch}/${file}" | ${MD5} >> "./build/${os}/${arch}/${file}.md5"
    done
  done
done

# upload to AWS S3
echo "Uploading builds to S3..."
aws s3 sync ./build/ s3://tools.nanopack.io/mist --grants read=uri=http://acs.amazonaws.com/groups/global/AllUsers --region us-east-1

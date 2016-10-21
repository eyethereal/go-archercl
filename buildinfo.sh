#!/bin/bash

bgo="`dirname $0`/build.go"

numCommits=`git rev-list HEAD | wc -l`
rpShort=`git rev-parse --short HEAD`
rpLong=`git rev-parse HEAD`

if [ `uname` == 'Darwin' ]; then
    epoch=`date +%s`
else
    refFile=/tmp/date_ref
    touch $refFile
fi

echo "Writing build info to $bgo"

echo "package config" > $bgo
echo "func init() {" >> $bgo
echo " BuildInfo = \`" >> $bgo

echo " git commits = $numCommits" >>$bgo
echo " git rev = \"$rpLong\"" >>$bgo
echo " git short = \"$rpShort\"" >>$bgo

if [ x != "x${TRAVIS}" ]; then
    # Running on travis

echo " travis: { " >>$bgo

echo "  server: true" >> $bgo
echo "  branch: \"$TRAVIS_BRANCH\"" >> $bgo
echo "  build id: \"$TRAVIS_BUILD_ID\"" >> $bgo
echo "  build number: \"$TRAVIS_BUILD_NUMBER\"" >> $bgo
echo "  job id: \"$TRAVIS_JOB_ID\"" >> $bgo
echo "  job number: \"$TRAVIS_JOB_NUMBER\"" >> $bgo
echo "  commit range: \"$TRAVIS_COMMIT_RANGE\"" >> $bgo

echo " } " >> $bgo

fi


echo " go: { " >>$bgo
echo "  version: \"`go version`\"" >> $bgo
echo " } " >> $bgo

echo " uname: \"`uname -a`\"" >>$bgo

if [ `uname` == 'Darwin' ]; then
    echo " time local: \"$(TZ='PDT+7' date -r $epoch +%c\ %Z)\"" >> $bgo
    echo " time stamp: \"$(TZ='PDT+7' date -r $epoch +%Y%m%d-%H%M%S)\"" >> $bgo
    echo " time epoch: $epoch"  >> $bgo
else
    echo " time local: \"$(TZ='PDT+7' date -r $refFile +%c\ %Z)\"" >> $bgo
    echo " time stamp: \"$(TZ='PDT+7' date -r $refFile +%Y%m%d-%H%M%S)\"" >> $bgo
    echo " time epoch: $(date -r $refFile +%s)"  >> $bgo    
fi


echo "\` } "  >> $bgo

echo -e "\n\033[1;31mcat of the build info file ....\n"
cat $bgo
echo -e "\n\033[0m"
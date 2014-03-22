#!/bin/bash

echo "mode: set" > acc.out

for dir in $(find . -maxdepth 10 -not -path './.git*' -type d);
do
	if ls $dir/*.go &> /dev/null;
	then
		returnval=`go test -coverprofile=profile.out $dir`
		echo ${returnval}
		if [[ ${returnval} != *FAIL* ]]
		then
    		if [ -f profile.out ]
    		then
        		cat profile.out | grep -v "mode: set" >> acc.out
    		fi
    	else
    		exit 1
    	fi
    fi
done
if [ -n "$COVERALLS" ]
then
	goveralls -v -coverprofile=acc.out $COVERALLS
fi

rm -rf ./profile.out
rm -rf ./acc.out

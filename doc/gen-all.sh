#!/bin/sh
for FILE in `ls *.dot`; do circo -Tpng $FILE -O; open $FILE.png; done

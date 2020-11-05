#!/bin/bash
VERSION='15_2'

FILE="zv$VERSION.yaml"
OUT="out$VERSION"

csplit --suppress-matched $FILE /---/ '{*}'
mkdir $OUT
mv xx* $OUT

for filename in $OUT/*; do
    name=$(sed -n '/name:/p' $filename | head -1 | awk '{print $2}')
    mv $filename "$OUT/$name.yaml"
done
#!/bin/bash

i=0
for img in "out"/*.jpg; do
        [ -e "${img}" ] || break
        echo mv "${img}" "out/$i.jpg"
        ((i++))
done
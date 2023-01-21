#!/bin/bash

dir='outmini'
i=0
for img in "$dir"/*.jpg; do
        [ -e "${img}" ] || break
        mv "${img}" "$dir/$i.jpg"
        ((i++))
done
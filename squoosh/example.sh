#!/bin/bash
podman build . -t squoosh
podman run -v ./photos2022/:/photos/ -v ./out/:/out/ localhost/squoosh:latest --mozjpeg {quality:75} -d /out /photos/
podman run -v ./photos2022/:/photos/ -v ./outmini/:/out/ localhost/squoosh:latest --mozjpeg {quality:15} -d /out /photos/

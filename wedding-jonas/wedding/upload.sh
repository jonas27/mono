if [[ "$*" == "build" ]] ; then
    flutter build web
fi
cp -f /home/joe/repos/mono/wedding/wedding/assets/images/app_icons/favicon.png /home/joe/repos/mono/wedding/wedding/build/web/
# cp -f /home/joe/repos/mono/wedding/wedding/assets/images/app_icons/Icon-192.png /home/joe/repos/mono/wedding/wedding/build/web/icons/
# cp -f /home/joe/repos/mono/wedding/wedding/assets/images/app_icons/Icon-512.png /home/joe/repos/mono/wedding/wedding/build/web/icons/
scp -r /home/joe/repos/mono/wedding/wedding/build/web/* joe@jonasburster.de:/volumes/mono/wedding/wedding/static
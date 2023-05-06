FROM ubuntu:22.04

RUN apt update && apt install wget -y
RUN wget https://github.com/Spotifyd/spotifyd/releases/download/v0.3.5/spotifyd-linux-full.tar.gz 
RUN tar -xf spotifyd-linux-full.tar.gz 
RUN apt install -y libdbus-1-3 libpulse-dev libasound2

COPY spotifyd.conf /etc/xdg/spotifyd/

# RUN /spotifyd --no-daemon 
CMD ["sleep", "1000"]

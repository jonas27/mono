FROM ubuntu:20.04 as builder
# RUN echo "http://dl-8.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apt update
RUN export DEBIAN_FRONTEND=noninteractive
RUN apt install -y tzdata
RUN apt install -y build-essential python3-dev python3-pip python3-numpy python3-scipy python3-matplotlib
RUN useradd joe
WORKDIR /home/joe
RUN chown -R joe:joe /home/joe 
USER joe
RUN pip3 install scikit-learn 
RUN pip3 install jupyter 
RUN mkdir ./jupyter
WORKDIR /home/joe/jupyter
CMD ["/home/joe/.local/bin/jupyter-notebook" , "--ip=0.0.0.0", "--port=8888" , "--allow-root"]


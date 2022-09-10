FROM golang:1.19-bullseye
WORKDIR /app
RUN apt-get update && apt-get install -y unar build-essential aria2
RUN aria2c -x16 -s16 http://www.jpnsec.org/files/competition2018/data/DB.zip \
    && unar ./DB.zip
RUN wget http://www.jpnsec.org/files/competition2018/data/moon_sop.zip \
    && unar ./moon_sop.zip \
    && wget http://www.jpnsec.org/files/competition2018/data/moon_mop.tgz \
    && tar xvf ./moon_mop.tgz
RUN cd moon_sop && make
RUN cd moon_mop && make
COPY go.mod go.sum wrapper.go ./
RUN ["go", "build"]
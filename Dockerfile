FROM golang:1.5

ADD . /go/src/github.com/fiorix/defacer
WORKDIR /go/src/github.com/fiorix/defacer

RUN echo deb http://httpredir.debian.org/debian jessie main > /etc/apt/sources.list
RUN echo deb http://httpredir.debian.org/debian jessie-updates main >> /etc/apt/sources.list
RUN echo deb http://security.debian.org/ jessie/updates main >> /etc/apt/sources.list

RUN apt-get update
RUN apt-get install -y \
	build-essential \
	libopencv-calib3d2.4 \
	libopencv-contrib2.4 \
	libopencv-core2.4 \
	libopencv-dev \
	libopencv-imgproc2.4 \
	libopencv-ocl2.4 \
	libopencv-stitching2.4 \
	libopencv-superres2.4 \
	libopencv-ts2.4 \
	libopencv-videostab2.4
RUN GO15VENDOREXPERIMENT=1 go install
RUN apt-get autoremove -y --purge build-essential libopencv-dev
RUN apt-get clean

ENTRYPOINT ["/go/bin/defacer"]

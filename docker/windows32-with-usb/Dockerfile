FROM therecipe/qt:windows_32_static-orig
MAINTAINER Dale Farnsworth

RUN apt-get update \
    && apt-get -y install libusb-1.0

RUN cd /usr/lib/mxe/usr/i686-w64-mingw32.static/lib && \
	wget https://farnsworth.org/dale/libusb-1.0.a

#RUN cd /usr/lib/mxe/usr/i686-w64-mingw32.static/include && \
#	wget https://farnsworth.org/dale/libusb.h

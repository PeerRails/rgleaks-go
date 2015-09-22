RgLeaks Go
====================

RgHost Scrapper
---------------------

### Requirements

* libvips >= 7.42.3
* Postgresql >= 9.3
* Go 1.4
* proxy! Maybe Tor, maybe not

### Install LibVips

Ubuntu 14.04 :
```shell
apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y \
  automake build-essential curl \
  gobject-introspection gtk-doc-tools libglib2.0-dev libjpeg-turbo8-dev libpng12-dev libwebp-dev libtiff5-dev libexif-dev libxml2-dev swig libmagickwand-dev

curl -O http://www.vips.ecs.soton.ac.uk/supported/7.42/vips-7.42.3.tar.gz && \
  tar zvxf vips-$LIBVIPS_VERSION.tar.gz && \
  cd vips-$LIBVIPS_VERSION && \
  ./configure --enable-debug=no --without-python --without-orc --without-fftw --without-gsf $1 && \
  make && \
  make install && \
  ldconfig
```

### Run RgLeaks

First set up environment

```shell
export DB_URL=postgres://lenny:123456@localhost/rgleaks-test?sslmode=disable
export IMG_DIR=/var/www/rgleaks/images
```

Build binary or use existing binary

```go
//rgleaks.go
package main

import (
	"github.com/peerrails/rgleaks-go"
	"time"
)

func main() {
	for {
		url := "http://rghost.ru/main"
		rgleaksgo.ScrapeRgHost(url)
		time.Sleep(10 * time.Second)
	}

}


//command-line
go build rgleaks.go
```

Install Tor with privoxy or just pull docker images

```bash
docker pull linuxconfig/instantprivacy
docker run --rm -p 8118:8118 linuxconfig/instantprivacy
```

Check connectiPon with torcheck argument

```bash
http_proxy=127.0.0.1:8118 rgleaks torcheck
```

You should see
>
>
>  Congratulations. This browser is configured to use Tor.
>
>

RUN and WATCH

```bash
http_proxy=127.0.0.1:8118 rgleaks
```

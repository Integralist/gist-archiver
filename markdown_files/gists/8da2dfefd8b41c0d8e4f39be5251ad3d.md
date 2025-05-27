# [Compile NGINX from source, including its dependencies] #nginx #open-source #compile #build

[View original Gist on GitHub](https://gist.github.com/Integralist/8da2dfefd8b41c0d8e4f39be5251ad3d)

**Tags:** #nginx #open-source #compile #build

## 1. Compile NGINX from source, with pre-built dependencies.Dockerfile

```dockerfile
FROM python:3.6.3-slim

# Installing NGINX (open-source):
# https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-open-source/#sources
#
# Locating available versions to build:
# http://nginx.org/download/

# Dependencies:
#
# NGINX: pcre, zlib, openssl (+ libssl-dev), build-essential (for c compiler)
RUN apt-get update && apt-get install -y wget libpcre3 libpcre3-dev zlib1g zlib1g-dev openssl libssl-dev build-essential

# Check Python
RUN python3 --version && pip3 --version
COPY ./requirements.txt /tmp/
RUN pip3 install -r /tmp/requirements.txt

# Compile NGINX (open-source)
WORKDIR /
ENV nginx_version="1.15.2"
RUN wget https://nginx.org/download/nginx-${nginx_version}.tar.gz
RUN tar zxf nginx-${nginx_version}.tar.gz
WORKDIR nginx-${nginx_version}
RUN ./configure \
      --user=www-data \
      --group=www-data \
      --with-pcre-jit \
      --with-debug \
      --with-http_ssl_module \
      --with-http_stub_status_module \
      --with-http_realip_module

# Possible future configuration addition:
# https://www.nginx.com/blog/thread-pools-boost-performance-9x/

RUN make
RUN make install
ENV PATH="/nginx-${nginx_version}/objs:$PATH"

WORKDIR /
COPY . /app

CMD python /app/template.py nginx.conf.j2 /nginx.conf && nginx -c /nginx.conf
```

## 2. Compile NGINX from source, including its dependencies.Dockerfile

```dockerfile
FROM python:3.6.3-slim

# Dependencies
# NGINX: pcre, zlib, openssl (+ libssl-dev), build-essential (for c compiler)
# RUN apt-get update && apt-get install -y wget libpcre3 libpcre3-dev zlib1g zlib1g-dev openssl libssl-dev build-essential
RUN apt-get update && apt-get install -y wget build-essential

# Check Python
RUN python3 --version && pip3 --version
COPY ./requirements.txt /tmp/
RUN pip3 install -r /tmp/requirements.txt

# Installing NGINX (open-source)
# https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-open-source/#sources

# PCRE
WORKDIR /
ENV pcre_version="8.42"
RUN wget ftp://ftp.csx.cam.ac.uk/pub/software/programming/pcre/pcre-${pcre_version}.tar.gz
RUN tar -zxf pcre-${pcre_version}.tar.gz
WORKDIR pcre-${pcre_version}
RUN ./configure
RUN make
RUN make install

# zlib
WORKDIR /
ENV zlib_version="1.2.11"
RUN wget http://zlib.net/zlib-${zlib_version}.tar.gz
RUN tar -zxf zlib-${zlib_version}.tar.gz
WORKDIR zlib-${zlib_version}
RUN ./configure
RUN make
RUN make install

# OpenSSL
WORKDIR /
ENV openssl_version="1.0.2o"
RUN wget http://www.openssl.org/source/openssl-${openssl_version}.tar.gz
RUN tar -zxf openssl-${openssl_version}.tar.gz
WORKDIR openssl-${openssl_version}
# RUN ./Configure LIST
RUN ./Configure linux-x86_64 --prefix=/usr
RUN make
RUN make install

# NGINX (open-source) http://nginx.org/download/
WORKDIR /
ENV nginx_version="1.15.2"
RUN wget https://nginx.org/download/nginx-${nginx_version}.tar.gz
RUN tar zxf nginx-${nginx_version}.tar.gz
WORKDIR nginx-${nginx_version}
RUN ./configure \
      --user=www-data \
      --group=www-data \
      --with-pcre-jit \
      --with-debug \
      --with-http_ssl_module \
      --with-http_stub_status_module \
      --with-http_realip_module

# Possible future configuration addition:
# https://www.nginx.com/blog/thread-pools-boost-performance-9x/

RUN make
RUN make install
ENV PATH="/nginx-${nginx_version}/objs:$PATH"

WORKDIR /
COPY . /app

CMD nginx -v
```


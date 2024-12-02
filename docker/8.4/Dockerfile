FROM ubuntu:24.04

LABEL maintainer="Taylor Otwell"

ARG WWWGROUP
ARG NODE_VERSION=22
ARG MYSQL_CLIENT="mysql-client"
ARG POSTGRES_VERSION=17

WORKDIR /var/www/html

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC
ENV SUPERVISOR_PHP_COMMAND="/usr/bin/php -d variables_order=EGPCS /var/www/html/artisan serve --host=0.0.0.0 --port=80"
ENV SUPERVISOR_PHP_USER="sail"

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN echo "Acquire::http::Pipeline-Depth 0;" > /etc/apt/apt.conf.d/99custom && \
    echo "Acquire::http::No-Cache true;" >> /etc/apt/apt.conf.d/99custom && \
    echo "Acquire::BrokenProxy    true;" >> /etc/apt/apt.conf.d/99custom

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y \
    git python3 python3-pip unzip plocate build-essential ninja-build \
    libssl-dev libgmp-dev zlib1g-dev openssl libpcre3-dev php-dev cmake

# Build libuv
RUN git clone --depth 1 -b v1.46.0 https://github.com/libuv/libuv.git /tmp/libuv \
    && cd /tmp/libuv \
    && mkdir build \
    && cd build \
    && cmake -DBUILD_TESTING=OFF -DBUILD_BENCHMARKS=OFF -DLIBUV_BUILD_SHARED=ON -DCMAKE_C_FLAGS="-fPIC" -DCMAKE_BUILD_TYPE="RelWithInfo" -G Ninja .. \
    && ninja install \
    && rm -rf /tmp/libuv

# Build ScyllaDB CPP driver
RUN git clone --depth 1 https://github.com/scylladb/cpp-driver.git /tmp/scyladb-driver \
    && cd /tmp/scyladb-driver \
    && mkdir build \
    && cd build \
    && cmake -DCASS_CPP_STANDARD=17 -DCASS_BUILD_STATIC=ON -DCASS_BUILD_SHARED=ON -DCASS_USE_STD_ATOMIC=ON -DCASS_USE_TIMERFD=ON -DCASS_USE_LIBSSH2=ON -DCASS_USE_ZLIB=ON -DCMAKE_C_FLAGS="-fPIC" -DCMAKE_CXX_FLAGS="-fPIC -Wno-error=redundant-move" -DCMAKE_BUILD_TYPE="RelWithInfo" -G Ninja .. \
    && ninja install \
    && rm -rf /tmp/scyladb-driver

RUN mkdir -p /etc/apt/keyrings \
    && apt-get install -y gnupg gosu curl ca-certificates zip unzip git supervisor sqlite3 libcap2-bin libpng-dev python3 dnsutils librsvg2-bin fswatch ffmpeg nano  \
    && curl -sS 'https://keyserver.ubuntu.com/pks/lookup?op=get&search=0x14aa40ec0831756756d7f66c4f4ea0aae5267a6c' | gpg --dearmor | tee /etc/apt/keyrings/ppa_ondrej_php.gpg > /dev/null \
    && echo "deb [signed-by=/etc/apt/keyrings/ppa_ondrej_php.gpg] https://ppa.launchpadcontent.net/ondrej/php/ubuntu noble main" > /etc/apt/sources.list.d/ppa_ondrej_php.list \
    && apt-get update \
    && apt-get install -y php8.4-cli php8.4-dev \
       php8.4-pgsql php8.4-sqlite3 php8.4-gd \
       php8.4-curl php8.4-mongodb \
       php8.4-imap php8.4-mysql php8.4-mbstring \
       php8.4-xml php8.4-zip php8.4-bcmath php8.4-soap \
       php8.4-intl php8.4-readline \
       php8.4-ldap \
       php8.4-msgpack php8.4-igbinary php8.4-redis \
#       php8.4-swoole \
       php8.4-memcached php8.4-pcov php8.4-imagick php8.4-xdebug \
    && curl -sLS https://getcomposer.org/installer | php -- --install-dir=/usr/bin/ --filename=composer \
    && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_VERSION.x nodistro main" > /etc/apt/sources.list.d/nodesource.list \
    && apt-get update \
    && apt-get install -y nodejs \
    && npm install -g npm \
    && npm install -g pnpm \
    && npm install -g bun \
    && curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | gpg --dearmor | tee /etc/apt/keyrings/yarn.gpg >/dev/null \
    && echo "deb [signed-by=/etc/apt/keyrings/yarn.gpg] https://dl.yarnpkg.com/debian/ stable main" > /etc/apt/sources.list.d/yarn.list \
    && curl -sS https://www.postgresql.org/media/keys/ACCC4CF8.asc | gpg --dearmor | tee /etc/apt/keyrings/pgdg.gpg >/dev/null \
    && echo "deb [signed-by=/etc/apt/keyrings/pgdg.gpg] http://apt.postgresql.org/pub/repos/apt noble-pgdg main" > /etc/apt/sources.list.d/pgdg.list \
    && apt-get update \
    && apt-get install -y yarn \
    && apt-get install -y $MYSQL_CLIENT \
    && apt-get install -y postgresql-client-$POSTGRES_VERSION \
    && apt-get -y autoremove \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN git clone https://github.com/he4rt/scylladb-php-driver.git /tmp/scylladb-php-driver \
    && cd /tmp/scylladb-php-driver \
    && sed -i '/find_package(Sanitizers REQUIRED)/d' CMakeLists.txt \
    && rm -rf build \
    && cmake -S . -B build \
        -DPHP_SCYLLADB_ENABLE_SANITIZERS=OFF \
        -DPHP_SCYLLADB_ENABLE_COVERAGE=OFF \
        -DPHP_SCYLLADB_OPTIMISE_FOR_CURRENT_MACHINE=ON \
        -DPHP_SCYLLADB_LIBUV_STATIC=OFF \
        -DPHP_SCYLLADB_LIBUV_FROM_SRC=OFF \
        -DPHP_SCYLLADB_LIBSCYLLADB_STATIC=OFF \
        -DPHP_SCYLLADB_LIBSCYLLADB_FROM_SRC=OFF \
        -DPHP_SCYLLADB_USE_LIBCASSANDRA=OFF \
        -DPHP_SCYLLADB_LIBCASSANDRA_STATIC=OFF \
        -DPHP_SCYLLADB_LIBCASSANDRA_FROM_SRC=OFF \
    && cmake --build build --target install

RUN EXTENSION_DIR=$(php-config8.4 --extension-dir) \
    && cp /tmp/scylladb-php-driver/build/cassandra.so "$EXTENSION_DIR/cassandra.so" \
    && cp /tmp/scylladb-php-driver/cassandra.ini /etc/php/8.4/cli/conf.d/10-cassandra.ini

RUN setcap "cap_net_bind_service=+ep" /usr/bin/php8.4

RUN userdel -r ubuntu
RUN groupadd --force -g $WWWGROUP sail
RUN useradd -ms /bin/bash --no-user-group -g $WWWGROUP -u 1337 sail

COPY start-container /usr/local/bin/start-container
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY php.ini /etc/php/8.4/cli/conf.d/99-sail.ini
RUN chmod +x /usr/local/bin/start-container

EXPOSE 80/tcp

ENTRYPOINT ["start-container"]

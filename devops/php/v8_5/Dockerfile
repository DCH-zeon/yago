FROM php:8.5.4-fpm-alpine3.23

# Встановлення системних залежностей та PHP-розширень
RUN apk add --no-cache \
    $PHPIZE_DEPS \
    git \
    unzip \
    curl \
    bash \
    su-exec \
    libzip-dev \
    libpng-dev \
    libjpeg-turbo-dev \
    libwebp-dev \
    libxpm-dev \
    libxml2-dev \
    imagemagick \
    libtool \
    pkgconf \
    imagemagick-dev \
    icu-dev \
    oniguruma-dev \
    postgresql-dev \
    rabbitmq-c-dev \
    util-linux-dev \
    linux-headers \
    nodejs \
    npm \
    && docker-php-ext-configure gd \
    --with-jpeg \
    --with-webp \
    --with-xpm

RUN pecl install redis imagick xdebug \
    && docker-php-ext-enable redis imagick xdebug \
    && docker-php-ext-install -j$(nproc) \
        pdo_pgsql \
        pdo_mysql \
        mbstring \
        zip \
        exif \
        pcntl \
        bcmath \
        gd \
        intl

# Встановлення Composer
COPY --from=composer:2 /usr/bin/composer /usr/bin/composer
# Базова конфігурація Xdebug
COPY ./devops/php/xdebug.ini /usr/local/etc/php/conf.d/xdebug.ini
# Конфігурація php
COPY ./devops/php/www.conf /usr/local/etc/php-fpm.d/www.conf

# Вказуємо sys_temp_dir для PHP
RUN echo "sys_temp_dir=/var/www/admin/storage/temp" > /usr/local/etc/php/conf.d/tempdir.ini \
    && echo "upload_tmp_dir=/var/www/admin/storage/temp" >> /usr/local/etc/php/conf.d/tempdir.ini

# Робоча директорія
WORKDIR /var/www/admin

COPY ./devops/php/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
# Запуск PHP-FPM
CMD ["php-fpm", "-F"]

#!/bin/sh
set -e

cd /var/www/admin

# Створення необхідних директорій
mkdir -p storage/temp storage/framework/views storage/framework/cache storage/framework/sessions

# Налаштування прав на весь проект для www-data
echo "Setting permissions on the entire project..."
chown -R www-data:www-data /var/www/admin
chmod -R 775 /var/www/admin

# Встановлення залежностей від імені www-data
if [ ! -d node_modules ]; then
  echo "Installing npm dependencies..."
  su-exec www-data npm install
fi

if [ ! -d vendor ]; then
  echo "Installing composer dependencies..."
  su-exec www-data composer install --no-interaction --prefer-dist
fi

# Генерація типів для Wayfinder (потрібно для запуску Vite)
echo "Generating Wayfinder types..."
su-exec www-data php artisan wayfinder:generate --with-form 2>/dev/null || echo "Wayfinder generation skipped"

echo "Starting Vite dev server..."
su-exec www-data npm run dev >/tmp/vite.log 2>&1 &

echo "Starting PHP-FPM..."
exec php-fpm -F
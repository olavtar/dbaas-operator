# Stage 1: MongoDB
FROM mongo:7.0.10-rc0-jammy AS mongo_stage

# Stage 2: Nginx
#FROM nginx:stable-perl AS nginx_stage

# Stage 3: Python
#FROM python:3.9 AS python_stage

# Stage 4: MySQL
#FROM mysql:8.0.37-debian AS mysql_stage

# Stage 5: PostgreSQL
#FROM postgres:bullseye AS postgres_stage

# Stage 6: PHP
#FROM php:zts-bullseye AS php_stage
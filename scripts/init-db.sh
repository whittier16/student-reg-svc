#!/bin/bash
/usr/bin/mysqld_safe --skip-grant-tables &
sleep 5
mysql -u root -e "CREATE DATABASE stdnt_reg"
mysql -u root mydb < /db/migration/dump.sql

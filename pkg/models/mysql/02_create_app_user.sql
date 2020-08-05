USE snippetbox;

CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE ON snippetbox.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY 'insecure';
update mysql.user set host = '%' where user='web';
commit;
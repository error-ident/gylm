# Set Up Guide

## Yandex

* Open [YM settings](https://music.yandex.ru/settings/account)
* Import music collections -> Link Last.fm account 

## Last.fm

* [Create API account](https://www.last.fm/api/account/create)
* Save it

## Server

* Clone repo: https://github.com/Innsmouth-trip/gylm

* Create .env file in repo folder

```
KEY=9b69**************************** | lastfm api key 
SECRET=de01**************************** | lastfm shared secret      
USERNAME=error_ident | lastfm username
LIMIT= 10 | the number of results to fetch per page 
YANDEX_URL=https://music.yandex.ru/users/Invisible-sleeper/artists | YM profile url
```

* Buuld for linux arch: GOOS=linux GOARCH=amd64 go build -o gylm

* Add gylm.conf to /etc/nginx/sites-available

```
server {
    server_name 92.119.90.17;
    location / {
        proxy_pass http://localhost:1984/yandex;
	expires    0;
    }
}
```

* Create symlink: sudo ln -s /etc/nginx/sites-available/gylm.conf /etc/nginx/sites-enabled/

* Restart Nginx: sudo systemctl restart nginx


## ReadMe

You can now use the following in your readme:

```[![yandex](http://92.119.90.17)](https://music.yandex.ru/users/Invisible-sleeper/albums)```
#
# go-everywhere
Proxy single page in whitelist with go.

- [x] http/https page
- [x] http/https git repository

## Deploy
Use source

    git clone https://github.com/MXWXZ/go-everywhere.git
    go run main.go

Or [release binary](https://github.com/MXWXZ/go-everywhere/releases)

Or Docker

    touch whitelist.txt
    chmod 777 whitelist.txt  # otherwise you can't sync the file inside docker
    sudo docker run -v $(pwd)/whitelist.txt:/app/whitelist.txt -p 8080:8080 -d imwxz/go-everywhere

View `http://localhost:8080/` for running check, view `http://localhost:8080/https://www.google.com/` for proxy.

## Config
Edit `whitelist.txt` to give access(do NOT change the filename!)

Each line in `whitelist.txt` indicate a website config, in regular expression format, like

    ^https://www\.google\.com.*

If you do not need whitelist, just insert this in the file

    .*

## Router
- `/`: running check
- `/xxx`: proxy url, **http/https prefix is a must**, like `/https://www.google.com`
- `/reload`: reload whitelist without restart server

## Environment variable
- `GIN_MODE=release` for release mode(when you deploy with source code)
- `PORT=8080` specify the listen port, default 8080

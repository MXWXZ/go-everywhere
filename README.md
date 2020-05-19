# go-everywhere
Proxy single page in whitelist with go.

## Deploy
Use source

    git clone https://github.com/MXWXZ/go-everywhere.git
    go run main.go

Or [release binary](https://github.com/MXWXZ/go-everywhere/releases)

## Config
Edit `whitelist.txt` to give access(do NOT change the filename!)

Each line in `whitelist.txt` indicate a website config, in regular expression format, like

    ^https://www\.google\.com.*

If you do not need whitelist, just insert this in the file

    .*

## Router
- `/`: running check
- `/?url=xxx`: proxy url, **http/https prefix is a must**, like `/?url=https://google.com`
- `/reload`: reload whitelist without restart server

## Environment variable
- `GIN_MODE=release` for release mode(when you deploy with source code)
- `PORT=8080` specify the listen port, default 8080

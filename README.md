# proxycmd

Connect SSH server with HTTP proxy

## Usage
Normally HTTP proxy with "Connect" only allows 443(HTTP) and 563(SNEWS), so make sure you've configure
your SSH server to use the aforementioned ports.

```bash
make
ssh -o ProxyCommand="./proxycmd -proxy http://http.proxy.addr:port -ssh-server ssh.server.addr:port" ssh.server.addr
```

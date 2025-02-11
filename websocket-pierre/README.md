```sh
curl -i --header "Upgrade: websocket" \
--header "Connection: Upgrade" \
--header "Sec-WebSocket-Key: YSBzYW1wbGUgMTYgYnl0ZQ==" \
--header "Sec-Websocket-Version: 13" \
localhost:8080
```


```sh
MacBook-Air-de-pierre:Websocket pierrecaboor$ curl -i --header "Upgrade: websocket" \
> --header "Connection: Upgrade" \
> --header "Sec-WebSocket-Key: YSBzYW1wbGUgMTYgYnl0ZQ==" \
> --header "Sec-Websocket-Version: 13" \
> localhost:8080
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: zkTGI6zVrOIDXiC4vnn1Rf37YFw=

curl: (52) Empty reply from server
```

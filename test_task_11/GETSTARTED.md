# Get Started

1. Test
- go test -v ./...

2. API doc
- to send command use - http://serveraddress/cmd
```
{
    "name":"SET",  /* SET, GET, REMOVE, KEYS */
    "type":"str",  /* str, slice, map */
    "key":"key",   /* key for stored value */
    "ttl":3600,    /* key ttl, second */
    "str":"value",  /* if type str */
    "slice":["value", "value"], /* if type slice */
    "map": {"key1":"value1","key2":"value2"}, /* if type map */
    "idx_key":"key1" /* for get slice, map by index/key */
}
```
3. Run
- Server
```
cd mcache-srv
go build
./mcache-srv

```
- Client
```
cd mcache-cln
go build
./mcache-cln
```

4. Example, client commands:
- set string value
```
enter command: set str key 3600 value
```
- keys list
```
enter command: keys
[key]
```
- set slice value
```
enter command: set slice key2 3600 q w e r ty
```
- keys list
```
enter command: keys
[key key2]
```
- get value
```
enter command: get key2
[q w e r ty]
```
- set map
```
enter command: set map m 1200 q qwerty a asdf
enter command: get m
map[a:asdf q:qwerty]
```
- get map item by key
```
enter command: get m a
asdf

```
- remove
```
enter command: remove m
enter command: keys
[key2 key]
```

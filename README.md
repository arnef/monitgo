# montigo

## start node
```sh
monitgo node --host 127.0.0.1 --port 5000
###
# OPTIONS:
#   --host value  (default: "127.0.0.1")
#   --port value  (default: 5000)
###
```

## start watcher

```sh
monitgo watch
###
# OPTIONS:
#   --interval value, -n value  (default: 60)
#   --no-bot                    (default: false)
###
```

### config.yml
```yaml
nodes:
  - name: x390
    host: 10.10.10.164
    port: 5000

telegram:
  token: supersecuretelegrambotapitoken
```
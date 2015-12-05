# meow

Run:

```bash
docker run -p 8080:8080 -d fiorix/defacer
curl localhost:8080/api/v1/deface?url=http://bit.ly/1gBahPH > faces.jpg
curl localhost:8080/api/v1/metrics | grep deface
```

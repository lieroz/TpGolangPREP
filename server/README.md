# Highload course

## HTTP server

<a name="drun"></a>
### Server start options

```bash
  $ go run *.go -p=:8080 -c=4 -wr=/var/www/html
``` 

* -p Port to run server on
* -c Number of cores to utilize
* -wr Webroot directory, all static lies here


<a name="drun"></a>
### Docker run

```bash
  $ docker build -t [NAME] .
  $ docker run --publish 8080:8080 --name [CONTAINER_NAME] [--rm [NAME]]
```

<a name="htest"></a>
### Http test suite

[Tests repo](https://github.com/init/http-test-suite)  

**All tests passed**  

```bash
  $ ./httptest.py
```

<a name="drun"></a>
### Run load on server

```bash
  $ ab -n 100000 -c 100 127.0.0.1:8080/httptest/dir2/page.html
``` 

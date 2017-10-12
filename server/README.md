# Technopark 3 semester

## Highload HW 1: HTTP server

<a name="drun"></a>
### Server start options

```bash
  $ go run *.go -p=:80 -c=0 -wr=/var/www/html -w=4
``` 

* -p Port to run server on
* -c Number of cores to utilize
* -wr Webroot directory, all static lies here
* -w Worker count

<a name="drun"></a>
### Docker run

```bash
  $ docker build -t [NAME] .
  $ docker run --publish 80:80 --name [CONTAINER_NAME] [--rm [NAME] or -t [NAME]] 
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
  $ ab -n 100000 -c 100 127.0.0.1:80/httptest/dir2/page.html
``` 

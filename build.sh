docker build -t webcrawler:v1 .
rm -rf github.com golang.org
docker container stop $(docker container ls)
# docker run  --log-driver=fluentd --log-opt tag=docker.{{.ID}} -v /etc/localtime:/etc/localtime -d -p 80:80 -p  8088:8088  webcrawler:v1
docker run  -d -p 80:80 -p  8088:8088  webcrawler:v1

cd /home/go/src/webcrawler
git pull
cd ..
cp -r github.com webcrawler
cp -r golang.org webcrawler
cd webcrawler
docker build -t webcrawler:v1 .
rm -rf github.com golang.org
docker container stop $(docker container ls)
# docker run  --log-driver=fluentd --log-opt tag=docker.{{.ID}} -v /etc/localtime:/etc/localtime -d -p 80:80 -p  8088:8088  webcrawler:v1
docker run  -d -p 80:80 -p  8088:8088  webcrawler:v1

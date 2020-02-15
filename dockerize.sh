echo "### Building image"
docker build -t ascii .
echo
echo "### Running image"
docker run -d -p 8080:8080 --name web ascii
echo
echo "### Images list"
docker images
echo
echo "### Containers list"
docker container ls
echo
echo "### Commandline in container"
docker exec -it web /bin/sh
docker build \
  -t "docker.stdlib.in/mango:6" \
  -f ./mango6.dockerfile .
docker push "docker.stdlib.in/mango:6"

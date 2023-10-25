TAG1="docker.stdlib.in/mango:5"

docker build \
  -t "$TAG1" \
  -f ./mango5.dockerfile .
docker push "$TAG1"

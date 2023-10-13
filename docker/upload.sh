TAG1="docker.dev-friends.com/mango:5"

docker build \
  -t "$TAG1" \
  -f ./mango5.dockerfile .
docker push "$TAG1"

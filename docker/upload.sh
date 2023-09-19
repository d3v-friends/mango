TAG1="docker.dev-friends.com/mango:latest"
TAG2="docker.dev-friends.com/mango:5"

docker build \
  -t "$TAG1" \
  -t "$TAG2" \
  -f ./mango5.dockerfile .
docker push "$TAG1"
docker push "$TAG2"

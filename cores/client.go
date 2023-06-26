package cores

import "go.mongodb.org/mongo-driver/mongo"

type (
	Client struct {
		client *mongo.Client
	}
)

func (x *Client) Client() *mongo.Client {
	return x.client
}

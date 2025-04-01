package mgOp

import "go.mongodb.org/mongo-driver/bson"

const (
	Match        = "$match"
	Limit        = "$limit"
	Lookup       = "$lookup"
	ElemMatch    = "$elemMatch"
	ArrayElemAt  = "$arrayElemAt"
	Push         = "$push"
	Pull         = "$pull"
	PullAll      = "$pullAll"
	Set          = "$set"
	Unset        = "$unset"
	Inc          = "$inc"
	First        = "$first"
	Lt           = "$lt"
	Lte          = "$lte"
	Gt           = "$gt"
	Gte          = "$gte"
	Eq           = "$eq"
	Ne           = "$ne"
	In           = "$in"
	Regex        = "$regex"
	Exists       = "$exists"
	Each         = "$each"
	Or           = "$or"
	Sort         = "$sort"
	Project      = "$project"
	AddFields    = "$addFields"
	Map          = "$map"
	MergeObjects = "$mergeObjects"
	Filter       = "$filter"
)

// LookupSingle
// https://www.mongodb.com/docs/manual/reference/operator/aggregation/lookup/#equality-match-with-a-single-join-condition
type LookupSingle struct {
	From         string `bson:"from,omitempty"`
	As           string `bson:"as,omitempty"`
	LocalField   string `bson:"localField,omitempty"`
	ForeignField string `bson:"foreignField,omitempty"`
}

// LookUpSubquery
// https://www.mongodb.com/docs/manual/reference/operator/aggregation/lookup/#join-conditions-and-subqueries-on-a-joined-collection
type LookUpSubquery struct {
	From     string `bson:"from,omitempty"`
	As       string `bson:"as,omitempty"`
	Let      bson.M `bson:"let,omitempty"`
	Pipeline bson.A `bson:"pipeline,omitempty"`
}

// LookupConcise
// concise: 간결한, 축약된
// https://www.mongodb.com/docs/manual/reference/operator/aggregation/lookup/#correlated-subqueries-using-concise-syntax
type LookupConcise struct {
	From         string `bson:"from,omitempty"`
	As           string `bson:"as,omitempty"`
	LocalField   string `bson:"localField"`
	ForeignField string `bson:"foreignField"`
	Let          bson.M `bson:"let,omitempty"`
	Pipeline     bson.A `bson:"pipeline,omitempty"`
}

// FilterArgs
// https://www.mongodb.com/docs/manual/reference/operator/aggregation/filter/#syntax
type FilterArgs struct {
	Input *string `bson:"input,omitempty"`
	As    *string `bson:"as,omitempty"`
	Cond  bson.M  `bson:"cond,omitempty"`
	Limit *uint64 `bson:"limit,omitempty"`
}

// MapArgs
// https://www.mongodb.com/docs/manual/reference/operator/aggregation/map/#syntax
type MapArgs struct {
	Input string  `bson:"input"`
	As    *string `bson:"as,omitempty"`
	In    bson.M  `bson:"in,omitempty"`
}

package mvars

const (
	// FInTrx transaction 을 위한 모델마다 가지고 있어야 하는 필드
	FInTrx     = "inTrx"
	FID        = "_id"
	FCreatedAt = "createdAt"
	FUpdatedAt = "updatedAt"
	FDeletedAt = "deletedAt"
	FErrors    = "errors"
)

const (
	OSet         = "$set"
	OPush        = "$push"
	OExist       = "$exist"
	OOr          = "$or"
	OUnset       = "$unset"
	OPull        = "$pull"
	ONE          = "$ne"
	OMatch       = "$match"
	OLookUp      = "$lookup"
	OReplaceRoot = "$replaceRoot"
	OType        = "$type"
	ONot         = "$not"
	OLimit       = "$limit"
	OSkip        = "$skip"
	OArrayElemAt = "$arrayElemAt"
	OAddField    = "$addFields"
	OAll         = "$all"
	OAddToSet    = "$addToSet"
)

const (
	VArrayType  = "array"
	VObjectType = "object"
)

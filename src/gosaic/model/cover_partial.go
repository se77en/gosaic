package model

import "image"

type CoverPartial struct {
	Id       int64 `db:"id"`
	CoverId  int64 `db:"cover_id"`
	AspectId int64 `db:"aspect_id"`
	X1       int64 `db:"x1"`
	Y1       int64 `db:"y1"`
	X2       int64 `db:"x2"`
	Y2       int64 `db:"y2"`
}

func (cp *CoverPartial) Rectangle() image.Rectangle {
	return image.Rect(int(cp.X1), int(cp.Y1), int(cp.X2), int(cp.Y2))
}

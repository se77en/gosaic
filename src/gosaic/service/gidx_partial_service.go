package service

import (
	"fmt"
	"gosaic/model"
	"gosaic/util"
	"sync"

	"gopkg.in/gorp.v1"
)

type GidxPartialService interface {
	Service
	Insert(*model.GidxPartial) error
	Update(*model.GidxPartial) error
	Delete(*model.GidxPartial) error
	Get(int64) (*model.GidxPartial, error)
	GetOneBy(string, interface{}) (*model.GidxPartial, error)
	ExistsBy(string, interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, interface{}) (int64, error)
	Find(*model.Gidx, *model.Aspect) (*model.GidxPartial, error)
	Create(*model.Gidx, *model.Aspect) (*model.GidxPartial, error)
	FindOrCreate(*model.Gidx, *model.Aspect) (*model.GidxPartial, error)
	FindMissing(*model.Aspect, string, int, int) ([]*model.Gidx, error)
}

type gidxPartialServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewGidxPartialService(dbMap *gorp.DbMap) GidxPartialService {
	return &gidxPartialServiceImpl{dbMap: dbMap}
}

func (s *gidxPartialServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *gidxPartialServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.GidxPartial{}, "gidx_partials").SetKeys(true, "id")
	return nil
}

func (s *gidxPartialServiceImpl) Insert(gidx_partial *model.GidxPartial) error {
	err := gidx_partial.EncodePixels()
	if err != nil {
		return err
	}
	return s.DbMap().Insert(gidx_partial)
}

func (s *gidxPartialServiceImpl) Update(gidx_partial *model.GidxPartial) error {
	err := gidx_partial.EncodePixels()
	if err != nil {
		return err
	}
	_, err = s.DbMap().Update(gidx_partial)
	return err
}

func (s *gidxPartialServiceImpl) Delete(gidx_partial *model.GidxPartial) error {
	_, err := s.DbMap().Delete(gidx_partial)
	return err
}

func (s *gidxPartialServiceImpl) Get(id int64) (*model.GidxPartial, error) {
	gidx_partial, err := s.DbMap().Get(model.GidxPartial{}, id)
	if err != nil {
		return nil, err
	} else if gidx_partial == nil {
		return nil, nil
	}

	gp, ok := gidx_partial.(*model.GidxPartial)
	if !ok {
		return nil, fmt.Errorf("Received struct is not a GidxPartial")
	}

	err = gp.DecodeData()
	if err != nil {
		return nil, err
	}

	return gp, nil
}

func (s *gidxPartialServiceImpl) GetOneBy(column string, value interface{}) (*model.GidxPartial, error) {
	var gidx_partial model.GidxPartial
	err := s.DbMap().SelectOne(&gidx_partial, "select * from gidx_partials where "+column+" = ? limit 1", value)
	return &gidx_partial, err
}

func (s *gidxPartialServiceImpl) ExistsBy(column string, value interface{}) (bool, error) {
	count, err := s.DbMap().SelectInt("select 1 from gidx_partials where "+column+" = ? limit 1", value)
	return count == 1, err
}

func (s *gidxPartialServiceImpl) Count() (int64, error) {
	return s.DbMap().SelectInt("select count(*) from gidx_partials")
}

func (s *gidxPartialServiceImpl) CountBy(column string, value interface{}) (int64, error) {
	return s.DbMap().SelectInt("select count(*) from gidx_partials where "+column+" = ?", value)
}

func (s *gidxPartialServiceImpl) doFind(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	p := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
	}

	err := s.DbMap().SelectOne(&p, "select * from gidx_partials where gidx_id = ? and aspect_id = ?", p.GidxId, p.AspectId)
	if err != nil {
		return nil, err
	}

	err = p.DecodeData()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *gidxPartialServiceImpl) Find(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doFind(gidx, aspect)
}

func (s *gidxPartialServiceImpl) doCreate(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	p := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
	}

	pixels, err := util.GetAspectLab(gidx, aspect)
	if err != nil {
		return nil, err
	}
	p.Pixels = pixels

	err = p.EncodePixels()
	if err != nil {
		return nil, err
	}

	err = s.Insert(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *gidxPartialServiceImpl) Create(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doCreate(gidx, aspect)
}

func (s *gidxPartialServiceImpl) FindOrCreate(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	p, err := s.doFind(gidx, aspect)
	if err == nil {
		return p, nil
	}

	// or create
	return s.doCreate(gidx, aspect)
}

func (s *gidxPartialServiceImpl) FindMissing(aspect *model.Aspect, order string, limit, offset int) ([]*model.Gidx, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := `
select * from gidx
where not exists (
	select 1 from gidx_partials
	where gidx_partials.gidx_id = gidx.id
	and gidx_partials.aspect_id = ?
)
order by ?
limit ?
offset ?
`

	var gidxs []*model.Gidx
	_, err := s.dbMap.Select(&gidxs, sql, aspect.Id, order, limit, offset)

	return gidxs, err
}

package service

import (
	"testing"

	"gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupPartialComparisonServiceTest() (PartialComparisonService, error) {
	dbMap, err := getTestDbMap()
	if err != nil {
		return nil, err
	}

	gidxService, err := getTestGidxService(dbMap)
	if err != nil {
		return nil, err
	}

	gidxPartialService, err := getTestGidxPartialService(dbMap)
	if err != nil {
		return nil, err
	}

	aspectService, err := getTestAspectService(dbMap)
	if err != nil {
		return nil, err
	}

	coverService, err := getTestCoverService(dbMap)
	if err != nil {
		return nil, err
	}

	coverPartialService, err := getTestCoverPartialService(dbMap)
	if err != nil {
		return nil, err
	}

	macroService, err := getTestMacroService(dbMap)
	if err != nil {
		return nil, err
	}

	macroPartialService, err := getTestMacroPartialService(dbMap)
	if err != nil {
		return nil, err
	}

	partialComparisonService, err := getTestPartialComparisonService(dbMap)
	if err != nil {
		return nil, err
	}

	aspect = model.Aspect{Columns: 239, Rows: 170}
	err = aspectService.Insert(&aspect)
	if err != nil {
		return nil, err
	}

	gidx = model.Gidx{
		AspectId:    aspect.Id,
		Path:        "testdata/shaq_bill.jpg",
		Md5sum:      "394c43174e42e043e7b9049e1bb10a39",
		Width:       uint(478),
		Height:      uint(340),
		Orientation: 1,
	}
	err = gidxService.Insert(&gidx)
	if err != nil {
		return nil, err
	}

	gidx2 := model.Gidx{
		AspectId:    aspect.Id,
		Path:        "testdata/eagle.jpg",
		Md5sum:      "5a19b84638fc471d8ec4167ea4e659fb",
		Width:       uint(512),
		Height:      uint(364),
		Orientation: 1,
	}
	err = gidxService.Insert(&gidx2)
	if err != nil {
		return nil, err
	}

	cover = model.Cover{Name: "test1", AspectId: aspect.Id, Width: 1, Height: 1}
	err = coverService.Insert(&cover)
	if err != nil {
		return nil, err
	}

	gp, err := gidxPartialService.FindOrCreate(&gidx, &aspect)
	if err != nil {
		return nil, err
	}
	gidxPartial = *gp

	_, err = gidxPartialService.FindOrCreate(&gidx2, &aspect)
	if err != nil {
		return nil, err
	}

	coverPartials := make([]model.CoverPartial, 6)
	for i := 0; i < 6; i++ {
		cp := model.CoverPartial{
			CoverId:  cover.Id,
			AspectId: aspect.Id,
			X1:       int64(i),
			Y1:       int64(i),
			X2:       int64(i + 1),
			Y2:       int64(i + 1),
		}
		err = coverPartialService.Insert(&cp)
		if err != nil {
			return nil, err
		}
		if i == 6 {
			coverPartial = cp
		} else {
			coverPartials[i] = cp
		}
	}

	macro = model.Macro{
		CoverId:     cover.Id,
		AspectId:    aspect.Id,
		Path:        "testdata/matterhorn.jpg",
		Md5sum:      "fcaadee574094a3ae04c6badbbb9ee5e",
		Width:       uint(696),
		Height:      uint(1024),
		Orientation: 1,
	}
	err = macroService.Insert(&macro)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 5; i++ {
		mp, err := macroPartialService.FindOrCreate(&macro, &coverPartials[i])
		if err != nil {
			return nil, err
		}
		if i == 0 {
			macroPartial = *mp
		}
	}

	return partialComparisonService, nil
}

func TestPartialComparisonServiceInsert(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err = partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	if pc.Id == int64(0) {
		t.Fatalf("Inserted partial comparison id not set")
	}

	pc2, err := partialComparisonService.Get(pc.Id)
	if err != nil {
		t.Fatalf("Error getting inserted partial comparison: %s\n", err.Error())
	} else if pc2 == nil {
		t.Fatalf("Partial comparison not inserted\n")
	}

	if pc.Id != pc2.Id ||
		pc.MacroPartialId != pc2.MacroPartialId ||
		pc.GidxPartialId != pc2.GidxPartialId ||
		pc.Dist != pc2.Dist {
		t.Fatal("Inserted macro partial data does not match")
	}
}

func TestPartialComparisonServiceUpdate(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err = partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	pc.Dist = 0.24
	err = partialComparisonService.Update(&pc)
	if err != nil {
		t.Fatalf("Error updating partial comparison: %s\n", err.Error())
	}

	pc2, err := partialComparisonService.Get(pc.Id)
	if err != nil {
		t.Fatalf("Error getting updated partial comparison: %s\n", err.Error())
	} else if pc2 == nil {
		t.Fatalf("Partial comparison not inserted\n")
	}

	if pc2.Dist != 0.24 {
		t.Fatal("Updated partial comparison data does not match")
	}
}

func TestPartialComparisonServiceDelete(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err = partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	err = partialComparisonService.Delete(&pc)
	if err != nil {
		t.Fatalf("Error deleting partial comparison: %s\n", err.Error())
	}

	pc2, err := partialComparisonService.Get(pc.Id)
	if err != nil {
		t.Fatalf("Error getting deleted partial comparison: %s\n", err.Error())
	} else if pc2 != nil {
		t.Fatalf("partial comparison not deleted\n")
	}
}

func TestPartialComparisonServiceGetOneBy(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err = partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	pc2, err := partialComparisonService.GetOneBy("macro_partial_id = ? and gidx_partial_id = ?", pc.MacroPartialId, pc.GidxPartialId)
	if err != nil {
		t.Fatalf("Error getting one by partial comparison: %s\n", err.Error())
	} else if pc2 == nil {
		t.Fatalf("partial comparison not found by\n")
	}

	if pc2.MacroPartialId != pc.MacroPartialId ||
		pc2.GidxPartialId != pc.GidxPartialId {
		t.Fatal("partial comparison macro id does not match")
	}
}

func TestPartialComparisonServiceGetOneByNot(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	_, err = partialComparisonService.GetOneBy("macro_partial_id = ? and gidx_partial_id = ?", macroPartial.Id, gidxPartial.Id)
	if err == nil {
		t.Fatalf("Getting one by partial comparison did not fail")
	}
}

func TestPartialComparisonServiceExistsBy(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err = partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	found, err := partialComparisonService.ExistsBy("macro_partial_id = ? and gidx_partial_id = ?", pc.MacroPartialId, pc.GidxPartialId)
	if err != nil {
		t.Fatalf("Error getting one by partial comparison: %s\n", err.Error())
	}

	if !found {
		t.Fatalf("Partial comparison not exists by\n")
	}
}

func TestPartialComparisonServiceExistsByNot(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	found, err := partialComparisonService.ExistsBy("macro_partial_id = ? and gidx_partial_id = ?", macroPartial.Id, gidxPartial.Id)
	if err != nil {
		t.Fatalf("Error getting exists by partial comparison: %s\n", err.Error())
	}

	if found {
		t.Fatalf("Partial comparison exists by\n")
	}
}

func TestPartialComparisonServiceCount(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err = partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	num, err := partialComparisonService.Count()
	if err != nil {
		t.Fatalf("Error counting partial comparison: %s\n", err.Error())
	}

	if num != int64(1) {
		t.Fatalf("Partial comparison count incorrect\n")
	}
}

func TestPartialComparisonServiceCountBy(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err = partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	num, err := partialComparisonService.CountBy("macro_partial_id = ? and gidx_partial_id = ?", pc.MacroPartialId, pc.GidxPartialId)
	if err != nil {
		t.Fatalf("Error counting by partial comparison: %s\n", err.Error())
	}

	if num != int64(1) {
		t.Fatalf("Partial comparison count incorrect\n")
	}
}

func TestPartialComparisonServiceFindAll(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err = partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	pcs, err := partialComparisonService.FindAll("id DESC", 1000, 0, "macro_partial_id = ?", macroPartial.Id)
	if err != nil {
		t.Fatalf("Error finding all partial comparisons: %s\n", err.Error())
	}

	if pcs == nil {
		t.Fatalf("No partial comparison slice returned for FindAll\n")
	}

	if len(pcs) != 1 {
		t.Fatal("Inserted partial comparison not found by FindAll")
	}

	pc2 := pcs[0]

	if pc2.MacroPartialId != pc.MacroPartialId ||
		pc2.GidxPartialId != pc.GidxPartialId {
		t.Fatal("partial comparison macro id does not match")
	}
}

func TestPartialComparisonServiceFindOrCreate(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	partialComparison, err := partialComparisonService.FindOrCreate(&macroPartial, &gidxPartial)
	if err != nil {
		t.Fatalf("Failed to FindOrCreate partialComparison: %s\n", err.Error())
	}

	if partialComparison.MacroPartialId != macroPartial.Id {
		t.Fatalf("partialComparison.MacroPartialId was %d, expected %d\n", partialComparison.MacroPartialId, macroPartial.Id)
	}

	if partialComparison.GidxPartialId != gidxPartial.Id {
		t.Fatalf("partialComparison.GidxPartialId was %d, expected %d\n", partialComparison.GidxPartialId, gidxPartial.Id)
	}

	if partialComparison.Dist == 0.0 {
		t.Fatalf("partial comparison dist was 0.0")
	}
}

func TestPartialComparisonServiceCountMissing(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	num, err := partialComparisonService.CountMissing(&macro)
	if err != nil {
		t.Fatalf("Error counting missing partial comparisons: %s\n", err.Error())
	}

	if num != 10 {
		t.Fatalf("Expected 10 missing partial comparisons, got %d\n", num)
	}
}

func TestPartialComparisonServiceFindMissing(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	macroGidxViews, err := partialComparisonService.FindMissing(&macro, 1000)
	if err != nil {
		t.Fatalf("Error finding missing partial comparisons: %s\n", err.Error())
	}

	if len(macroGidxViews) != 10 {
		t.Fatalf("Expected 10 missing partial comparisons, got %d\n", len(macroGidxViews))
	}
}

func TestPartialComparisonServiceCreateFromView(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	macroGidxViews, err := partialComparisonService.FindMissing(&macro, 1000)
	if err != nil {
		t.Fatalf("Error finding missing partial comparisons: %s\n", err.Error())
	}

	if len(macroGidxViews) != 10 {
		t.Fatalf("Expected 10 missing partial comparisons, got %d\n", len(macroGidxViews))
	}

	view := macroGidxViews[0]
	pc, err := partialComparisonService.CreateFromView(view)
	if err != nil {
		t.Fatalf("Error creating partial comparison from view: %s\n", err.Error())
	}

	if pc == nil {
		t.Fatal("Partial comparison not created from view")
	}

	if pc.Id == int64(0) {
		t.Fatal("Partial comparison from view not given id")
	}

	if pc.Dist == 0.0 {
		t.Fatal("Partial comparison dist not calculated")
	}

}
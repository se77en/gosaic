package controller

import (
	"fmt"
	"testing"
)

type argTestIn struct {
	width, height, size, minDepth, maxDepth, minArea, maxArea int
}

type argTestOut struct {
	size, minDepth, maxDepth, minArea, maxArea int
	err                                        string
}

func TestMacroQuadFixArgs(t *testing.T) {
	for _, tt := range []struct {
		t argTestIn
		e argTestOut
	}{
		//(1200, 1200, -1, -1, -1, -1, -1) => (291, 4, 6, 1225, 0, ), expect (291, 5, 6, 1225, 0, )
		//(2400, 2400, -1, -1, -1, -1, -1) => (506, 4, 8, 1225, 160000, ), expect (506, 5, 8, 1225, 40000, )
		//(3600, 3600, -1, -1, -1, -1, -1) => (700, 5, 9, 1794, 360000, ), expect (700, 5, 9, 1794, 90000, )

		{
			argTestIn{100, 100, -1, -1, -1, -1, -1},
			argTestOut{40, 3, 4, 1225, 0, ""},
		},
		{
			argTestIn{300, 300, -1, -1, -1, -1, -1},
			argTestOut{96, 3, 4, 1225, 0, ""},
		},
		{
			argTestIn{600, 600, -1, -1, -1, -1, -1},
			argTestOut{167, 4, 5, 1225, 0, ""},
		},
		{
			argTestIn{1200, 1200, -1, -1, -1, -1, -1},
			argTestOut{291, 4, 6, 1225, 0, ""},
		},
		{
			argTestIn{2400, 2400, -1, -1, -1, -1, -1},
			argTestOut{506, 4, 8, 1225, 160000, ""},
		},
		{
			argTestIn{3600, 3600, -1, -1, -1, -1, -1},
			argTestOut{700, 5, 9, 1794, 360000, ""},
		},
	} {
		size, minDepth, maxDepth, minArea, maxArea, err := macroQuadFixArgs(
			tt.t.width, tt.t.height, tt.t.size, tt.t.minDepth, tt.t.maxDepth, tt.t.minArea, tt.t.maxArea)
		var errStr string
		if err != nil {
			errStr = err.Error()
		} else {
			errStr = ""
		}
		if tt.e.size != size ||
			tt.e.minDepth != minDepth ||
			tt.e.maxDepth != maxDepth ||
			tt.e.minArea != minArea ||
			tt.e.maxArea != maxArea ||
			tt.e.err != errStr {
			t.Errorf("macroQuadFixArgs(%d, %d, %d, %d, %d, %d, %d) => (%d, %d, %d, %d, %d, %s), expect (%d, %d, %d, %d, %d, %s)",
				tt.t.width, tt.t.height, tt.t.size, tt.t.minDepth, tt.t.maxDepth, tt.t.minArea, tt.t.maxArea,
				size, minDepth, maxDepth, minArea, maxArea, errStr,
				tt.e.size, tt.e.minDepth, tt.e.maxDepth, tt.e.minArea, tt.e.maxArea, tt.e.err)
		}
	}
}

func TestMacroQuad(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	cover, macro := MacroQuad(env, "testdata/jumping_bunny.jpg", 200, 200, 10, -1, 2, 50, -1, "", "")
	if cover == nil || macro == nil {
		fmt.Println(out.String())
		t.Fatal("Failed to create cover or macro")
	}

	expect := []string{
		"Building macro quad with 10 splits, 34 partials, min depth 2, max depth 2, min area 50...",
	}

	testResultExpect(t, out.String(), expect)
}

package transform

import "testing"

func TestNew(t *testing.T) {
	options := &ReaderOptions{
		Chunks:   64,
		Filename: "../../hack/Morgendämmerung.mp3",
		Mode:     TransformerModeRootMeanSquare,
	}

	ctx, err := New(options)
	if err != nil {
		t.Fatal(err)
	}

	blocks := ctx.Blocks()
	if len(blocks) != 64 {
		t.Fatalf("wrong number of chunks, expected %d, found %d", 64, len(blocks))
	}

	for _, sample := range blocks {
		if sample != 0 {
			return
		}
	}
	t.Fatalf("block slice only contains 0 entries, expected at least one non-null sample")
}

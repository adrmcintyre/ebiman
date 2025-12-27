package video

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shader.kage
var shaderProgram []byte

// Init initialises the shader for post-processing video output.
func (v *Video) Init() error {
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil {
		return err
	}
	v.shader = shader
	return nil
}

func (v *Video) SetChromaShift(f float64) {
	v.chromaShift = f
}

// PostProcess the video frame to simulate a phosphor display.
func (v *Video) PostProcess(dst ebiten.FinalScreen, src *ebiten.Image) {
	srcBounds := src.Bounds()
	dstBounds := dst.Bounds()

	var vert [4]ebiten.Vertex
	// set the source image sampling coordinates
	vert[0].SrcX = float32(srcBounds.Min.X) // top-left
	vert[0].SrcY = float32(srcBounds.Min.Y) // top-left
	vert[1].SrcX = float32(srcBounds.Max.X) // top-right
	vert[1].SrcY = float32(srcBounds.Min.Y) // top-right
	vert[2].SrcX = float32(srcBounds.Min.X) // bottom-left
	vert[2].SrcY = float32(srcBounds.Max.Y) // bottom-left
	vert[3].SrcX = float32(srcBounds.Max.X) // bottom-right
	vert[3].SrcY = float32(srcBounds.Max.Y) // bottom-right

	// set the destination image target coordinates
	vert[0].DstX = float32(dstBounds.Min.X + v.offsetX) // top-left
	vert[0].DstY = float32(dstBounds.Min.Y + v.offsetY) // top-left
	vert[1].DstX = float32(dstBounds.Max.X + v.offsetX) // top-right
	vert[1].DstY = float32(dstBounds.Min.Y + v.offsetY) // top-right
	vert[2].DstX = float32(dstBounds.Min.X + v.offsetX) // bottom-left
	vert[2].DstY = float32(dstBounds.Max.Y + v.offsetY) // bottom-left
	vert[3].DstX = float32(dstBounds.Max.X + v.offsetX) // bottom-right
	vert[3].DstY = float32(dstBounds.Max.Y + v.offsetY) // bottom-right

	// triangle shader options
	var shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderOpts.Images[0] = src
	shaderOpts.Uniforms = map[string]any{
		"ChromaShift": v.chromaShift,
	}

	// draw shader
	indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	dst.DrawTrianglesShader(vert[:], indices, v.shader, &shaderOpts)
}

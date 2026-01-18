package units

// ResolveWidth resolves a dimension as a width value in pixels
func (d Dimension) ResolveWidth(ctx ResolveContext) int {
	switch d.Unit {
	case Em:
		return int(d.Value * float64(ctx.EmHeight))
	case Percent:
		return int(d.Value / 100.0 * float64(ctx.WorkArea.Width))
	default:
		return int(d.Value)
	}
}

// ResolveHeight resolves a dimension as a height value in pixels
func (d Dimension) ResolveHeight(ctx ResolveContext) int {
	switch d.Unit {
	case Em:
		return int(d.Value * float64(ctx.EmHeight))
	case Percent:
		return int(d.Value / 100.0 * float64(ctx.WorkArea.Height))
	default:
		return int(d.Value)
	}
}

// ResolveSize resolves both width and height dimensions to pixels
func ResolveSize(width, height Dimension, ctx ResolveContext) (w, h int) {
	return width.ResolveWidth(ctx), height.ResolveHeight(ctx)
}

package text_alignment_constants

type HAlignment string
type VAlignment string

const (
	Left           HAlignment = "left"
	Center         HAlignment = "center"
	Right          HAlignment = "right"
	Justify        HAlignment = "justify"
	Distributed    HAlignment = "distributed"
	Top            VAlignment = "top"
	CenterVertical VAlignment = "center"
	Bottom         VAlignment = "bottom"
	JustifyV       VAlignment = "justify"
	DistributedV   VAlignment = "distributed"
)

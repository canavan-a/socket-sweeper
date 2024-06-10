package board

type Voxel struct {
	IsBomb bool
	Number int
	IsOpen bool
}

func (v *Voxel) Open() {
	v.IsOpen = true
}

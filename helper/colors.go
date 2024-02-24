package helper

func PickColor() string {
	colors := GetColors()

	return colors[Random(len(colors)-1)]
}

func GetColors() []string {
	return []string{"blue", "green", "yellow", "pink", "violet", "orange"}
}

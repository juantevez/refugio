package domain

// Point representa una coordenada geográfica (GEOGRAPHY(POINT, 4326)).
// Longitud va primero para ser compatible con ST_MakePoint(long, lat).
type Point struct {
	Lat  float64
	Long float64
}

// SearchArea define el área de búsqueda para encontrar reportes cercanos.
type SearchArea struct {
	Center       Point
	RadiusMeters int
}

package postgis

// Registry holds all tokens to parse postgis SQL extensions.
type Registry struct {
}

// Name of this registry.
func (Registry) Name() string {
	return "postgis"
}

// Operators supported by postgres v14.
func (Registry) Operators() []string {
	return []string{
		"&&&",
		"|=|",
		"<#>",
		"<<->>",
		"<<#>>",
	}
}

// Functions supported by postgres v14.
func (Registry) Functions() []string {
	return []string{
		"AddGeometryColumn",
		"DropGeometryColumn",
		"DropGeometryTable",
		"Find_SRID",
		"Populate_Geometry_ColumnS",
		"UpdateGeometrySRID",
		"ST_Collect",
		"ST_LineFromMultipoint",
		"ST_MakeEnvelope",
		"ST_MakeLine",
		"ST_MakePoint",
		"ST_MakePointM",
		"ST_MakePolygon",
		"ST_Point",
		"ST_PointZ",
		"ST_PointM",
		"ST_PointZM",
		"ST_Polygon",
		"ST_TileEnvelope",
		"ST_HexagonGrid",
		"ST_Hexagon",
		"ST_SquareGrid",
		"ST_Square",
		"GeometryType",
		"ST_Boundary",
		"ST_BoundingDiagonal",
		"ST_CoordDim",
		"ST_Dimension",
		"ST_Dump",
		"ST_DumpPointS",
		"ST_DumpSegments",
		"ST_DumpRingS",
		"ST_EndPoint",
		"ST_Envelope",
		"ST_ExteriorRing",
		"ST_GeometryN",
		"ST_GeometryType",
		"ST_HasArc",
		"ST_InteriorRingN",
		"ST_IsClosed",
		"ST_IsCollection",
		"ST_IsEmpty",
		"ST_IsPolygonCCW",
		"ST_IsPolygonCW",
		"ST_IsRing",
		"ST_IsSimple",
		"ST_M",
		"ST_MemSize",
		"ST_NDims",
		"ST_NPoints",
		"ST_NRings",
		"ST_NumGeometries",
		"ST_NumInteriorRings",
		"ST_NumInteriorRing",
		"ST_NumPatches",
		"ST_NumPoints",
		"ST_PatchN",
		"ST_PointN",
		"ST_Points",
		"ST_StartPoint",
		"ST_Summary",
		"ST_X",
		"ST_Y",
		"ST_Z",
		"ST_Zmflag",
		"ST_AddPoint",
		"ST_CollectionExtract",
		"ST_CollectionHomogenize",
		"ST_CurveToLine",
		"ST_Scroll",
		"ST_FlipCoordinates",
		"ST_Force2D",
		"ST_Force3D",
		"ST_Force3DZ",
		"ST_Force3DM",
		"ST_Force4D",
		"ST_ForcePolygonCCW",
		"ST_ForceCollection",
		"ST_ForcePolygonCW",
		"ST_ForceSFS",
		"ST_ForceRHR",
		"ST_ForceCurve",
		"ST_LineToCurve",
		"ST_Multi",
		"ST_Normalize",
		"ST_QuantizeCoordinates",
		"ST_RemovePoint",
		"ST_RemoveRepeatedPoints",
		"ST_Reverse",
		"ST_Segmentize",
		"ST_SetPoint",
		"ST_ShiftLongitude",
		"ST_WrapX",
		"ST_SnapToGrid",
		"ST_Snap",
		"ST_SwapOrdinates",
		"ST_IsValid",
		"ST_IsValidDetail",
		"ST_IsValidReason",
		"ST_MakeValid",
		"ST_SetSRID",
		"ST_SRID",
		"ST_Transform",
		// ... https://postgis.net/docs/reference.html §  8.8
	}
}

// Types supported by postgres v14.
func (Registry) Types() []string {
	return []string{
		"GEOGRAPHY",
		"GEOMETRY",
		"GEOMETRY_Dump",
		"BOX2D",
		"BOX3D",
	}
}

// ReservedValues supported by postgres v14.
func (Registry) ReservedValues() []string {
	return []string{}
}

// ConstantBuilders yields contant builders supported by postgres v14.
func (Registry) ConstantBuilders() []string {
	return []string{}
}

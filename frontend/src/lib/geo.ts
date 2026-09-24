export type LatLng = [number, number];

/**
 * Ambil koordinat [lng, lat] dari geometri yang bisa berupa:
 * - string JSON (mis. hasil SQL ST_AsGeoJSON)
 * - object { coordinates: [lng, lat] }
 * - undefined/null (fallback [0, 0], aman — tidak pernah throw)
 */
export function getGeoCoordinates(geom: unknown): LatLng {
  if (typeof geom === "string") {
    try {
      const parsed = JSON.parse(geom);
      return Array.isArray(parsed?.coordinates) ? parsed.coordinates : [0, 0];
    } catch {
      return [0, 0];
    }
  }
  if (geom && typeof geom === "object") {
    const coordinates = (geom as { coordinates?: unknown }).coordinates;
    if (Array.isArray(coordinates)) return coordinates as LatLng;
  }
  return [0, 0];
}

/** Buang titik yang persis berurutan duplikat (server sudah mengembalikan endpoint). */
export function dedupeCoords(coords: LatLng[]): LatLng[] {
  const out: LatLng[] = [];
  for (const c of coords) {
    const prev = out[out.length - 1];
    if (!prev || prev[0] !== c[0] || prev[1] !== c[1]) {
      out.push(c);
    }
  }
  return out;
}
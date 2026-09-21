import { json } from '@sveltejs/kit';

export async function GET({ url, fetch }) {
  const start = url.searchParams.get('start');
  const end = url.searchParams.get('end');
  
  if (!start || !end) {
    return json({ error: "Missing start or end coordinates" }, { status: 400 });
  }

  const sLng = parseFloat(start.split(',')[0]);
  const sLat = parseFloat(start.split(',')[1]);
  const eLng = parseFloat(end.split(',')[0]);
  const eLat = parseFloat(end.split(',')[1]);

  // Poligon Kawasan Rawan Bencana (KRB III) & Jalur Lahar Utama (Source: BMKG & NASA)
  const hazardPolygon = [
    [
      [112.290, -7.900], // North of crater
      [112.330, -7.920], // East of crater
      [112.320, -7.960], // South of crater
      [112.280, -7.960], // South-West of crater
      [112.230, -7.930], // South of Ngobo
      [112.210, -7.900], // End of Ngobo
      [112.230, -7.880], // North of Ngobo
      [112.290, -7.900]  // Close polygon
    ]
  ];

  // Strategi 1: OpenRouteService dengan Fitur Dynamic Hazard Avoidance (Sesuai Proposal)
  try {
    const orsUrl = `https://api.openrouteservice.org/v2/directions/driving-car/geojson`;
    const body = {
      coordinates: [[sLng, sLat], [eLng, eLat]],
      options: {
        avoid_polygons: {
          type: "Polygon",
          coordinates: hazardPolygon
        }
      }
    };

    const res = await fetch(orsUrl, {
      method: 'POST',
      headers: {
        'Authorization': '5b3ce3597851110001cf6248b3b4e6e1e1a14a0e8c68e8d25f2e58e3',
        'Content-Type': 'application/json',
        'Accept': 'application/json, application/geo+json, application/gpx+xml, img/png; charset=utf-8'
      },
      body: JSON.stringify(body)
    });

    if (res.ok) {
      const orsData = await res.json();
      if (orsData.features && orsData.features.length > 0) {
        const coords = orsData.features[0].geometry.coordinates;
        const props = orsData.features[0].properties.summary;
        return json({
          routes: [{
            geometry: { coordinates: coords, type: "LineString" },
            distance: props.distance,
            duration: props.duration,
            isSafe: true // Flag rute ini dijamin aman dari zona merah
          }]
        });
      }
    }
  } catch (e) {
    console.error("ORS Hazard Avoidance gagal:", e);
  }

  // Strategi 2: Fallback ke OSRM standar jika titik awal/akhir terperangkap di dalam poligon
  try {
    const osrmUrl = `https://router.project-osrm.org/route/v1/driving/${start};${end}?overview=full&geometries=geojson&alternatives=false`;
    const res = await fetch(osrmUrl, {
      headers: { 'User-Agent': 'LOKARI-Disaster-Evacuation-App/1.0' }
    });
    
    if (res.ok) {
      const data = await res.json();
      if (data.routes && data.routes.length > 0) {
        return json({ ...data, routes: [{ ...data.routes[0], isSafe: false }] }); // Rute darurat (mungkin melintasi zona bahaya)
      }
    }
  } catch (e) {
    console.error("OSRM gagal:", e);
  }

  return json({ error: "Semua routing engine gagal" }, { status: 502 });
}

import { json } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export async function GET({ url, fetch }) {
  const start = url.searchParams.get('start');
  const end = url.searchParams.get('end');
  
  if (!start || !end) {
    return json({ error: "Missing start or end coordinates" }, { status: 400 });
  }

  // Validasi ketat: 4 nilai harus finite & dalam rentang geografis.
  // Number (bukan parseFloat) menolak input seperti "12.3abc".
  const [sLngRaw, sLatRaw] = start.split(',');
  const [eLngRaw, eLatRaw] = end.split(',');
  const sLng = Number(sLngRaw);
  const sLat = Number(sLatRaw);
  const eLng = Number(eLngRaw);
  const eLat = Number(eLatRaw);

  const invalid =
    ![sLng, sLat, eLng, eLat].every(Number.isFinite) ||
    Math.abs(sLng) > 180 || Math.abs(eLng) > 180 ||
    Math.abs(sLat) > 90 || Math.abs(eLat) > 90;

  if (invalid) {
    return json(
      { error: "Invalid coordinates: expected start=lng,lat&end=lng,lat" },
      { status: 400 }
    );
  }

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
  const orsApiKey = env.ORS_API_KEY ?? '';

  try {
    if (!orsApiKey) {
      throw new Error("ORS_API_KEY tidak diset — langsung pakai fallback OSRM");
    }

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
        'Authorization': orsApiKey,
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
            isSafe: true, // Flag rute ini dijamin aman dari zona merah
            status: "safe"
          }]
        });
      }
    } else {
      const errText = await res.text();
      console.error("ORS API Error:", res.status, errText);
    }
  } catch (e) {
    console.error("ORS Hazard Avoidance gagal:", e);
  }

  // Strategi 2: Fallback ke OSRM standar jika titik awal/akhir terperangkap di dalam poligon
  try {
    const osrmUrl = `https://router.project-osrm.org/route/v1/driving/${sLng},${sLat};${eLng},${eLat}?overview=full&geometries=geojson&alternatives=false`;
    const res = await fetch(osrmUrl, {
      headers: { 'User-Agent': 'LOKARI-Disaster-Evacuation-App/1.0' }
    });
    
    if (res.ok) {
      const data = await res.json();
      if (data.routes && data.routes.length > 0) {
        return json({ ...data, routes: [{ ...data.routes[0], isSafe: false, status: "fallback" }] }); // Rute darurat / Fallback
      }
    }
  } catch (e) {
    console.error("OSRM gagal:", e);
  }

  return json({ error: "Semua routing engine gagal" }, { status: 502 });
}

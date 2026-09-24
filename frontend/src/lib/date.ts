/**
 * Format timestamp RFC3339 dari backend (mis. "2026-09-23T22:44:39+07:00")
 * ke bentuk tanggal Indonesia yang mudah dibaca: "24 Sep 2026, 22:44 WIB".
 *
 * Selalu dirender dalam zona Asia/Jakarta (terlepas dari zona browser),
 * karena backend sudah menuliskan offset Asia/Jakarta.
 * Kembalikan null bila input kosong/tidak dapat di-parse — pemanggil
 * yang memutuskan fallback (mis. i18n "Baru Saja").
 */
export function formatNewsDate(createdAt?: string | null): string | null {
  if (!createdAt) return null;
  const d = new Date(createdAt);
  if (Number.isNaN(d.getTime())) return null;

  const parts = new Intl.DateTimeFormat("id-ID", {
    timeZone: "Asia/Jakarta",
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hourCycle: "h23",
    timeZoneName: "short",
  }).formatToParts(d);

  const get = (type: string) => parts.find((p) => p.type === type)?.value ?? "";
  return `${get("day")} ${get("month")} ${get("year")}, ${get("hour")}:${get("minute")} ${get("timeZoneName")}`;
}
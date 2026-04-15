/**
 * 从 YouTube embed / watch / youtu.be 链接解析视频 id，用于官方缩略图。
 */
export function youtubeVideoIdFromUrl(url: string): string | null {
  try {
    const u = new URL(url, "https://www.youtube.com");
    if (u.hostname === "youtu.be") {
      const id = u.pathname.replace(/^\//, "").split("/")[0];
      return id || null;
    }
    if (u.pathname.startsWith("/embed/")) {
      const id = u.pathname.replace(/^\/embed\//, "").split("/")[0];
      return id || null;
    }
    const v = u.searchParams.get("v");
    if (v) {
      return v;
    }
  } catch {
    // ignore
  }
  const embed = url.match(/youtube\.com\/embed\/([^?&/]+)/);
  if (embed?.[1]) {
    return embed[1];
  }
  const short = url.match(/youtu\.be\/([^/?]+)/);
  if (short?.[1]) {
    return short[1];
  }
  const vParam = url.match(/[?&]v=([^&]+)/);
  if (vParam?.[1]) {
    return vParam[1];
  }
  return null;
}

export function youtubePosterUrl(url: string): string | undefined {
  const id = youtubeVideoIdFromUrl(url);
  return id ? `https://img.youtube.com/vi/${id}/hqdefault.jpg` : undefined;
}

/** B 站播放器页 URL 中的 bvid */
export function bilibiliBvidFromPlayerUrl(url: string): string | null {
  try {
    const u = new URL(url);
    const bvid = u.searchParams.get("bvid");
    if (bvid) {
      return bvid;
    }
  } catch {
    // ignore
  }
  const m = url.match(/bvid=([^&]+)/i);
  return m?.[1] || null;
}

/**
 * 服务端代理 B 站稿件信息，避免浏览器直连 api.bilibili.com 的 CORS 限制。
 */
export default defineEventHandler(async (event) => {
  const raw = String(getQuery(event).bvid || "").trim();
  if (!/^BV[\w]+$/i.test(raw)) {
    return { pic: null as string | null };
  }
  try {
    const res = await $fetch<{ data?: { pic?: string } }>(
      `https://api.bilibili.com/x/web-interface/view?bvid=${encodeURIComponent(raw)}`,
      {
        timeout: 10_000,
        headers: {
          "User-Agent": "Mozilla/5.0 (compatible; Moments/1.0)",
        },
      },
    );
    const pic = res?.data?.pic?.trim();
    return { pic: pic || null };
  } catch {
    return { pic: null as string | null };
  }
});

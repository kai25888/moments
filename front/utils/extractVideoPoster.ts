/**
 * 在浏览器内从本地视频文件截取一帧，生成 JPEG Blob（用于上传为封面）。
 * 使用 blob: URL，画布不会被跨域污染。
 */
export async function extractVideoPosterBlob(
  file: File,
  options?: { maxWidth?: number; quality?: number; seekRatio?: number },
): Promise<Blob | null> {
  const maxWidth = options?.maxWidth ?? 960;
  const quality = options?.quality ?? 0.82;
  const seekRatio = options?.seekRatio ?? 0.02;

  const objectUrl = URL.createObjectURL(file);
  const video = document.createElement("video");
  video.muted = true;
  video.playsInline = true;
  video.preload = "auto";
  video.src = objectUrl;

  const cleanup = () => {
    URL.revokeObjectURL(objectUrl);
    video.removeAttribute("src");
    video.load();
  };

  try {
    await new Promise<void>((resolve, reject) => {
      const onMeta = () => {
        video.removeEventListener("loadedmetadata", onMeta);
        resolve();
      };
      video.addEventListener("loadedmetadata", onMeta);
      video.addEventListener(
        "error",
        () => reject(new Error("video metadata load failed")),
        { once: true },
      );
      video.load();
    });

    const duration = Number.isFinite(video.duration) && video.duration > 0 ? video.duration : 0;
    const seekTo =
      duration > 0.2 ? Math.min(Math.max(duration * seekRatio, 0.05), duration - 0.05) : 0.05;

    await new Promise<void>((resolve, reject) => {
      const onSeeked = () => {
        video.removeEventListener("seeked", onSeeked);
        resolve();
      };
      video.addEventListener("seeked", onSeeked);
      video.addEventListener(
        "error",
        () => reject(new Error("video seek failed")),
        { once: true },
      );
      video.currentTime = seekTo;
    });

    const vw = video.videoWidth;
    const vh = video.videoHeight;
    if (!vw || !vh) {
      return null;
    }

    const scale = Math.min(1, maxWidth / vw);
    const cw = Math.max(1, Math.round(vw * scale));
    const ch = Math.max(1, Math.round(vh * scale));

    const canvas = document.createElement("canvas");
    canvas.width = cw;
    canvas.height = ch;
    const ctx = canvas.getContext("2d");
    if (!ctx) {
      return null;
    }
    ctx.drawImage(video, 0, 0, cw, ch);

    return await new Promise<Blob | null>((resolve) => {
      canvas.toBlob((b) => resolve(b), "image/jpeg", quality);
    });
  } catch {
    return null;
  } finally {
    cleanup();
  }
}

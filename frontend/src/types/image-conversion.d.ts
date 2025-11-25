declare module 'image-conversion' {
  export function compressAccurately(
    file: Blob,
    options: { size?: number; accuracy?: number; type?: string }
  ): Promise<Blob>
  export function compress(
    file: Blob,
    options?: { quality?: number; type?: string; width?: number; height?: number }
  ): Promise<Blob>
}

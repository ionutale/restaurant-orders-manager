# 0040 — Performance: Image upload and serving

## What to build

Test and optimize the file upload and serving pipeline. Verify upload time, file size limits, image compression, and caching headers.

## Acceptance criteria

- [ ] Test upload of 1KB, 100KB, 1MB, 5MB, 10MB images
- [ ] 10MB upload returns 413 (payload too large) or is rejected gracefully
- [ ] Upload speeds are measured and documented
- [ ] Served images have correct Cache-Control headers (immutable, max-age)
- [ ] Nginx serves uploaded files with proper Content-Type
- [ ] Uploading a non-image file (e.g. .exe, .pdf) is rejected
- [ ] Uploading an image with wrong extension but valid content is handled
- [ ] Concurrency: 10 simultaneous uploads do not cause errors
- [ ] Disk space monitoring: uploads directory does not grow unbounded

## Blocked by

- 0006 — Allergens + dish suggestions (file upload was added here)

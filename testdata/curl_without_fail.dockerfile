FROM scratch

# Keep --fail non-adjacent to curl and target URL
RUN curl --tlsv1.2 --fail --http2 https://example.com/file.json

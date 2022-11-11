
## Send story from a client to the server.
The story contains formatted text in an HTML file and other file types which are photos, gifs, videos and in some cases audios and other file types.

- The entire files gets sent from a client to server as a chunk.
- Every time a user sends a story, a new directory gets created in the cloud, named as story id.
- And the files(chunks) gets saved into the dir(story id)

Example:
|                   folder              |               Files                   |
|---------------------------------------|---------------------------------------|
| 36rouwYPVYRodi4BoW2Lav8tnGzpuqjzkR8a  | index.html, vid.mp4                   |

# mp3dir

mp3dir is a zero-configuration tool for creating a strictly lossy copy of a music
folder.

## Example

```
./mp3dir -i ~/Music -o ~/Music.mp3
```

## Dependencies

Both `ffmpeg` and `ffprobe` need to be in PATH at runtime.

## Transformation matrix

mp3dir will transform files according to the following rules:

| input         | transformation    |
|---------------|-------------------|
| flac (.flac)  | convert to mp3 v0 |
| alac (.m4a)   | convert to mp3 v0 |
| mp3 (.mp3)    | copy              |
| aac (.m4a)    | copy              |

## License

Copyright © 2021 Martin Frederic 

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the
“Software”), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS
OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

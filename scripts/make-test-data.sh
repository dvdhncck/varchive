
rm -rf test
mkdir -p test/one

ffmpeg -f lavfi -i testsrc=duration=5:size=640x320:rate=25 "test/one/video-5s.mpg"
ffmpeg -f lavfi -i testsrc=duration=7:size=640x320:rate=25 "test/one/video-7s.mpg"
ffmpeg -f lavfi -i testsrc=duration=9:size=640x320:rate=25 "test/one/video-9s.mpg"

ffmpeg -f lavfi -i "sine=frequency=1000:duration=5" "test/one/audio-1kHz5seconds.mp3"
ffmpeg -f lavfi -i "sine=frequency=5000:duration=7" "test/one/audio-5kHz7seconds.mp3"
ffmpeg -f lavfi -i "sine=frequency=9000:duration=9" "test/one/audio-9kHz9seconds.mp3"

ffmpeg -i "test/one/video-5s.mpg" -i "test/one/audio-1kHz5seconds.mp3" \
       -c copy -map 0:v:0 -map 1:a:0 -shortest "test/one/sample 001.mpg"
ffmpeg -i "test/one/video-7s.mpg" -i "test/one/audio-5kHz7seconds.mp3" \
       -c copy -map 0:v:0 -map 1:a:0 -shortest "test/one/sample 002.mpg"
ffmpeg -i "test/one/video-9s.mpg" -i "test/one/audio-9kHz9seconds.mp3" \
       -c copy -map 0:v:0 -map 1:a:0 -shortest "test/one/sample 003.mpg"

rm -rf test/one/audio-*
rm -rf test/one/video-*
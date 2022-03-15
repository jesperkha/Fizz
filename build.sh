set -e

echo "[INFO] Adding libraries"
python lib/build.py

echo "[INFO] Generating library docs"
python lib/autodocs.py

echo "[INFO] Building binary"
[ ! -d "./bin" ] && mkdir bin && echo "[INFO] Created /bin directory"

CMD="go build -o bin/fizz.exe ."
echo "[CMD] $CMD"
$CMD

echo "[INFO] Finished build of " && ./bin/fizz --version
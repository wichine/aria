# !/bin/bash
version=0.1.0
set -e

echo "   _____             .__          ";
echo "  /  _  \   _______  |__| _____   ";
echo " /  /_\  \  \_  __ \ |  | \__  \  ";
echo "/    |    \  |  | \/ |  |  / __ \_";
echo "\____|__  /  |__|    |__| (____  /";
echo "        \/                     \/ ";
echo "                                  ";


echo "Start building ...";
echo " ";

cur_dir="$(cd `dirname $0`;pwd)"
echo "Change to build dir ..."
cd ${cur_dir}
echo "[OK]"

echo "Inject assets to assets.go ..."
go run *.go selfbuild inject
echo "[OK]"

echo "Build aria ..."
go build -o aria-$(echo $(uname -s) | tr '[A-Z]' '[a-z]')-${version}
echo "[OK]"

echo "Restore assets.go to template ..."
go run *.go selfbuild restore
echo "[OK]"

echo "[Finish]"

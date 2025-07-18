set -e
set -o pipefail

#------

rm -rf cmd/mcp-victorialogs/resources/vm

git clone --no-checkout --depth=1 https://github.com/VictoriaMetrics/VictoriaMetrics.git cmd/mcp-victorialogs/resources/vm
cd cmd/mcp-victorialogs/resources/vm

git sparse-checkout init --cone
git sparse-checkout set docs
git checkout master
rm -rf ./.git
rm -f ./docs/Makefile ./Makefile ./LICENSE ./*.md ./*.mod ./*.sum ./*.zip ./.golangci.yml ./.wwhrd.yml ./.gitignore ./.dockerignore ./codecov.yml

rm -rf ./docs/anomaly-detection
rm -rf ./docs/victoriametrics-cloud
rm -rf ./docs/victoriametrics

find docs/guides/*  | grep -v logs | xargs -I {} rm -rf {}

cd -

#------

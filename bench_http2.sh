#!/usr/bin/env bash
set -euo pipefail

if [[ ${1-} == "" ]]; then
  echo "Usage: $0 <path-to-ripestat-binary> [resource ...]" >&2
  exit 1
fi

BINARY=$1
shift

if [[ ! -x $BINARY ]]; then
  echo "Binary '$BINARY' is not executable" >&2
  exit 1
fi

if [[ $# -eq 0 ]]; then
  # Default sample set hits both ASN and IP paths.
  RESOURCES=(13335 15169 8.8.8.8 1.1.1.1 2001:4860:4860::8888)
else
  RESOURCES=($@)
fi

RUNS=${BENCH_RUNS:-3}
CONCURRENCY=${RIPESTAT_MAX_CONCURRENCY:-8}

echo "Benchmarking $BINARY using RIPESTAT_MAX_CONCURRENCY=$CONCURRENCY"
echo "Resources: ${RESOURCES[*]}"
echo "Runs: $RUNS"

for ((i=1; i<=RUNS; i++)); do
  echo "-- Run $i/$RUNS --"
  TIMEFORMAT='%3R seconds'
  time "$BINARY" "${RESOURCES[@]}" >/dev/null
  echo
done

echo "Benchmark complete"

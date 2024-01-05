#!/usr/bin/env bash

SCRIPT_DIR=$(dirname "${0}")

source "${SCRIPT_DIR}"/set_vars.sh

# Public: Compiles the PyTeal to TEAL source code.
#
# Examples:
#
#   ./bin/compile.sh
#
# Returns exit code 0.
function main {
  local compiled_filename
  local build_dir
  local src_dir
  local src_files

  set_vars

  build_dir="${PWD}/.build"
  src_dir="${PWD}/contract"

  # if the .build/ directory does not exist, create it
  if [[ ! -d "${build_dir}" ]]; then
    printf "%b creating %b directory... \n" "${INFO_PREFIX}" "${build_dir}"
    mkdir -p "${build_dir}"
  fi

  src_files=$(ls -p "${src_dir}" | grep -v /) # only get files

  for src_filename in ${src_files}; do
    compiled_filename="${src_filename%.py}.teal"

    # remove the previous compiled teal file, if it exists
    if [[ -f "${build_dir}/${compiled_filename}" ]]; then
      rm "${build_dir}/${compiled_filename}"
    fi

    # compile pyteal code to teal
    python3 "${src_dir}/${src_filename}" >> "${build_dir}/${compiled_filename}"

    printf "%b compiled teal code to %b file... \n" "${INFO_PREFIX}" "${build_dir}/${compiled_filename}"
  done

  exit 0
}

# And so, it begins...
main

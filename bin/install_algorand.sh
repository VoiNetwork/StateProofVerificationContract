#!/usr/bin/env bash

SCRIPT_DIR=$(dirname "${0}")

source "${SCRIPT_DIR}"/set_vars.sh

# Public: Creates (or re-creates) a .algorand directory and clones the Algorand sandbox to it.
#
# Examples
#
#   ./bin/install_algorand.sh
#
# Returns exit code 0.
function main {
  local algorand_dir

  set_vars

  algorand_dir="${PWD}/.algorand"

  if [[ -d "${algorand_dir}" ]]; then
    printf "%b deleting previous .algorand directory... \n" "${INFO_PREFIX}"
    rm -rf "${algorand_dir}"
  fi

  mkdir -p "${algorand_dir}"

  # clone algorand sandbox to the .algorand directory
  git clone git@github.com:algorand/sandbox.git "${algorand_dir}"

  exit 0
}

# And so, it begins...
main

#!/usr/bin/env bash

SCRIPT_DIR=$(dirname "${0}")

source "${SCRIPT_DIR}"/set_vars.sh

# Public: Starts the Algorand sandbox. If the .algorand sandbox does not exist, it will install the sandbox.
#
# Examples
#
#   ./bin/start_algorand.sh
#
# Returns exit code 0.
function main {
  local algorand_dir

  set_vars

  algorand_dir="${PWD}/.algorand"

  if [[ ! -d "${algorand_dir}" ]]; then
    printf "%b no .algorand/ directory exists. installing algorand sandbox... \n" "${ERROR_PREFIX}" "${SCRIPT_DIR}"
    ./"${SCRIPT_DIR}"/install_algorand.sh
  fi

  # start private network, if this is the first time, this will take a while to download dependencies
  "${algorand_dir}"/sandbox up dev

  # keep algod container open and print log information
  "${algorand_dir}"/sandbox tail algod

  # when cacenlled clean up
  "${algorand_dir}"/sandbox down

  exit 0
}

# And so, it begins...
main

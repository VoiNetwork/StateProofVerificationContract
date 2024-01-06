#!/usr/bin/env bash

SCRIPT_DIR=$(dirname "${0}")

source "${SCRIPT_DIR}"/set_vars.sh

# Public: Starts up a private Algorand network, compiles the application code and runs tests against it.
#
# Examples
#
#   ./bin/test.sh
#
# Returns exit code 0.
function main {
  local algorand_dir
  local is_algod_running

  set_vars

  algorand_dir="${PWD}/.algorand"

  if [[ ! -d "${algorand_dir}" ]]; then
    printf "%b no .algorand/ directory exists. installing algorand sandbox... \n" "${ERROR_PREFIX}" "${SCRIPT_DIR}"
    ./"${SCRIPT_DIR}"/install_algorand.sh
  fi

  printf "%b staring private algorand network... \n" "${INFO_PREFIX}"

  is_algod_running=$(docker ps | grep algorand-sandbox-algod)

  # if it is not already running start a  private network
  if [[ -z "${is_algod_running}" ]]; then
     printf "%b algorand sandbox not running, initializing... \n" "${INFO_PREFIX}"
     "${algorand_dir}"/sandbox up dev
  fi

  # compile the teal code
  ./"${SCRIPT_DIR}"/compile.sh

  printf "%b running tests... \n" "${INFO_PREFIX}"

  # run tests
  go test

  # if the algod was not already running, stop the private network
  if [[ -z "${is_algod_running}" ]]; then
     printf "%b stopping algorand sandbox... \n" "${INFO_PREFIX}"
     "${algorand_dir}"/sandbox up dev
  fi

  exit 0
}

# And so, it begins...
main

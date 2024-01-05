#!/usr/bin/env bash

SCRIPT_DIR=$(dirname "${0}")

source "${SCRIPT_DIR}"/set_vars.sh

# Public: Starts up a private Algorand network, compiles the contract code and runs tests against it.
#
# Examples
#
#   ./bin/test.sh
#
# Returns exit code 0.
function main {
  local algorand_dir
  local version

  set_vars

  version=$(<VERSION) # get the version from the version file
  algorand_dir="${PWD}/.algorand"

  if [[ ! -d "${algorand_dir}" ]]; then
    printf "%b no .algorand/ directory exists. installing algorand sandbox... \n" "${ERROR_PREFIX}" "${SCRIPT_DIR}"
    ./"${SCRIPT_DIR}"/install_algorand.sh
  fi

  printf "%b staring private algorand network... \n" "${INFO_PREFIX}"

  # start private network, if this is the first time, this will take a while to download dependencies
  "${algorand_dir}"/sandbox up dev

  # compile the teal code
  ./"${SCRIPT_DIR}"/compile.sh

  # run tests
  go test -ldflags "-X stateproofverificationcontract.Version=$version"

  # stop network
  "${algorand_dir}"/sandbox down

  exit 0
}

# And so, it begins...
main

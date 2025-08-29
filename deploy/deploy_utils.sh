#!/bin/bash

function create_dirs () {
  local BASE_DIR=$1
  local DIRS=("${@:2}")
  for DIR in "${DIRS[@]}"; do
    echo "CREATE: $DIR in $BASE_DIR"
    mkdir -p "$BASE_DIR"/"$DIR"
  done
}

function create_log_dirs () {
  BASE_DIR=$1
  create_dirs "$BASE_DIR" "${@:2}"
}

function copy_symlink () {
  local SOURCE=$1
  local TARGET=$2
  rm -f "$TARGET"
  ln -sf "$SOURCE" "$TARGET"
}

function give_execution_permissions_if_shell() {
    local file_path=$1

    # Check if the file exists
    if [[ ! -f $file_path ]]; then
        echo "File $file_path does not exist."
        return 1
    fi

    # Use the file command to check the file type
    file_type=$(file "$file_path")

    # Check if the file type indicates a shell script
    if [[ $file_type == *"shell script"* ]]; then
        # Give execution rights
        chmod +x "$file_path"
        echo "Execution granted: $file_path"
        return 0
    else
        return 1
    fi
}

function update_symlinks () {
  local SOURCE=$1
  local TARGET=$2
  local ROOT_ORIGINAL=${3:-false}
  find "$SOURCE" -type f | while read -r sfile; do
    name=$(basename "$sfile")
    if [ "$ROOT_ORIGINAL" == true ]; then
      chown root:root "$sfile"
      chmod 0440 "$sfile"
    fi
    copy_symlink "$sfile" "$TARGET"/"$name"
    echo "REFRESH $name DONE"
  done
}

function apply_execution_permissions () {
  BASE_DIR=$1
  _apply_exec_perms "$BASE_DIR" "${@:2}"
}

function _apply_exec_perms () {
  local BASE_DIR=$1
  local DIRS=("${@:2}")
  for DIR in "${DIRS[@]}"; do
    give_execution_permissions "$BASE_DIR"/"$DIR"
  done
}

function give_execution_permissions () {
  local SOURCE=$1
  find "$SOURCE" -type f | while read -r sfile; do
    name=$(basename "$sfile")
    give_execution_permissions_if_shell "$sfile"
  done
}

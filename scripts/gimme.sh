#!/bin/bash

reserved="sync init discover list"

if [[ "$reserved" = *"$1"* ]]; then
  gimme-core "$1"
else

  gimmeRaw=$(echo "$@" | xargs gimme-core --manifest)
  if [[ $? == 0 ]]; then
    IFS=' ' read -r -a warp <<< "$(echo "$gimmeRaw" | grep "\[WARP\]")"
    environ="$(echo "$gimmeRaw" | grep "\[ENV\]")"
    gimmeCommands="$(echo "$gimmeRaw" | grep "\[CMD\]")"

    warpDest="${warp[1]}"
    cd "$warpDest" || return

    # Perform exports

    if [[ -n $environ ]]; then
      echo "Exporting the following variables:"
      while IFS= read -r line; do
        read -r -a tokens <<< "$line"
        export "${tokens[1]}"="${tokens[2]}"
        echo -e "\t${tokens[1]}=${tokens[2]}"
      done <<< "$environ"
    fi

    # Execute commands

    echo "The following commands are requested:"
    while IFS= read -r line; do
      cmdStripped=${line/"[CMD] "/}
      echo -e "\t$cmdStripped"
    done <<< "$(echo "$gimmeCommands" | grep "\[CMD\]")"

    if read -r -p "Ok with you? (Y/N): " confirm && [[ $confirm == [yY] ]] || [[ $confirm == [yY][eE][sS] ]]; then
      while IFS= read -r line; do
        cmdStripped=${line/"[CMD]"/}
        read -r -a tokens <<< "$cmdStripped"
        command="${tokens[0]}"
        args="${tokens[@]:1}"
        echo -e "$(eval "$command" "$args")"
      done <<< "$(echo "$gimmeCommands" | grep "\[CMD\]")"
    else
      echo "No injection commands executed!"
    fi
  fi
fi
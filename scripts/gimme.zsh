#!/bin/zsh

reserved="sync init discover"

if [[ "$reserved" = *"$1"* ]]; then
  gimme-core "$1"
else

  gimmeRaw=$(echo "$@" | xargs gimme-core --manifest)
  if [[ $? == 0 ]]; then
    IFS=' ' read -r -A warp <<< "$(echo "$gimmeRaw" | grep "\[WARP\]")"
    environ="$(echo "$gimmeRaw" | grep "\[ENV\]")"
    gimmeCommands="$(echo "$gimmeRaw" | grep "\[CMD\]")"

    warpDest="${warp[2]}"
    cd "$warpDest" || return

    # Perform exports

    echo "Exporting the following variables:"
    while IFS= read -r line; do
      read -r -A tokens <<< "$line"
      export "${tokens[2]}"="${tokens[3]}"
      echo -e "\t${tokens[2]}=${tokens[3]}"
    done <<< "$environ"

    # Execute commands

    echo "The following commands are requested:"
    while IFS= read -r line; do
      cmdStripped=${line/"[CMD] "/}
      echo -e "\t$cmdStripped"
    done <<< "$(echo "$gimmeCommands" | grep "\[CMD\]")"

    if read -r "confirm?Ok with you? (Y/N)[default no]: " && [[ $confirm == [yY] ]] || [[ $confirm == [yY][eE][sS] ]]; then
      while IFS= read -r line; do
        cmdStripped=${line/"[CMD]"/}
        read -r -A tokens <<< "$cmdStripped"
        command="${tokens[1]}"
        args="${tokens[@]:1}"
        echo -e "$(eval "$command" "$args")"
      done <<< "$(echo "$gimmeCommands" | grep "\[CMD\]")"
    else
      echo "No injection commands executed!"
    fi
  fi
fi